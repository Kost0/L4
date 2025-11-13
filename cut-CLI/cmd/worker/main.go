package main

import (
	"log"
	"net"

	"github.com/Kost0/L4/internal/gen/tasks"
	"github.com/Kost0/L4/internal/workers"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

func main() {
	port := pflag.StringP("port", "p", ":8081", "port to listen on")

	pflag.Parse()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	tasks.RegisterDistributedCutServer(s, &workers.WorkerServer{})

	log.Printf("Starting worker on port %s", *port)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
