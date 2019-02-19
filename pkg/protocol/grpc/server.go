package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/docker/docker/client"
	v1 "github.com/onuryartasi/scaler/pkg/api/v1"
	service "github.com/onuryartasi/scaler/pkg/service/v1"
	"google.golang.org/grpc"
)

var timeout time.Duration = 10 * time.Second

func RunServer() error {
	listen, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatalf("Run server error : &v", err)
	}

	s := &service.Service{}
	server := grpc.NewServer()
	v1.RegisterContainerServiceServer(server, s)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("Shutting gRPC Server and stopping container...")
			cli, err := client.NewClientWithOpts(client.FromEnv)
			if err != nil {
				log.Println("Container cants stop")
			}
			projects := *service.GetProjects()
			for _, project := range projects {
				for _, containerID := range project.Containers {
					cli.ContainerStop(context.Background(), containerID, &timeout)
				}
			}
			server.GracefulStop()
			<-c
		}
	}()

	log.Println("Starting gRPC Server")

	return server.Serve(listen)
}
