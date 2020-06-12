package main

import (
	"flag"
	"io"
	"log"
	"os"
)

const (
	defaultToken  = "THIS_IS_NOT_A_VALID_TOKEN"
	defaulLogDir  = "log.txt"
	defaultWsPort = 8080
)

var (
	discordToken            = flag.String("dctoken", defaultToken, "token to connect Discord bot")
	logFileDir              = flag.String("log", defaulLogDir, "log file location")
	disableDiscordClient    = flag.Bool("no-dcbot", false, "disable integrated Discord bot")
	enableWebsocketEndpoint = flag.Bool("ws", false, "enable websocket endpoint")
	websocketPort           = flag.Int("ws-port", defaultWsPort, "websocket port")

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	flag.Parse()

	file, err := os.OpenFile(*logFileDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	multi := io.MultiWriter(file, os.Stdout)

	Trace = log.New(multi,
		"[TRACE] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(multi,
		"[INFO] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(multi,
		"[WARNING] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(multi,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	log.New(file, "----------\n[STARTING] ", log.Ldate|log.Ltime).Println("")
}

func main() {
	logFlagInfo()
	/*
		Trace.Println("I have something standard to say")
		Info.Println("Special Information")
		Warning.Println("There is something you need to know about")
		Error.Println("Something has failed")
	*/

	Trace.Println(*discordToken)
	Trace.Println(*logFileDir)
	Trace.Println(*disableDiscordClient)
	Trace.Println(*enableWebsocketEndpoint)
	Trace.Println(*websocketPort)
}

func logFlagInfo() {
	if *logFileDir != defaulLogDir {
		Info.Println("Log file location set to:", *logFileDir)
	}

	if *disableDiscordClient {
		Info.Println("Integrated Discord bot disabled")
	}

	if *discordToken != defaultToken {
		if *disableDiscordClient {
			Warning.Println("Discord bot token set but the Discord bot is disabled")
		} else {
			Info.Println("Discord bot token set")
		}
	} else {
		if !*disableDiscordClient {
			Warning.Println("No Discord bot token set")
		}
	}

	if *enableWebsocketEndpoint {
		Info.Println("Websocket endpoint enabled")
	}

	if *websocketPort != defaultWsPort {
		if *enableWebsocketEndpoint {
			Info.Println("Websocket port set to:", *websocketPort)
		} else {
			Warning.Println("Websocket port set to:", *websocketPort, "but websocket is not enabled")
		}
	}
}
