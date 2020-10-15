package config

import (
	"encoding/json"
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
		// TODO: логирование
		conf := createDefaultConfig()
		writeConfigInFile(conf)
		return conf
	}

	instance := config{}
	err = json.Unmarshal(file, &instance)
	if err != nil {
		// TODO: логирование
		conf := createDefaultConfig()
		writeConfigInFile(conf)
		return conf
	}

	return &instance
}

func writeConfigInFile(config *config) {
	jsonString, err := json.Marshal(config)
	// TODO: логирование
	file, err := os.Create("config.json")
	// TODO: логирование
	defer func() {
		err = file.Close()
		// TODO: логирование
	}()
	_, err = file.Write(jsonString)
	// TODO: логирование
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
