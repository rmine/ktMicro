package ktMicro

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
)

type Code struct {
	yamlFile *ymlFile
	dataMap  map[int]string
	rwMutex  *sync.RWMutex
}

func NewCode() *Code {
	m := &Code{}
	m.rwMutex = new(sync.RWMutex)
	return m
}

//==========protocol
func (m *Code) Load() {
	m.yamlFile = NewYamlFile("configFiles/code.yml")
	err := m.yamlFile.Load(&m.dataMap)
	if err != nil {
		log.Fatalln("Code Load err:", err)
		return
	}
}

func (m *Code) UpdateConfig() {
}

func (m *Code) GetValue(key string) (data interface{}, err error) {
	intKey, err := strconv.Atoi(key)
	if err != nil {
		return
	}
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()
	if len(key) == 0 {
		err = errors.New("code GetValue key can not be blank!")
		return
	}
	data = m.dataMap[intKey]
	return
}

func (m *Code) GetMsgValue(code int) (data string, err error) {
	val, err := m.GetValue(strconv.Itoa(code))
	if err != nil {
		return
	}
	if v, ok := val.(string); ok {
		data = v
		return
	} else {
		err = errors.New(fmt.Sprintf("code GetStringValue value not a string:%v,%v", code, val))
	}
	return
}
