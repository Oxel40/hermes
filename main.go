package main

import (
	"flag"

	"github.com/Oxel40/hermes/internal/configuration"
	"github.com/Oxel40/hermes/internal/logging"
	"github.com/Oxel40/hermes/internal/token"
	"github.com/Oxel40/hermes/internal/web"
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
	// Setup log
	log := logging.GetLogger(*logFileDir)
	// Create empty TokenMaps
	serviceTokenMap := token.TokenMap{make(map[string]string), make(map[string]string)}
	communicatorTokenMap := token.TokenMap{make(map[string]string), make(map[string]string)}

	// Setup Config
	var config configuration.Config
	config.AttatchLogger(log)
	config.AttatchServiceTokenMap(&serviceTokenMap)
	config.AttatchCommunicatorTokenMap(&communicatorTokenMap)
	config.AttatchConfigFile("config.json")

	// Setup WebEndpoint
	var webEndpoint web.Web
	webEndpoint.AttatchLogger(log)
	webEndpoint.AttatchServiceTokenMap(&serviceTokenMap)
	webEndpoint.AttatchCommunicatorTokenMap(&communicatorTokenMap)
	webEndpoint.AttatchConfig(&config)

	// Start subroutines
	go config.Subroutine()
	go webEndpoint.Subroutine(*httpPort)

	/* for {
	time.Sleep(3 * time.Second)
	log.Trace.Println(config)
	} */
	/*
		Trace.Println("I have something standard to say")
		Info.Println("Special Information")
		Warning.Println("There is something you need to know about")
		Error.Println("Something has failed")
	*/
	select {}
}
