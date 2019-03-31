package main

import (
	"flag"
	"runtime"

	"github.com/APwhitehat/zookeeper/pkg/bimock"
	"github.com/sirupsen/logrus"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := flag.String("p", "9000", "Port")

	flag.Parse()
	if port == nil {
		logrus.Fatal("Port is nil")
	}

	server := bimock.NewServer(*port)
	server.SetupHTTP()
}
