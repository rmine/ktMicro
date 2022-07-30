package ktMicro

import "testing"

func Test_GetConfigSwitchData(t *testing.T) {
	data, err := GetConfigSwitch("mysql")
	if err != nil {
		t.Error("err", err)
	} else {
		t.Log("data", data)
	}
}
