package ktMicro

import (
	"github.com/rmine/ktMicro/util/fileUtil"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	ginLogger   = LoggerFatory("gin")
	mysqlLogger = LoggerFatory("mysql")
	redisLogger = LoggerFatory("redis")
	mongoLogger = LoggerFatory("mongo")
	mqLogger    = LoggerFatory("mq")
)

func Newlogger() *logrus.Logger {
	m := logrus.New()
	return m
}

func LoggerFatory(node string) *logrus.Logger {
	m := Newlogger()
	loggerFilePath, err := LogFilePath(node)
	if err != nil {
		log.Fatalln("LoggerFatory get LogFilePath error:", node, err)
	}
	var absPath string
	if strings.HasPrefix(loggerFilePath, "/") || strings.HasPrefix(loggerFilePath, "~/") {
		absPath = loggerFilePath
	} else {
		absPath, err = filepath.Abs(loggerFilePath)
		if err != nil {
			log.Fatalln("LoggerFatory absPath error:", absPath, err)
		}
	}

	createSucc, err := fileUtil.CreateDir(absPath)
	if err != nil || !createSucc {
		log.Fatalln("LoggerFatory CreateDir error:", absPath, err)
	}

	loggerFileName, err := LogFileName(node)
	if err != nil {
		log.Fatalln("LoggerFatory get loggerFileName error:", node, err)
	}
	// 设置日志级别
	m.SetLevel(logrus.DebugLevel)

	/* 单个文件日志输出
	// 写入文件
	src, err := fileUtil.CreateFile(absPath, loggerFileName)
	if err != nil {
		log.Fatalln("LoggerFatory OpenFile error:", loggerFilePath, loggerFileName, err)
	}
	// 设置输出
	m.Out = src
	*/

	if AppIsDev() {
		m.Out = os.Stdout
	} else {
		//切割日志
		lumberjackFileName := path.Join(absPath, loggerFileName)
		var lumberjackMaxSize int
		lumberjackMaxSize, err = LogIntValue(node, "MaxSize")
		if err != nil {
			lumberjackMaxSize = 20
		}
		var lumberjackMaxBackups int
		lumberjackMaxBackups, err = LogIntValue(node, "MaxBackups")
		if err != nil {
			lumberjackMaxBackups = 20
		}
		var lumberjackMaxAge int
		lumberjackMaxAge, err = LogIntValue(node, "MaxAge")
		if err != nil {
			lumberjackMaxAge = 20
		}
		var lumberjackCompress bool
		lumberjackCompress, err = LogBoolValue(node, "Compress")
		if err != nil {
			lumberjackCompress = true
		}

		m.Out = &lumberjack.Logger{
			Filename:   lumberjackFileName,
			MaxSize:    lumberjackMaxSize, // megabytes
			MaxBackups: lumberjackMaxBackups,
			MaxAge:     lumberjackMaxAge,   //days
			Compress:   lumberjackCompress, // disabled by default
		}
	}
	return m
}

func GinLogger() *logrus.Logger {
	return ginLogger
}

func MysqlLogger() *logrus.Logger {
	return mysqlLogger
}

func RedisLogger() *logrus.Logger {
	return redisLogger
}

func MongoLogger() *logrus.Logger {
	return mongoLogger
}

func MqLogger() *logrus.Logger {
	return mqLogger
}
