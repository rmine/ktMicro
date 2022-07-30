package ktMicro

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"strings"
	"time"
)

type HttpInterface struct {
	Client   *HttpClient
	Response *HttpResponse

	CodeString string //映射json中code节点
	MsgString  string //映射json中msg节点
	DataString string //映射json中data节点

	SuccCode int           //请求成功对应的成功code
	url      string        //接口url
	method   string        //get,post,patch,put,delete,随意大小写
	timeOut  time.Duration //超时时间,单位秒

}

type HttpInterfaceResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	codeStringList = []string{"code", "errno", "errcode", "errCode", "errorCode", "resultCode", "returncode", "result", "retCode"}
	msgStringList  = []string{"message", "errmsg", "errMsg", "msg", "errormsg", "retMsg"}
	dataStringList = []string{"data", "content"}
)

func NewHttpInterface() *HttpInterface {
	m := &HttpInterface{}
	m.SuccCode = 0
	m.method = "get"
	m.timeOut = 10
	m.Client = NewHttpClient(m.timeOut)
	return m
}

func (m *HttpInterface) SetUrl(url string) {
	m.url = url
}

func (m *HttpInterface) SetTimeOut(timeOut time.Duration) {
	m.timeOut = timeOut
}

func (m *HttpInterface) Method(method string) {
	if method != "" && len(method) > 0 {
		m.method = strings.ToUpper(method)
	}
}

func (m *HttpInterface) Get(params ...interface{}) (data interface{}, msg string, err error) {
	url := m.url
	if url == "" {
		err = errors.New("[url] must not be blank!")
		return
	}
	resp, _, err := m.Client.Get(url, params)
	if err != nil {
		return
	}
	return m.getDataFromResponse(resp)
}

func (m *HttpInterface) Post(params ...interface{}) (data interface{}, msg string, err error) {
	url := m.url
	if url == "" {
		err = errors.New("[url] must not be blank!")
		return
	}
	resp, _, err := m.Client.Post(url, params)
	if err != nil {
		return
	}
	return m.getDataFromResponse(resp)
}

func (m *HttpInterface) PostJson(params ...interface{}) (data interface{}, msg string, err error) {
	url := m.url
	if url == "" {
		err = errors.New("[url] must not be blank!")
		return
	}
	resp, _, err := m.Client.PostJson(url, params)
	if err != nil {
		return
	}
	return m.getDataFromResponse(resp)
}

func (m *HttpInterface) DoRequest(method string, params ...interface{}) (data interface{}, msg string, err error) {
	url := m.url
	if url == "" {
		err = errors.New("[url] must not be blank!")
		return
	}
	url, err = m.getMethod(method)
	if err != nil {
		return
	}
	for _, p := range params {
		url = addParams(url, convertToUrlValues(p))
	}
	resp, _, err := m.Client.Do(method, url, nil, nil)
	if err != nil {
		return
	}
	return m.getDataFromResponse(resp)
}

func (m *HttpInterface) SetSuccCode(code int) {
	m.SuccCode = code
}

func (m *HttpInterface) SetCodeString(codeString string) {
	if len(codeString) > 0 {
		m.CodeString = codeString
	}
}

func (m *HttpInterface) SetMsgString(msgString string) {
	if len(msgString) > 0 {
		m.MsgString = msgString
	}
}

func (m *HttpInterface) SetDataString(dataString string) {
	if len(dataString) > 0 {
		m.DataString = dataString
	}
}

//解析data节点结构数据
func ParseData(jsonData interface{}, value interface{}) error {
	if decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: value, WeaklyTypedInput: true}); err == nil {
		return decoder.Decode(jsonData)
	} else {
		return err
	}
}

//==========private methods
func (m *HttpInterface) getMethod(method string) (data string, err error) {
	data = method
	if data == "" {
		data = m.url
	}
	if data == "" {
		err = errors.New("[method] must not be blank!")
		return
	}
	return strings.ToUpper(data), nil
}

func (m *HttpInterface) getDataFromResponse(resp *HttpResponse) (data interface{}, msg string, err error) {
	mapData, err := resp.ResponseToMap()
	if err != nil {
		return
	}
	codeString := getHitCodeString(m.CodeString, mapData)
	if len(codeString) == 0 {
		err = errors.New("Response struct [codeString] not correct!")
		return
	}

	msgString := getHitMsgString(m.MsgString, mapData)
	if len(msgString) == 0 {
		err = errors.New("Response struct [msgString] not correct!")
		return
	}

	respCode := int(mapData[codeString].(float64))
	//成功
	if respCode == m.SuccCode {
		dataString := getHitDataString(m.DataString, mapData)
		if len(dataString) == 0 {
			data = nil
		} else {
			data = mapData[dataString]
		}
		if v, ok := mapData[msgString].(string); ok {
			msg = v
		}
	} else {
		//失败
		if v, ok := mapData[msgString].(string); ok {
			err = errors.New(v)
			msg = v
		} else {
			err = errors.New("HttpInterface failed!")
		}
	}
	return
}

func getHitCodeString(codeString string, mapData map[string]interface{}) (data string) {
	if len(codeString) > 0 {
		return codeString
	}
	for _, v := range codeStringList {
		if mapData[v] != nil {
			return v
		}
	}
	return
}

func getHitMsgString(msgString string, mapData map[string]interface{}) (data string) {
	if len(msgString) > 0 {
		return msgString
	}
	for _, v := range msgStringList {
		if mapData[v] != nil {
			return v
		}
	}
	return
}

func getHitDataString(dataString string, mapData map[string]interface{}) (data string) {
	if len(dataString) > 0 {
		return dataString
	}
	for _, v := range dataStringList {
		if mapData[v] != nil {
			return v
		}
	}
	return
}
