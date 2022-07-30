package ktMicro

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
)

type Config struct {
	yamlFile *ymlFile
	dataMap  map[string]interface{}
	rwMutex  *sync.RWMutex
}

func NewConfig() *Config {
	m := &Config{}
	m.rwMutex = new(sync.RWMutex)
	return m
}

//==========protocol
func (m *Config) Load() {
	m.yamlFile = NewYamlFile("configFiles/config.yml")
	err := m.yamlFile.Load(&m.dataMap)
	if err != nil {
		log.Fatalln("Config Load err:", err)
		return
	}
}

func (m *Config) UpdateConfig() {
}

func (m *Config) GetValue(key string) (data interface{}, err error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	if len(key) == 0 {
		err = errors.New("config GetValue key can not be blank!")
		return
	}
	data = m.dataMap[key]
	return
}

//==========public methods
func (m *Config) GetAppConfigValue(key string) (data interface{}, err error) {
	appConfig, err := m.GetValue("AppConfig")
	if err != nil {
		return
	}
	if mapData, ok := appConfig.(map[interface{}]interface{}); ok {
		for k, v := range mapData {
			if k.(string) == key {
				data = v
				return
			}
		}
	} else {
		err = errors.New(fmt.Sprintf("GetAppConfigData appConfig not a map[interface{}]interface{}:%v", appConfig))
	}
	return
}

func (m *Config) GetAppConfigStringValue(key string) (data string, err error) {
	val, err := m.GetAppConfigValue(key)
	if err != nil {
		return
	}
	if v, ok := val.(string); ok {
		data = v
		return
	} else {
		err = errors.New(fmt.Sprintf("GetAppConfigStringValue value not a string:%v,%v", key, val))
	}
	return
}

func (m *Config) GetAppConfigIntValue(key string) (data int, err error) {
	val, err := m.GetAppConfigValue(key)
	if err != nil {
		return
	}
	if v, ok := val.(int); ok {
		data = v
		return
	} else {
		err = errors.New(fmt.Sprintf("GetAppConfigStringValue value not a int:%v,%v", key, val))
	}
	return
}

func (m *Config) AppIsDev() (data bool) {
	val, _ := m.GetAppConfigStringValue("AppEnv")
	return val == "develop"
}

func (m *Config) AppIsTest() (data bool) {
	val, _ := m.GetAppConfigStringValue("AppEnv")
	return val == "test"
}

func (m *Config) AppIsProduct() (data bool) {
	val, _ := m.GetAppConfigStringValue("AppEnv")
	return val == "product"
}

func (m *Config) AppServerPort() (data string) {
	data, _ = m.GetAppConfigStringValue("ServerPort")
	if strings.Contains(data, ":") {
		return data
	} else {
		return ":" + data
	}
}

//test
func GetValueWithKey() {
	app := NewConfig()
	app.Load()

	val, _ := app.GetAppConfigStringValue("AppEnv")
	fmt.Println("AppEnv val", val)

	valint, _ := app.GetAppConfigIntValue("ServerPort")
	fmt.Println("ServerPort val", valint)

	dev := app.AppIsDev()
	fmt.Println("dev val", dev)
}
