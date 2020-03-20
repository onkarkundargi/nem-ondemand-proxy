/*
 * Copyright 2018-present Open Networking Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Implements a client for nem-ondemend-proxy.
package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	pb "github.com/opencord/voltha-protos/v3/go/voltha"
	"google.golang.org/grpc"
	"log"
	"time"
)

/*
 * TODO: Consider refactoring so that the kafka and grpc clients are
 * initialized once rather than for each request that is handled.
 *
 */

type OnDemandHandler struct {
}

func NewOnDemandHandler() *OnDemandHandler {
	var handler OnDemandHandler
	return &handler
}

func (handler *OnDemandHandler) HandleRequest(device_id *string) (*pb.Event, error) {
	// Set up a connection to the server.
	log.Printf("voltha grpc client started, address=%s ...", GlobalConfig.Server)
	conn, err := grpc.Dial(GlobalConfig.Server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()
	c := pb.NewVolthaServiceClient(conn)
	id, err := uuid.NewUUID()
	log.Printf("ID: %s", id.String())
	if err != nil {
		log.Printf("did not generate uuid: %v", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Printf("Calling StartOmciTestAction")
	r, err := c.StartOmciTestAction(ctx, &pb.OmciTestRequest{Id: *device_id, Uuid: id.String()})
	if err != nil {
		return nil, fmt.Errorf("start-omci-test-action-failed: %v", err)
	}
	log.Printf("Result: %s", r.Result)
	djson, _ := json.Marshal(r.Result)
	result := &pb.Event{}
	if string(djson) == "0" {
		config := sarama.NewConfig()
		config.ClientID = "go-kafka-consumer"
		config.Consumer.Return.Errors = true

		brokers := []string{GlobalConfig.Kafka}
		// Create new consumer
		master, err := sarama.NewConsumer(brokers, config)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := master.Close(); err != nil {
				panic(err)
			}
		}()

		topic := "voltha.events"
		consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		// Get signnal for finish
		doneCh := make(chan struct{})
		go func() {
			// TODO: Needs a timeout in here
			for {
				select {
				case err := <-consumer.Errors():
					fmt.Println(err)
				case msg := <-consumer.Messages():
					unpackResult := &pb.Event{}
					var err error
					if err = proto.Unmarshal(msg.Value, unpackResult); err != nil {
						fmt.Println("Error while doing unmarshal", err)
					}
					kpi_event2 := unpackResult.GetKpiEvent2()
					if (kpi_event2 != nil) && (kpi_event2.SliceData[0].Metadata.Uuid == id.String()) {
						result = unpackResult
						close(doneCh)
						return
					}
				}
			}
		}()
		<-doneCh
		log.Printf("Result: %s", result)
	}
	return result, nil
}
