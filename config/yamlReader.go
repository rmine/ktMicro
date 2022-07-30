package ktMicro

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rmine/ktMicro/util/fileUtil"
	"gopkg.in/yaml.v2"
	"log"
	"sync"
)

type YamlFile interface {
	Load()
	UpdateConfig()
	GetValue(key string) (data interface{}, err error)
}

type ymlFile struct {
	path string
}

func NewYamlFile(path string) *ymlFile {
	m := &ymlFile{path: path}
	return m
}

func (m *ymlFile) Load(model interface{}) (err error) {
	var yamlfile []byte
	yamlfile, err = fileUtil.ReadFile(m.path)
	if err != nil {
		yamlfile, err = fileUtil.ReadFile("../" + m.path)
		if err != nil {
			return
		}
	}

	err = yaml.Unmarshal(yamlfile, model)
	if err != nil {
		return
	}
	return
}

func (m *ymlFile) UpdateConfig(*sync.RWMutex) (err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/tmp/foo")
	if err != nil {
		log.Fatal(err)
	}
	<-done

	return err
}
