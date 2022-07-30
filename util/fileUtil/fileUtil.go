package fileUtil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(filePath string) (b bool, err error) {
	exist, err := PathExists(filePath)
	if err != nil {
		return
	}
	if exist {
		b = true
		return
	} else {
		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return
		}
		err = os.Chmod(filePath, os.ModePerm)
		if err != nil {
			return
		}
		b = true
	}
	return
}

func CreateFile(filePath string, fileName string) (file *os.File, err error) {
	if filePath == "" || fileName == "" {
		err = errors.New("FileUtil filePath or filename is blank!")
		return
	}
	cdb, err := CreateDir(filePath)
	if err != nil {
		return
	}
	if !cdb {
		err = errors.New(fmt.Sprintf("CreateDir error:%v,%v", filePath, fileName))
		return
	}
	fullfile := path.Join(filePath, fileName)
	// 写入文件
	resultfile, err := os.OpenFile(fullfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln("FileUtil CreateFile error:", err)
	}
	file = resultfile
	return
}

func ReadFile(path string) (data []byte, err error) {
	var absPath string
	absPath, err = filepath.Abs(path)
	if err != nil {
		return
	} else {
		data, err = ioutil.ReadFile(absPath)
		if err != nil {
			return
		}
	}
	return
}
