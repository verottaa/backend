package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"verottaa/models"
)

type config struct {
	Ports    ports    `json:"ports"`
	Database database `json:"database"`
}

type ports struct {
	Api    string `json:"api"`
	Static string `json:"static"`
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
	GetApiPort() string
	GetStaticPort() string
}

type Databaser interface {
	GetDatabaseHost() string
}

func handleError(err error) {
	if err != nil {
		fmt.Println("[ERROR]:", err)
	}
}

var destroyCh = make(chan bool)
var getApiPortCh = make(chan chan string)
var getStaticPortCh = make(chan chan string)
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
				case ch := <-getApiPortCh:
					ch <- configInstance.Ports.Api
				case ch := <-getStaticPortCh:
					ch <- configInstance.Ports.Static
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

func createDefaultPorts() *ports {
	instance := new(ports)
	instance.Static = ":32678"
	instance.Api = ":8089"
	return instance
}

func createDefaultConfig() *config {
	instance := new(config)
	instance.Ports = *createDefaultPorts()
	instance.Database = *createDefaultDatabase()
	return instance
}

func createConfig() *config {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		handleError(err)
		conf := createDefaultConfig()
		writeConfigInFile(conf)
		return conf
	}

	instance := config{}
	err = json.Unmarshal(file, &instance)
	if err != nil {
		handleError(err)
		conf := createDefaultConfig()
		writeConfigInFile(conf)
		return conf
	}

	return &instance
}

func writeConfigInFile(config *config) {
	jsonString, err := json.Marshal(config)
	handleError(err)
	file, err := os.Create("config.json")
	handleError(err)
	defer func() {
		err = file.Close()
		handleError(err)
	}()
	_, err = file.Write(jsonString)
	handleError(err)
}

func (c config) Destroy() {
	destroyCh <- true
	close(getApiPortCh)
	close(getStaticPortCh)
	close(destroyCh)
	configInstance = nil
}

func (c config) GetApiPort() string {
	resCh := make(chan string)
	defer close(resCh)
	getApiPortCh <- resCh
	return <-resCh
}

func (c config) GetStaticPort() string {
	resCh := make(chan string)
	defer close(resCh)
	getStaticPortCh <- resCh
	return <-resCh
}

func (c config) GetDatabaseHost() string {
	resCh := make(chan string)
	defer close(resCh)
	getDatabaseHost <- resCh
	return <-resCh
}
