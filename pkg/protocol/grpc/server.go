package grpc

import (
	"github.com/onuryartasi/scaler/pkg/api/v1"
	"net"
	"log"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	service "github.com/onuryartasi/scaler/pkg/service/v1"
)


func RunServer() error{
	listen, err := net.Listen("tcp",":4444")
	if err != nil {
		log.Fatalf("Run server error : &v",err)
	}

	s := &service.Service{}
	server := grpc.NewServer()
	v1.RegisterContainerServiceServer(server,s)
	c := make(chan os.Signal,1)
	signal.Notify(c,os.Interrupt)

	go func() {for range c{
		log.Println("Shutting gRPC Server...")
		server.GracefulStop()
		<-c
	}}()

	log.Println("Starting gRPC Server")

	return server.Serve(listen)
}