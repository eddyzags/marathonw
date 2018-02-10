# Resolver

Resolver is a grpc resolver using Marathon framework for service discovery. A minimalistic library that uses Marathon framework, which is a container orchestration platform for Apache Mesos and Datacenter Operating (DC/OS), to interact with other services using name-based system discovery. Plus, when scaling service, It gives the ability to integrate round robin balancing into a grpc client.

# Features

* Round robin load distribution
* Service name discovery (collision supported)
* High availability service discovery with Marathon

# Dependencies

* [Marathon](https://mesosphere.github.io/marathon): A production-grade container orchestration platform for Mesosphere's Datacenter.
* [gRPC-Go](https://github.com/grpc/grpc-go): Go implementation of gRPC. A high performance, open source, general RPC framework.

# Installation

Install the resolver using the "go get" command:

`go get github.com/eddyzags/resolver`

Import the library into a project:

`import "github.com/eddyzags/resolver"`

# Usage

Resolver uses Marathon label feature in order to identify services.
The marathon label is composed of the service unique name and the port
index. Those two informations will allow the resolver to identify the
service's tasks and on which port the grpc client should establish a connection.
Let's start with a simple app definition:

```json
{
  "id": "my-app"
  "cpus": 0.1,
  "mem": 64,
  "container": {
    "type": "DOCKER",
    "docker": {
      "image": "eddyzags/healthy:latest",
      "network": "HOST"
    },
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 0
      },
      {
        "containerPort": 4242,
        "hostPort": 0
      }
    ]
  },
  "labels": {
    "RESOLVER_0_NAME": "my-app-service"
  }
}
```

Here, we have just defined an application called my-app with a service resolver name `my-app-service` which points to the port index 0.

A service resolver name can be defined using a labels map:

`"RESOLVER_{PORTINDEX}_NAME": "{NAME}"`

Once we deployed the application in Marathon, the service can be discovered through its name in the grpc client instantiation.

```golang
package main

import (
       "log"

       "github.com/eddyzags/resolver"

       "google.golang.org/grpc"
       pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
     resolver, err := resolver.New("marathon.mesos:8080")
     if err != nil {
        log.Fatalf("couldn't instantiation resolver: %v", err)
     }

     b := grpc.RoundRobin(resolver)

     conn, err := grpc.Dial("my-app-discovery", grpc.WithBalancer(b))
     if err != nil {
        log.Fatalf("couldn't dial grpc server: %v", err)
     }
     defer conn.Close()

     c := pb.NewGreeterClient(conn)

     r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: *name})
     if err != nil {
         log.Fatalf("couldn't send say hello request: %v", err)
     }

     log.Printf("Response: %s\n", r.Message)
}
```
