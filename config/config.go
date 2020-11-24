package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sync"
	"verottaa/models"
)

type config struct {
	Port     string   `json:"port"`
	Database database `json:"database"`
}

type database struct {
	Host string `json:"host"`
}

type Configurable interface {
	models.Destroyable
	Portable
	Databaser
}

type Portable interface {
	GetPort() string
}

type Databaser interface {
	GetDatabaseHost() string
}

var destroyCh = make(chan bool)
var getPortCh = make(chan chan string)
var getDatabaseHost = make(chan chan string)

var configInstance *config
var once sync.Once

func GetConfiguration() Configurable {
	once.Do(func() {
		configInstance = createConfig()
		go func() {
			for
			{
				select {
				case ch := <-getPortCh:
					ch <- configInstance.Port
				case ch := <-getDatabaseHost:
					ch <- configInstance.Database.Host
				case <-destroyCh:
					return
				}
			}
		}()
	})

	return configInstance
}

func createDefaultDatabase() *database {
	instance := new(database)
	instance.Host = "mongodb://localhost:27017"
	return instance
}

func createDefaultConfig() *config {
	instance := new(config)
	instance.Port = ":8080"
	instance.Database = *createDefaultDatabase()
	return instance
}

func createConfig() *config {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "config",
			"function": "createConfig",
			"error":    err,
			"cause":    "ReadFile",
		}).Error("Unexpected error")
		conf := createDefaultConfig()
		writeConfigInFile(conf)
		return conf
	}

	instance := config{}
	err = json.Unmarshal(file, &instance)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "config",
			"function": "createConfig",
			"error":    err,
			"cause":    "Unmarshalling",
		}).Error("Unexpected error")
		conf := createDefaultConfig()
		writeConfigInFile(conf)
		return conf
	}

	return &instance
}

func writeConfigInFile(config *config) {
	jsonString, err := json.Marshal(config)
	log.WithFields(log.Fields{
		"package":  "config",
		"function": "writeConfigInFile",
		"error":    err,
		"cause":    "Marshalling",
	}).Error("Unexpected error")

	file, err := os.Create("config.json")
	log.WithFields(log.Fields{
		"package":  "config",
		"function": "writeConfigInFile",
		"error":    err,
		"cause":    "Creating config file",
	}).Error("Unexpected error")

	defer func() {
		err = file.Close()
		log.WithFields(log.Fields{
			"package":  "config",
			"function": "writeConfigInFile",
			"error":    err,
			"cause":    "File cannot close",
		}).Error("Unexpected error")
	}()

	_, err = file.Write(jsonString)
	log.WithFields(log.Fields{
		"package":  "config",
		"function": "writeConfigInFile",
		"error":    err,
		"cause":    "Cannot write file",
	}).Error("Unexpected error")
}

func (c config) Destroy() {
	destroyCh <- true
	close(getPortCh)
	close(destroyCh)
	configInstance = nil
}

func (c config) GetPort() string {
	resCh := make(chan string)
	defer close(resCh)
	getPortCh <- resCh
	return <-resCh
}

func (c config) GetDatabaseHost() string {
	resCh := make(chan string)
	defer close(resCh)
	getDatabaseHost <- resCh
	return <-resCh
}
