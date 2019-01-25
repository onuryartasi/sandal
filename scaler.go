package main

import (
	"flag"
	"fmt"
	"log"
)
const (
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
	--min  <min-value>   			Minimum container to run (default is 1)
	--max  <max-value>   			Maximum container to run (0 is unlimited, default is 3)

`
func usage() {
	log.Fatalf(InfoColor,usageStr)
}

func main(){


	var image string
	var minValue string
	var maxValue string


	flag.StringVar(&image,"image","","Container's image for scale")
	flag.StringVar(&minValue,"min","1","Minimum container to run (default is 1)")
	flag.StringVar(&maxValue,"max","3","Maximum container to run (0 is unlimited, default is 3)")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if len(image) < 1{

		log.Printf(ErrorColor,"Error: An image must be specified.")

		usage()
	}
	fmt.Printf("Container created with image: %s, min: %s, max: %s",image,minValue,maxValue)

}
