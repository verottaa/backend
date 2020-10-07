package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"verottaa/models"
	logpack "verottaa/utils/logger"
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

var logTag = "CONFIG"
var logger *logpack.Logger

func init() {
	logger = logpack.CreateLogger(logTag)
}

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
		logger.Error(err)
		conf := createDefaultConfig()
		writeConfigInFile(conf)
		return conf
	}

	instance := config{}
	err = json.Unmarshal(file, &instance)
	if err != nil {
		logger.Error(err)
		conf := createDefaultConfig()
		writeConfigInFile(conf)
		return conf
	}

	return &instance
}

func writeConfigInFile(config *config) {
	jsonString, err := json.Marshal(config)
	logger.Error(err)
	file, err := os.Create("config.json")
	logger.Error(err)
	defer func() {
		err = file.Close()
		logger.Error(err)
	}()
	_, err = file.Write(jsonString)
	logger.Error(err)
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
