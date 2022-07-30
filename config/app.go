package ktMicro

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type App struct {
	AppEnv string

	yamlFile *ymlFile
	dataMap  map[string]interface{}
	rwMutex  *sync.RWMutex
}

func NewApp() *App {
	m := &App{}
	m.rwMutex = new(sync.RWMutex)
	return m
}

func (m *App) Load() {
	path := fmt.Sprintf("configFiles/%v/app.yml", m.AppEnv)
	m.yamlFile = NewYamlFile(path)
	err := m.yamlFile.Load(&m.dataMap)
	if err != nil {
		log.Fatalln("App Load err:", err)
		return
	}
}

func (m *App) UpdateConfig() {
}

func (m *App) GetValue(key string) (data interface{}, err error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	if len(key) == 0 {
		err = errors.New("app GetValue key can not be blank!")
		return
	}
	data = m.dataMap[key]
	return
}

//==========public methods
func (m *App) GetStringValue(key string) (data string, err error) {
	val, err := m.GetValue(key)
	if err != nil {
		return
	}
	if v, ok := val.(string); ok {
		data = v
		return
	} else {
		err = errors.New(fmt.Sprintf("app GetStringValue value not a string:%v,%v", key, val))
	}
	return
}
