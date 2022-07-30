package ktMicro

import (
	"sync"
)

type ConfigManager struct {
	Config       *Config
	App          *App
	Code         *Code
	Log          *Log
	ConfigSwitch *ConfigSwitch
}

var (
	configManagerOnce          sync.Once
	shareInstanceConfigManager *ConfigManager
)

func SharedConfigManager() *ConfigManager {
	configManagerOnce.Do(func() {
		shareInstanceConfigManager = &ConfigManager{}
		initAllModel(shareInstanceConfigManager)
	})
	return shareInstanceConfigManager
}

func init() {
	SharedConfigManager()
}

func initAllModel(shareInstance *ConfigManager) {
	shareInstance.Config = NewConfig()
	shareInstance.Config.Load()

	appEnv, err := shareInstance.Config.GetAppConfigStringValue("AppEnv")
	if err != nil {
		//保底设置开发环境
		appEnv = "develop"
	}

	shareInstance.App = NewApp()
	shareInstance.App.AppEnv = appEnv
	shareInstance.App.Load()

	shareInstance.Code = NewCode()
	shareInstance.Code.Load()

	shareInstance.Log = NewLog()
	shareInstance.Log.AppEnv = appEnv
	shareInstance.Log.Load()

	shareInstance.ConfigSwitch = NewConfigSwitch()
	shareInstance.ConfigSwitch.Load()
}

func Load() {

}

//===============快捷方法 start ===============
//app环境
func AppEnv() (data string) {
	return SharedConfigManager().App.AppEnv
}

func AppIsDev() (data bool) {
	return SharedConfigManager().Config.AppIsDev()
}

//code错误码
func GetMsgValue(code int) (data string) {
	data, err := SharedConfigManager().Code.GetMsgValue(code)
	if err != nil {
		data = "服务器繁忙，请稍后再试"
		return
	}
	return
}

//serverport
func AppServerPort() (data string) {
	return SharedConfigManager().Config.AppServerPort()
}

func LogFilePath(node string) (data string, err error) {
	return SharedConfigManager().Log.GetLogNodeStringValue(node, "filePath")
}

func LogFileName(node string) (data string, err error) {
	return SharedConfigManager().Log.GetLogNodeStringValue(node, "fileName")
}

func LogIntValue(node string, key string) (data int, err error) {
	return SharedConfigManager().Log.GetLogNodeIntValue(node, key)
}

func LogBoolValue(node string, key string) (data bool, err error) {
	return SharedConfigManager().Log.GetLogNodeBoolValue(node, key)
}

func GetConfigSwitch(key string) (data bool, err error) {
	return SharedConfigManager().ConfigSwitch.GetConfigSwitch(key)
}

//===============快捷方法 end ===============
