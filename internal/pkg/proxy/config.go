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

// Implements global configuration for nem-ondemand-proxy
package proxy

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type OutputType uint8

type GrpcConfigSpec struct {
	Timeout time.Duration `yaml:"timeout"`
}

type TlsConfigSpec struct {
	UseTls bool   `yaml:"useTls"`
	CACert string `yaml:"caCert"`
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
	Verify string `yaml:"verify"`
}

type GlobalConfigSpec struct {
	Server string         `yaml:"server"`
	Kafka  string         `yaml:"kafka"`
	Local  string         `yaml:"local"`
	Tls    TlsConfigSpec  `yaml:"tls"`
	Grpc   GrpcConfigSpec `yaml:"grpc"`
}

var (
	CharReplacer = strings.NewReplacer("\\t", "\t", "\\n", "\n")

	GlobalConfig = GlobalConfigSpec{
		Server: "voltha-rw-core.voltha:50057",
		Kafka:  "voltha-kafka.voltha:9092",
		Local:  "0.0.0.0:50052",
		Tls: TlsConfigSpec{
			UseTls: false,
		},
		Grpc: GrpcConfigSpec{
			Timeout: time.Minute * 5,
		},
	}

	GlobalCommandOptions = make(map[string]map[string]string)

	GlobalOptions struct {
		Config string `short:"c" long:"config" env:"PROXYCONFIG" value-name:"FILE" default:"" description:"Location of proxy config file"`
		Server string `short:"s" long:"server" default:"" value-name:"SERVER:PORT" description:"IP/Host and port of VOLTHA"`
		Kafka  string `short:"k" long:"kafka" default:"" value-name:"SERVER:PORT" description:"IP/Host and port of Kafka"`
		Local  string `short:"l" long:"local" default:"" value-name:"SERVER:PORT" description:"IP/Host and port to listen on"`

		// The following are not necessarily implemented yet.
		Debug   bool   `short:"d" long:"debug" description:"Enable debug mode"`
		Timeout string `short:"t" long:"timeout" description:"API call timeout duration" value-name:"DURATION" default:""`
		UseTLS  bool   `long:"tls" description:"Use TLS"`
		CACert  string `long:"tlscacert" value-name:"CA_CERT_FILE" description:"Trust certs signed only by this CA"`
		Cert    string `long:"tlscert" value-name:"CERT_FILE" description:"Path to TLS vertificate file"`
		Key     string `long:"tlskey" value-name:"KEY_FILE" description:"Path to TLS key file"`
		Verify  bool   `long:"tlsverify" description:"Use TLS and verify the remote"`
	}

	Debug = log.New(os.Stdout, "DEBUG: ", 0)
	Info  = log.New(os.Stdout, "INFO: ", 0)
	Warn  = log.New(os.Stderr, "WARN: ", 0)
	Error = log.New(os.Stderr, "ERROR: ", 0)
)

func ParseCommandLine() {
	parser := flags.NewNamedParser(path.Base(os.Args[0]),
		flags.HelpFlag|flags.PassDoubleDash|flags.PassAfterNonOption)
	_, err := parser.AddGroup("Global Options", "", &GlobalOptions)
	if err != nil {
		panic(err)
	}

	_, err = parser.ParseArgs(os.Args[1:])
	if err != nil {
		_, ok := err.(*flags.Error)
		if ok {
			real := err.(*flags.Error)
			if real.Type == flags.ErrHelp {
				os.Stdout.WriteString(err.Error() + "\n")
				os.Exit(0)
			}
		}

		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err.Error())

		os.Exit(1)
	}
}

func ProcessGlobalOptions() {
	if len(GlobalOptions.Config) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			Warn.Printf("Unable to discover the user's home directory: %s", err)
			home = "~"
		}
		GlobalOptions.Config = filepath.Join(home, ".nem", "config")
	}

	if info, err := os.Stat(GlobalOptions.Config); err == nil && !info.IsDir() {
		configFile, err := ioutil.ReadFile(GlobalOptions.Config)
		if err != nil {
			Error.Fatalf("Unable to read the configuration file '%s': %s",
				GlobalOptions.Config, err.Error())
		}
		if err = yaml.Unmarshal(configFile, &GlobalConfig); err != nil {
			Error.Fatalf("Unable to parse the configuration file '%s': %s",
				GlobalOptions.Config, err.Error())
		}
	}

	// Override from command line
	if GlobalOptions.Server != "" {
		GlobalConfig.Server = GlobalOptions.Server
	}
	if GlobalOptions.Kafka != "" {
		GlobalConfig.Kafka = GlobalOptions.Kafka
	}
	if GlobalOptions.Local != "" {
		GlobalConfig.Local = GlobalOptions.Local
	}

	if GlobalOptions.Timeout != "" {
		timeout, err := time.ParseDuration(GlobalOptions.Timeout)
		if err != nil {
			Error.Fatalf("Unable to parse specified timeout duration '%s': %s",
				GlobalOptions.Timeout, err.Error())
		}
		GlobalConfig.Grpc.Timeout = timeout
	}
}

func ShowGlobalOptions() {
	log.Printf("Configuration:")
	log.Printf("    Voltha gRPC Server: %v", GlobalConfig.Server)
	log.Printf("    Kafka: %v", GlobalConfig.Kafka)
	log.Printf("    Listen Address: %v", GlobalConfig.Local)
}
