package resolver

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func newGRPCServer() (*grpc.Server, string, error) {
	lis, err := net.Listen("tcp", "localhost:")
	if err != nil {
		return nil, "", err
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	go func() {
		_ = s.Serve(lis)
	}()

	return s, lis.Addr().String(), nil
}

func TestProbeExecWithoutError(t *testing.T) {
	assert := assert.New(t)

	grpcServer, addr, err := newGRPCServer()
	assert.NoError(err, "an unexpected error occured in grpc server instantiation")

	defer grpcServer.Stop()

	probe, err := newProbe(addr, time.Second*5)
	assert.NoError(err, "an unexpected error occured in probe instantiation")

	out := probe.exec()

	grpcServer.GracefulStop()

	res := <-out
	assert.Equal(connectivity.Idle, res, "The connectivity states should be idle")

	res = <-out
	assert.Equal(connectivity.Connecting, res, "The connectivity state should be connecting")

	res = <-out
	assert.Equal(connectivity.TransientFailure, res, "The connectivity state should be failure")

	probe.close()
}
