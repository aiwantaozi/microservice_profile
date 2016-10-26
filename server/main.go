package main

import (
	"flag"
	configs "github.com/aiwantaozi/microservice_profile/configs"
	interceptor "github.com/aiwantaozi/microservice_profile/middleware"
	profile_proto "github.com/aiwantaozi/microservice_profile/proto/profile"
	profile "github.com/aiwantaozi/microservice_profile/server/profile"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/mwitkow/go-grpc-middleware"
	"google.golang.org/grpc"
	"net"
)

func Run() error {
	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptor.LoggerStreamInterceptor, grpc_prometheus.StreamServerInterceptor, interceptor.AuthStreamInterceptor)))

	s := grpc.NewServer(opts...)

	profile_proto.RegisterProfileServiceServer(s, profile.NewProfileServer())

	s.Serve(l)
	return nil
}

func main() {
	configs.InitMongo()
	flag.Parse()
	defer glog.Flush()

	if err := Run(); err != nil {
		glog.Fatal(err)
	}

}
