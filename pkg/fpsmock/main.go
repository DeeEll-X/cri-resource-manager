// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"
	"net"

	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
	pb "github.com/intel/cri-resource-manager/pkg/cri/fpsserver"
)

var (
	addr = flag.String("addr", "/var/run/cri-resmgr/cri-resmgr-fps.sock", "the address to connect to")
)

func main() {
	flag.Parse()

	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial("unix", addr)
	}

	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithDialer(dialer))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Set up a connection to the server.
	c := pb.NewFPSServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for true {
		_, err = c.FPSDrop(ctx, &pb.FPSDropRequest{Id: "12"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting")
        time.Sleep(time.Second)
    }	
}