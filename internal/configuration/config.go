package configuration

import (
	"encoding/json"
	"os"

	"github.com/Oxel40/hermes/internal/logging"
	"github.com/Oxel40/hermes/internal/token"
	"github.com/fsnotify/fsnotify"
)

// Config stores configuration information
type Config struct {
	Services         []Service      `json:"services"`
	Recipiens        []Recipient    `json:"recipients"`
	Communicators    []Communicator `json:"communicator"`
	DiscordBot       DiscordBot     `json:"discord-bot"`
	fileDir          string
	attachedTokenMap *token.TokenMap
	log              *logging.Logger
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

// Communicator ...
type Communicator struct {
	Name    string `json:"name"`
	IDIndex int    `json:"id-index"`
}

// DiscordBot ...
type DiscordBot struct {
	Token   string `json:"token"`
	IDIndex int    `json:"id-index"`
}

// AttatchTokenMap attaches the `tokenMap` to the `Config` to be updated on config updates
func (config *Config) AttatchTokenMap(tokenMap *token.TokenMap) {
	config.attachedTokenMap = tokenMap
}

// AttatchConfigFile loads a config from a file and watches for changes in the config during runtime
func (config *Config) AttatchConfigFile(fileDir string) {
	config.fileDir = fileDir

	config.parseConfig(fileDir)
	config.attachedTokenMap.LoadFromFile("tokens.txt")

	config.updateTokenMap()
	config.attachedTokenMap.SaveToFile("tokens.txt")
	config.log.Info.Println("Config loaded")
}

// AttatchLogger attatches a `*logging.Logger` to a `Config` to be used in the config subroutine
func (config *Config) AttatchLogger(log *logging.Logger) {
	config.log = log
}

// StartConfigSubroutine ...
func (config *Config) StartConfigSubroutine() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		config.log.Error.Fatal(err)
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				config.log.Info.Println("Event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					config.log.Info.Println("Config file modified:", event.Name)
					err = config.parseConfig(config.fileDir)
					if err != nil {
						config.log.Error.Fatalln("Failed to parse config:", err)
					}
					config.updateTokenMap()
					config.attachedTokenMap.SaveToFile("tokens.txt")
					config.log.Info.Println("Config reloaded")

				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					config.log.Error.Fatalln("Config file removed:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				config.log.Error.Println(err)
			}
		}
	}()

	err = watcher.Add(config.fileDir)
	if err != nil {
		config.log.Error.Fatal(err)
	}
}

func (config *Config) parseConfig(fileDir string) error {
	// Parse config
	confReader, err := os.Open(fileDir)
	defer confReader.Close()
	if err != nil {
		return err
	}
	err = json.NewDecoder(confReader).Decode(config)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) getAllServiceAndCommunicatorNames() []string {
	var out []string
	for _, service := range config.Services {
		out = append(out, service.Name)
	}
	for _, communicator := range config.Communicators {
		out = append(out, communicator.Name)
	}
	return out
}

func (config *Config) updateTokenMap() {
	names := config.getAllServiceAndCommunicatorNames()
	config.attachedTokenMap.Add(names...)
	keys := config.attachedTokenMap.GetNames()
	for _, key := range keys {
		isPresent := false
		for _, name := range names {
			if key == name {
				isPresent = true
				break
			}
		}
		if !isPresent {
			config.attachedTokenMap.Remove(key)
		}
	}
}
