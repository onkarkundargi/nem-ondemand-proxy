/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for nem-ondemand-proxy
package main

import (
	"context"
	"github.com/opencord/nem-ondemand-proxy/protos/nem_ondemand_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50052"
)

//server is used to implement the grpc server for the proxy
type server struct {
}

func (s *server) OmciTest(ctx context.Context, id *nem_ondemand_api.OnuID) (*nem_ondemand_api.ResponseTest, error) {
	log.Printf("Request Received from operator client: %s", id.Id)
	resp, err := connect(&id.Id)
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	log.Printf("Result received from voltha-grpc-client: %s", resp.String())
	if len(resp.String()) > 0 {
		return &nem_ondemand_api.ResponseTest{Result: resp.String()}, nil
	}
	return &nem_ondemand_api.ResponseTest{Result: "FAILURE"}, nil
}

func main() {
	log.Printf("Starting nem_grpc_server")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	nem_ondemand_api.RegisterNemServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
