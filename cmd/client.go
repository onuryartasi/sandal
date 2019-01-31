package main

import (
	"google.golang.org/grpc"
	"github.com/onuryartasi/scaler/pkg/api/v1"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"fmt"
	"log"
	"flag"
	"os"
	"io"
)

type container struct {
	*v1.Container
}

var (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)


var usageStr = `
Usage: scaler [options]

Options:
	--image   <image-url>           Container's image for scale
	--min  	  <min-value>   	    Minimum container to run (default is 1)
	--max     <max-value>   		Maximum container to run (0 is unlimited, default is 3)
`

func usage() {
	log.Fatalf(InfoColor,usageStr)
}

func connect() v1.ContainerServiceClient{
	conn,err := grpc.Dial(":4444",grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := v1.NewContainerServiceClient(conn)
	return client
}
func main(){
	var image string
	var minValue string
	var maxValue string

	flag.NewFlagSet("list",flag.ExitOnError)
	flag.NewFlagSet("stop",flag.ExitOnError)
	create := flag.NewFlagSet("create",flag.ExitOnError)
	create.StringVar(&image,"image","","Container's image for scale")
	create.StringVar(&minValue,"min","1","Minimum container to run (default is 1)")
	create.StringVar(&maxValue,"max","3","Maximum container to run (0 is unlimited, default is 3)")

	log.SetFlags(0)
	flag.Usage = usage
	if len(os.Args)<2{
		usage()
	}
	switch os.Args[1] {
	case "list":

		client := connect()
		resp, err := client.ContainerList(context.Background(),&empty.Empty{})
		if err != nil {
			log.Fatalf("Container List Error : %s",err)
		}
		for _,container := range resp.GetContainer(){
			fmt.Println(container.Id,container.Names,container.Image)
		}

	case "create":
		create.Parse(os.Args[2:])
		if len(image) < 1{
			log.Printf(ErrorColor,"Error: An image must be specified.")
			usage()
		}
		log.Printf("Container created with image: %s, min: %s, max: %s",image,minValue,maxValue)
		client := connect()
		resp,err := client.ContainerCreate(context.Background(),&v1.ContainerConfig{Image:image})
		if err != nil{
			log.Printf(ErrorColor,"Error: Contaner Create error")
		}
		fmt.Println(resp.GetId())
	case "start":
		client := connect()
		containerId := string(os.Args[2])
		resp,err := client.ContainerStart(context.Background(),&v1.ContainerId{ContainerId:containerId})
		if err != nil{
			log.Printf(ErrorColor,"Error: Contaner start error")
		}
		fmt.Println(resp)
	case "stop":
		client := connect()
		containerId := string(os.Args[2])
		resp,err := client.ContainerStop(context.Background(),&v1.ContainerId{ContainerId:containerId})
		if err != nil{
			log.Printf(ErrorColor,"Error: Contaner Stop error")
		}
		fmt.Println(resp)
	case "stat":
		client := connect()
		containerId := string(os.Args[2])
		stream,err := client.ContainerStatStream(context.Background(),&v1.ContainerId{ContainerId:containerId})
		if err != nil {
			panic(err)
		}
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			log.Println(data)
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.SetFlags(0)
	flag.Usage = usage
}
