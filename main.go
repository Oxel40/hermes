package main

import (
	"flag"
	"time"

	"github.com/Oxel40/hermes/internal/configuration"
	"github.com/Oxel40/hermes/internal/logging"
	"github.com/Oxel40/hermes/internal/token"
)

const (
	defaulLogDir  = "log.txt"
	defaultWsPort = 8080
)

var (
	logFileDir    = flag.String("log", defaulLogDir, "log file location")
	websocketPort = flag.Int("ws-port", defaultWsPort, "websocket port")

	config   configuration.Config
	tokenMap token.TokenMap
	log      *logging.Logger
)

func init() {
	// Parse flags
	flag.Parse()

	// Setup log
	log = logging.GetLogger(*logFileDir)

	tokenMap = token.TokenMap{make(map[string]string), make(map[string]string)}
	config.AttatchTokenMap(&tokenMap)
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
