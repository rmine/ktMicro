package ktMicro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type ResponseInfo struct {
	Gzipped bool              `json:"gzipped"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
	Form    map[string]string `json:"form"`
	Files   map[string]string `json:"files"`
	Headers map[string]string `json:"headers"`
	Cookies map[string]string `json:"cookies"`
}

func TestHTTPClient_Get(t *testing.T) {
	client := NewHttpClient(10)
	client.ReuseCookie(true)
	client.ReuseHeaders(true)

	cookie1 := &http.Cookie{Name: "cookie1", Value: "1111"}
	cookie2 := &http.Cookie{Name: "cookie2", Value: "2222"}
	cookieList := make([]*http.Cookie, 0)
	cookieList = append(cookieList, cookie1)
	client = client.AddCookieList(cookieList)
	client = client.AddCookie(cookie2)

	headerMap := make(map[string]string)
	headerMap["headerMap1"] = "headerMap1value"
	headerMap["headerMap2"] = "headerMap2value"
	client = client.AddHeader("header1", "head1value")
	client = client.AddHeaders(headerMap)

	params1 := make(map[string]string)
	params1["qq1"] = "newsget1"
	params1["qq2"] = "newsget2"
	resp, code, err := client.Get("http://httpbin.org/get", params1)
	if err != nil {
		t.Error("err", err)
	}
	var info ResponseInfo
	err = resp.ResponseToStruct(&info)
	if err != nil {
		t.Error("ResponseToStruct", err)
	}
	//data = data.(ResponseInfo)
	t.Log("code", code)
	t.Log("data", info)
	t.Log("origin", info.Url)

	//test reuse function
	resp2, code2, err2 := client.Get("http://httpbin.org/get", map[string]string{"qq2": "newsget2"})
	bytes2, err2 := client.ReadResponse(resp2)
	if err2 != nil {
		t.Log("ReadResponse", err2)
	}
	var info2 ResponseInfo
	err2 = json.Unmarshal(bytes2, &info2)
	if err2 != nil {
		fmt.Printf("Unmarshal err %v\n", err2)
	}
	t.Log("code2", code2)
	t.Log("info2", info2)
}

func TestHTTPClient_Post(t *testing.T) {
	client := DefaultHttpClient()
	resp, code, err := client.Post("http://httpbin.org/post", map[string]string{
		"p1": "newspost",
		"p2": "newspost2",
	})
	if err != nil {
		t.Error("err", err)
	}
	bytes, readerr := client.ReadResponse(resp)
	if readerr != nil {
		t.Error("readerr", readerr)
	}
	var info ResponseInfo
	merr := json.Unmarshal(bytes, &info)
	if merr != nil {
		t.Error("merr", merr)
	}
	t.Log("code", code)
	t.Log("info", info)
}

func TestHTTPClient_Patch(t *testing.T) {
	client := DefaultHttpClient()
	resp, code, err := client.Patch("http://httpbin.org/patch", map[string]string{
		"p1": "newspatch1",
		"p2": "newspatch2",
	})
	var info ResponseInfo
	err = resp.ResponseToStruct(&info)
	if err != nil {
		t.Error("ResponseToStruct", err)
	}
	//data = data.(ResponseInfo)
	t.Log("code", code)
	t.Log("data", info)
	t.Log("origin", info.Url)
}

func TestHTTPClient_Put(t *testing.T) {
	client := DefaultHttpClient()
	resp, code, err := client.Put("http://httpbin.org/put", map[string]string{
		"p1": "PUT1",
		"p2": "PUT2",
	})
	var info ResponseInfo
	err = resp.ResponseToStruct(&info)
	if err != nil {
		t.Error("ResponseToStruct", err)
	}
	//data = data.(ResponseInfo)
	t.Log("code", code)
	t.Log("data", info)
	t.Log("origin", info.Url)
}

func TestHTTPClient_Delete(t *testing.T) {
	client := DefaultHttpClient()
	resp, code, err := client.Delete("http://httpbin.org/delete", map[string]string{
		"p1": "DELETE1",
		"p2": "DELETE2",
	})
	var info ResponseInfo
	err = resp.ResponseToStruct(&info)
	if err != nil {
		t.Error("ResponseToStruct", err)
	}
	//data = data.(ResponseInfo)
	t.Log("code", code)
	t.Log("data", info)
	t.Log("origin", info.Url)
}
