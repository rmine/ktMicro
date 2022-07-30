package ktMicro

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type Log struct {
	AppEnv string

	yamlFile *ymlFile
	dataMap  map[string]interface{}
	rwMutex  *sync.RWMutex
}

func NewLog() *Log {
	m := &Log{}
	m.rwMutex = new(sync.RWMutex)
	return m
}

//==========protocol
func (m *Log) Load() {
	path := fmt.Sprintf("configFiles/%v/log.yml", m.AppEnv)
	m.yamlFile = NewYamlFile(path)
	err := m.yamlFile.Load(&m.dataMap)
	if err != nil {
		log.Fatalln("Log Load err:", err)
		return
	}
}

func (m *Log) UpdateConfig() {
}

func (m *Log) GetValue(key string) interface{} {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	if len(key) == 0 {
		return nil
	}
	return m.dataMap[key]
}

//========  public methods
func (m *Log) GetLogNodeValue(node string, key string) (data interface{}, err error) {
	if len(key) == 0 || len(node) == 0 {
		err = errors.New(fmt.Sprintf("log ginNode node or key is blank:%v", key))
		return
	}
	nodeVal := m.GetValue(node)
	if nodeVal == nil {
		err = errors.New(fmt.Sprintf("log nodeVal value is blank"))
		return
	}
	if mapData, ok := nodeVal.(map[interface{}]interface{}); ok {
		for k, v := range mapData {
			if k.(string) == key {
				data = v
				return
			}
		}
	} else {
		err = errors.New(fmt.Sprintf("GetLogNodeValue nodeVal not a map[string]string: %v", nodeVal))
	}
	return
}

func (m *Log) GetLogNodeStringValue(node string, key string) (data string, err error) {
	val, err := m.GetLogNodeValue(node, key)
	if err != nil {
		return
	}
	if v, ok := val.(string); ok {
		data = v
		return
	} else {
		err = errors.New(fmt.Sprintf("GetlogNodeStringValue value not a string:%v,%v", key, val))
	}
	return
}

func (m *Log) GetLogNodeIntValue(node string, key string) (data int, err error) {
	val, err := m.GetLogNodeValue(node, key)
	if err != nil {
		return
	}
	if v, ok := val.(int); ok {
		data = v
		return
	} else {
		err = errors.New(fmt.Sprintf("GetlogNodeIntValue value not a int:%v,%v", key, val))
	}
	return
}

func (m *Log) GetLogNodeBoolValue(node string, key string) (data bool, err error) {
	val, err := m.GetLogNodeValue(node, key)
	if err != nil {
		return
	}
	if v, ok := val.(bool); ok {
		data = v
		return
	} else {
		err = errors.New(fmt.Sprintf("GetlogNodeBoolValue value not a bool:%v,%v", key, val))
	}
	return
}
