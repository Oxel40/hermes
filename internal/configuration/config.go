package configuration

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Oxel40/hermes/internal/logging"
	"github.com/Oxel40/hermes/internal/token"
	"github.com/fsnotify/fsnotify"
)

// Config stores configuration information
type Config struct {
	Services             []Service      `json:"services"`
	Recipiens            []Recipient    `json:"recipients"`
	Communicators        []Communicator `json:"communicators"`
	fileDir              string
	communicatorTokenMap *token.TokenMap
	serviceTokenMap      *token.TokenMap
	log                  *logging.Logger
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

// AttatchServiceTokenMap attaches the `tokenMap` to the `Config` to be updated on config updates
func (config *Config) AttatchServiceTokenMap(tokenMap *token.TokenMap) {
	config.serviceTokenMap = tokenMap
}

// AttatchCommunicatorTokenMap attaches the `tokenMap` to the `Config` to be updated on config updates
func (config *Config) AttatchCommunicatorTokenMap(tokenMap *token.TokenMap) {
	config.communicatorTokenMap = tokenMap
}

// AttatchConfigFile loads a config from a file and watches for changes in the config during runtime
func (config *Config) AttatchConfigFile(fileDir string) {
	config.fileDir = fileDir

	err := config.parseConfig(config.fileDir)
	if err != nil {
		config.log.Error.Fatalln("Failed to parse config:", err)
	}
	errs := config.loadTokenMaps()
	for _, err := range errs {
		if err != nil {
			config.log.Error.Fatalln(err)
		}
	}

	config.updateTokenMaps()

	errs = config.saveTokenMaps()
	for _, err := range errs {
		if err != nil {
			config.log.Error.Fatalln(err)
		}
	}

	config.log.Info.Println("Config loaded")
}

// AttatchLogger attatches a `*logging.Logger` to a `Config` to be used in the config subroutine
func (config *Config) AttatchLogger(log *logging.Logger) {
	config.log = log
}

// Subroutine ...
func (config *Config) Subroutine() {
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
					config.updateTokenMaps()
					config.saveTokenMaps()
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

	countMap := make(map[string]int)
	for _, val := range config.getServiceNames() {
		if len(val) <= 0 {
			return errors.New("Services need to have names that are atleast one character long, \"" + val + "\" does not")
		}
		countMap[val]++
		if countMap[val] > 1 {
			return errors.New("Services need to have uniqe names, \"" + val + "\" encountered more than once")
		}
	}

	countMap = make(map[string]int)
	for _, val := range config.getCommunicatorNames() {
		if len(val) <= 0 {
			return errors.New("Communicators need to have names that are atleast one character long, \"" + val + "\" does not")
		}
		countMap[val]++
		if countMap[val] > 1 {
			return errors.New("Communicators need to have uniqe names, \"" + val + "\" encountered more than once")
		}
	}

	return nil
}

func (config *Config) getServiceNames() []string {
	var out []string
	for _, service := range config.Services {
		out = append(out, service.Name)
	}
	return out
}

func (config *Config) getCommunicatorNames() []string {
	var out []string
	for _, communicator := range config.Communicators {
		out = append(out, communicator.Name)
	}
	return out
}

func updateTokenMap(tokenMap *token.TokenMap, names []string) {
	tokenMap.Add(names...)
	keys := tokenMap.GetNames()
	for _, key := range keys {
		isPresent := false
		for _, name := range names {
			if key == name {
				isPresent = true
				break
			}
		}
		if !isPresent {
			tokenMap.Remove(key)
		}
	}
}

func (config *Config) updateTokenMaps() {
	updateTokenMap(config.serviceTokenMap, config.getServiceNames())
	updateTokenMap(config.communicatorTokenMap, config.getCommunicatorNames())
}

func (config *Config) saveTokenMaps() []error {
	err := make([]error, 2)
	err[0] = config.serviceTokenMap.SaveToFile("ServiceTokens.txt")
	err[1] = config.communicatorTokenMap.SaveToFile("CommunicatorTokens.txt")
	return err
}

func (config *Config) loadTokenMaps() []error {
	err := make([]error, 2)
	if e := config.serviceTokenMap.LoadFromFile("ServiceTokens.txt"); !os.IsNotExist(e) {
		err[0] = e
	}
	if e := config.communicatorTokenMap.LoadFromFile("CommunicatorTokens.txt"); !os.IsNotExist(e) {
		err[1] = e
	}
	return err
}
