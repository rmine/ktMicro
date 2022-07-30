package ktMicro

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type ConfigSwitch struct {
	yamlFile *ymlFile
	dataMap  map[string]bool
	rwMutex  *sync.RWMutex
}

func NewConfigSwitch() *ConfigSwitch {
	m := &ConfigSwitch{}
	m.rwMutex = new(sync.RWMutex)
	return m
}

//==========protocol
func (m *ConfigSwitch) Load() {
	m.yamlFile = NewYamlFile("configFiles/configSwitch.yml")
	err := m.yamlFile.Load(&m.dataMap)
	if err != nil {
		log.Fatalln("ConfigSwitch Load err:", err)
		return
	}
}

func (m *ConfigSwitch) UpdateConfig() {
}

func (m *ConfigSwitch) GetValue(key string) (data interface{}, err error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	if len(key) == 0 {
		err = errors.New("configSwitch GetValue key can not be blank!")
		return
	}
	data = m.dataMap[key]
	return
}

func (m *ConfigSwitch) GetConfigSwitch(key string) (data bool, err error) {
	val, err := m.GetValue(key)
	if err != nil {
		return
	}
	if v, ok := val.(bool); ok {
		data = v
		return
	} else {
		err = errors.New(fmt.Sprintf("configSwitch GetConfigSwitch value not a bool:%v,%v", key, val))
	}
	return
}
