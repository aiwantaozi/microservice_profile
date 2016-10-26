package main

import (
	"flag"
	"net/http"

	profile "github.com/aiwantaozi/microservice_profile/proto/profile"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	profileEndpoint = flag.String("profile_endpoint", "localhost:9092", "endpoint of ProfileService")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	err := profile.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, *profileEndpoint, dialOpts)
	if err != nil {
		return err
	}

	http.ListenAndServe(":8082", mux)
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
