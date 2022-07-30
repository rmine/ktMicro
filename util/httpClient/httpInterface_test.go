package ktMicro

import "testing"

type ImIWantData struct {
	ServeTrack ServeTrack `mapstructure:"serve_track" json:"serve_track"`
}

type ServeTrack struct {
	Image          string `json:"image"`
	Title          string `json:"title"`
	TitleFontcolor string `json:"title_fontcolor"`
	Title_BG       string `json:"title_BG"`
	Link           string `json:"link"`
	Id             string `mapstructure:"_id" json:"_id"`
}

func TestInterface_Get(t *testing.T) {

	params := make(map[string]string)
	params["param1"] = "p1"
	params["param2"] = "p2"

	httpInterface := NewHttpInterface()
	httpInterface.url = "http://httpbin.org/get"
	httpInterface.SuccCode = 0
	//httpInterface.MsgString = "message"
	resp, msg, err := httpInterface.Get(params)
	if err != nil {
		t.Error("err", err)
		return
	} else {
		t.Log("resp", resp)
	}
	t.Log("msg", msg)
	var imIWantData ImIWantData
	err = ParseData(resp, &imIWantData)
	if err != nil {
		t.Error("ParseData", err)
		return
	} else {
		t.Log("Data", imIWantData.ServeTrack)
		t.Log("Data id", imIWantData.ServeTrack.Id)
	}
}
