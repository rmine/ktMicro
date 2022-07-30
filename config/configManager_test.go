package ktMicro

import (
	"fmt"
	"testing"
)

func Test_GetConfigData(t *testing.T) {
	val, _ := SharedConfigManager().App.GetStringValue("testModelString")
	fmt.Println("testModelString val", val)

	valcode, _ := SharedConfigManager().Code.GetMsgValue(1002)
	fmt.Println("1001 val", valcode)

	logGinFileName, _ := SharedConfigManager().Log.GetLogNodeStringValue("gin", "fileName")
	fmt.Println("logGinFileName", logGinFileName)
}
