# marathonw

Marathonw is a grpc watcher for Marathon. A minimalistic library that gives the ability to integrate round robin balancing into a grpc client. It uses Marathon framework, which is a container orchestration platform for Apache Mesos and DC/OS, to interact with a service using a name-based system discovery.

# Features

* Round robin load distribution
* Service name discovery (collision supported)
* High availability service discovery with Marathon

# Dependencies

* [Marathon](https://mesosphere.github.io/marathon): A production-grade container orchestration platform for Mesosphere's Datacenter

# Installation

Install Marathonw using the "go get" command:

`go get github.com/eddyzags/marathonw`

Import the library into a project:

`import "github.com/eddyzags/marathonw"`

# Usage

Marathonw uses marathon labels in order to contacts another service. Let's start with an simple app definition with a marathonw label

```
{
  "id": "my-app"
  "cpus": 0.1,
  "mem": 16,
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
      }
    ]
  },
  "labels": {
    "MARATHONW_0_NAME": "my-app-discovery"
  }
}
```

Here we defined an application called my-app with a marathonw service discovery name my-app-discovery which points to the port index 0.

Service discovery name can be defined using a label variable:

`"MARATHONW_{PORTINDEX}_NAME": "{VALUE}"`

Once we deployed the application in Marathon, the service can me discovered through its name in the grpc client instantiation.

```golang
func main() {
     b := grpc.RoundRobin(marathonw.NewResolver("http://marathon.mesos:8080"))

     conn, err := grpc.Dial("my-app-discovery", grpc.WithBalancer(b))
     if err != nil {
       log.Fatalf("Failed to dial grpc server: %v", err)
     }
     defer conn.Close()

     c := pb.NewGreeterClient(conn)

     r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: *name})
     if err != nil {
         log.Fatalf("Failed to send say hello request: %v", err)
     }

     fmt.Printf("Response: %s\n", r.Message)
}
```