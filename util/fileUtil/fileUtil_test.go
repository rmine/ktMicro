package fileUtil

import (
	"log"
	"testing"
)

func Test_CreateDir(t *testing.T) {
	data, err := CreateDir("./../../data/logs/gin")
	if err != nil {
		t.Log("err:", err)
	} else {
		log.Println("data:", data)
	}
}
