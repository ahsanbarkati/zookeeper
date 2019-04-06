package main

import (
	"flag"
	"runtime"

	"github.com/APwhitehat/zookeeper/pkg/simulator"
	"github.com/sirupsen/logrus"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := flag.String("p", "10000", "Port")

	flag.Parse()
	if port == nil {
		logrus.Fatal("Port is nil")
	}

	server := simulator.NewServer(*port)
	server.SetupHTTP()
}
