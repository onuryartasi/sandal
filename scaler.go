package main

import (
	"flag"
	"log"
	"github.com/onuryartasi/scaler-api"
	"os"
	"fmt"
)

var usageStr = `
Usage: scaler [options]

Options:
	--image   <image-url>           Container's image for scale
	--min  	  <min-value>   	    Minimum container to run (default is 1)
	--max     <max-value>   		Maximum container to run (0 is unlimited, default is 3)
`
func usage() {
	log.Fatalf(scaler_api.InfoColor,usageStr)
}

func main(){
	var image string
	var minValue string
	var maxValue string

	list := flag.NewFlagSet("list",flag.ExitOnError)
	create := flag.NewFlagSet("create",flag.ExitOnError)
	create.StringVar(&image,"image","","Container's image for scale")
	create.StringVar(&minValue,"min","1","Minimum container to run (default is 1)")
	create.StringVar(&maxValue,"max","3","Maximum container to run (0 is unlimited, default is 3)")

	log.SetFlags(0)
	flag.Usage = usage

	switch os.Args[1] {
	case "list":
		list.Parse(os.Args[2:])
		containers := scaler_api.ContainerList()
		for _,container := range containers{
			fmt.Println(container.Names,container.Image)
		}
	case "create":
		create.Parse(os.Args[2:])
		if len(image) < 1{
			log.Printf(scaler_api.ErrorColor,"Error: An image must be specified.")
			usage()
		}
		log.Printf("Container created with image: %s, min: %s, max: %s",image,minValue,maxValue)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}




}
