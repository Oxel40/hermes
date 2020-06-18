package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

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

	log.Info.Println("hermes is now running. Press CTRL-C to exit.")
	// Wait here until CTRL-C or other term signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Info.Println("hermes terminated, exiting...")

	// Cleanly close down the webendpoint
	webEndpoint.Close()
}
