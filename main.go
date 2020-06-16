package main

import (
	"flag"
	"time"

	"github.com/Oxel40/hermes/internal/configuration"
	"github.com/Oxel40/hermes/internal/logging"
	"github.com/Oxel40/hermes/internal/token"
)

const (
	defaulLogDir = "log.txt"
	defaultPort  = 8080
)

var (
	logFileDir = flag.String("log", defaulLogDir, "log file location")
	httpPort   = flag.Int("port", defaultPort, "http/websocket port")
)

func init() {
	// Parse flags
	flag.Parse()
}

func main() {
	var config configuration.Config
	var tokenMap token.TokenMap
	var log *logging.Logger

	// Setup log
	log = logging.GetLogger(*logFileDir)
	// Create empty TokenMap
	tokenMap = token.TokenMap{make(map[string]string), make(map[string]string)}

	config.AttatchLogger(log)
	config.AttatchTokenMap(&tokenMap)
	config.AttatchConfigFile("config.json")

	config.StartConfigSubroutine()
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
