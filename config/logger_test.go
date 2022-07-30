package ktMicro

import "testing"

func Test_GenLogger(t *testing.T) {
	g1 := ginLogger
	g2 := ginLogger
	g1.Info("gin 123123")
	g2.WithField("log1", "log1value").Debug("test debug")
	MysqlLogger().Info("mysql 123123")
}
