package configuration

import (
	"encoding/json"
	"os"

	"github.com/Oxel40/hermes/internal/logging"
	"github.com/fsnotify/fsnotify"
)

// Config stores configuration information
type Config struct {
	Services   []Service    `json:"services"`
	Recipiens  []Recipient  `json:"recipients"`
	WSClients  []WSClient   `json:"ws-clients"`
	DiscordBot []DiscordBot `json:"discord-bot"`
}

// Service ...
type Service struct {
	Name string `json:"name"`
}

// Recipient ...
type Recipient struct {
	Name          string   `json:"name"`
	IDs           []string `json:"ids"`
	Subscriptions []string `json:"subscriptions"`
}

// WSClient ...
type WSClient struct {
	Name    string `json:"name"`
	IDIndex int    `json:"id-index"`
}

// DiscordBot ...
type DiscordBot struct {
	Token   string `json:"token"`
	IDIndex int    `json:"id-index"`
}

// AttatchConfigFile loads a config from a file and watches for changes in the config during runtime
func (config *Config) AttatchConfigFile(fileDir string, log *logging.Logger) {
	config.parseConfig(fileDir, log)
	log.Info.Println("Config loaded")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error.Fatal(err)
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Info.Println("Event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Info.Println("Config file modified:", event.Name)
					config.parseConfig(fileDir, log)
					log.Info.Println("Config reloaded")
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Error.Fatalln("Config file removed:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error.Println(err)
			}
		}
	}()

	err = watcher.Add(fileDir)
	if err != nil {
		log.Error.Fatal(err)
	}
}

func (config *Config) parseConfig(fileDir string, log *logging.Logger) {
	// Parse config
	confReader, err := os.Open(fileDir)
	defer confReader.Close()
	if err != nil {
		log.Error.Fatalln("Failed to open config file:", err)
	}
	json.NewDecoder(confReader).Decode(config)
}
