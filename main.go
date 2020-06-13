package main

import (
	"flag"
	"time"

	"github.com/Oxel40/hermes/internal/configuration"
	"github.com/Oxel40/hermes/internal/logging"
)

const (
	defaulLogDir  = "log.txt"
	defaultWsPort = 8080
)

var (
	logFileDir    = flag.String("log", defaulLogDir, "log file location")
	websocketPort = flag.Int("ws-port", defaultWsPort, "websocket port")

	config configuration.Config
	log    *logging.Logger
)

func init() {
	// Parse flags
	flag.Parse()

	// Setup log
	log = logging.GetLogger(*logFileDir)

	config.AttatchConfigFile("config.json", log)
}

func main() {
	/*
		Trace.Println("I have something standard to say")
		Info.Println("Special Information")
		Warning.Println("There is something you need to know about")
		Error.Println("Something has failed")
	*/
	for {
		time.Sleep(3 * time.Second)
		log.Trace.Println(config)
	}
}
