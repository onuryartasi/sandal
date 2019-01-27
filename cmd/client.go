package main

import (
	"google.golang.org/grpc"
	"github.com/onuryartasi/scaler/pkg/api/v1"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"fmt"
)

type container struct {
	*v1.Container
}

func main(){
	conn,err := grpc.Dial(":4444",grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := v1.NewContainerListServiceClient(conn)

	resp, err := client.ContainerList(context.Background(),&empty.Empty{})
	for _,container := range resp.GetContainer(){
		fmt.Println(container.Id,container.Names,container.Image)
	}

}
