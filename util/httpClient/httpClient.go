package ktMicro

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HttpClient struct {
	Client *http.Client
	//private fields
	reqHeaders   map[string]string
	reuseHeaders bool
	reqCookies   []*http.Cookie
	reuseCookie  bool
}

type HttpResponse struct {
	*http.Response
}

func DefaultHttpClient() *HttpClient {
	return NewHttpClient(5)
}

//单位秒
func NewHttpClient(timeout time.Duration) *HttpClient {
	client := &http.Client{}
	client.Timeout = timeout * time.Second
	c := &HttpClient{Client: client}
	c.reqHeaders = make(map[string]string)
	c.reqCookies = make([]*http.Cookie, 0)
	return c
}

func (c *HttpClient) ReuseCookie(reuse bool) *HttpClient {
	c.reuseCookie = reuse
	return c
}

func (c *HttpClient) ReuseHeaders(reuse bool) *HttpClient {
	c.reuseHeaders = reuse
	return c
}

func (c *HttpClient) AddHeader(k string, v string) *HttpClient {
	if k != "" && v != "" {
		c.reqHeaders[k] = v
	}
	return c
}

func (c *HttpClient) AddHeaders(mapData map[string]string) *HttpClient {
	for k, v := range mapData {
		if k != "" && v != "" {
			c.reqHeaders[k] = v
		}
	}
	return c
}

func (c *HttpClient) AddCookie(cookies ...*http.Cookie) *HttpClient {
	if cookies != nil {
		c.reqCookies = append(c.reqCookies, cookies...)
	}
	return c
}

func (c *HttpClient) AddCookieList(cookieList []*http.Cookie) *HttpClient {
	for _, v := range cookieList {
		c.reqCookies = append(c.reqCookies, v)
	}
	return c
}

func (c *HttpClient) Do(method string, url string, headers map[string]string,
	body io.Reader) (response *HttpResponse, code int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	headers = mergeHeaders(headers, c.reqHeaders)
	cookies := c.reqCookies

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return serverError(err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	if !c.reuseHeaders && len(c.reqHeaders) > 0 {
		c.reqHeaders = nil
	}
	if !c.reuseCookie && len(c.reqCookies) > 0 {
		c.reqCookies = nil
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return serverError(err)
	} else {
		if resp.StatusCode != http.StatusOK {
			return serverError(errors.New(resp.Status))
		}
	}
	statusCode := resp.StatusCode
	return &HttpResponse{resp}, statusCode, err
}

func (m *HttpResponse) ResponseToMap() (data map[string]interface{}, err error) {
	bytesData, err := getRespBytes(m)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		return
	}
	return
}

func (m *HttpResponse) ResponseToStruct(model interface{}) (err error) {
	bytesData, err := getRespBytes(m)
	if err != nil {
		return
	}
	mapResponse := make(map[string]interface{})
	err = json.Unmarshal(bytesData, &mapResponse)
	if err != nil {
		return
	}
	err = mapstructure.Decode(mapResponse, model)
	if err != nil {
		return
	}
	return nil
}

func (m *HttpResponse) ResponseToString() (data string, err error) {
	bytesData, err := getRespBytes(m)
	if err != nil {
		return
	}
	return string(bytesData), nil
}

func (c *HttpClient) ReadResponse(resp *HttpResponse) ([]byte, error) {
	return getRespBytes(resp)
}

func getRespBytes(m *HttpResponse) (bytes []byte, err error) {
	var reader io.ReadCloser
	if m.Header != nil {
		switch m.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err = gzip.NewReader(m.Body)
			if err != nil {
				return
			}
		default:
			reader = m.Body
		}
	} else {
		reader = m.Body
	}

	defer reader.Close()

	bytes, err = ioutil.ReadAll(reader)
	if err != nil {
		//log
		return
	}
	return
}

//########## Get
func (c *HttpClient) Get(url string, params ...interface{}) (*HttpResponse, int, error) {
	for _, p := range params {
		url = addParams(url, convertToUrlValues(p))
	}
	return c.Do(http.MethodGet, url, nil, nil)
}

//########## Post
func (c *HttpClient) Post(url string, params interface{}) (respose *HttpResponse, code int, err error) {
	paramsValues := convertToUrlValues(params)
	// 文件上传
	if checkParamFile(paramsValues) {
		return c.PostMultipart(url, params)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	body := strings.NewReader(paramsValues.Encode())

	return c.Do(http.MethodPost, url, headers, body)
}

func (c *HttpClient) PostMultipart(url string, params interface{}) (response *HttpResponse, code int, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	paramsValues := convertToUrlValues(params)

	for k, v := range paramsValues {
		for _, pv := range v {
			//文件类型以@起始
			if k[0] == '@' {
				err := addFormFile(writer, k[1:], pv)
				if err != nil {
					return serverError(err)
				}
			} else {
				err = writer.WriteField(k, pv)
				return
			}
		}
	}
	headers := make(map[string]string)

	headers["Content-Type"] = writer.FormDataContentType()
	err = writer.Close()
	if err != nil {
		return serverError(err)
	}
	return c.Do(http.MethodPost, url, headers, body)
}

func (c *HttpClient) PostMultipartWithByte(url string, params interface{}) (
	[]byte, int, error) {
	resp, _, err := c.PostMultipart(url, params)
	return c.readResponseToByte(resp, err)
}

func (c *HttpClient) PostJson(url string, data interface{}) (*HttpResponse, int, error) {
	return c.requestJson(http.MethodPost, url, data)
}

func (c *HttpClient) requestJson(method string, url string, data interface{}) (*HttpResponse, int, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	var body []byte
	switch datatype := data.(type) {
	case []byte:
		body = datatype
	case string:
		body = []byte(datatype)
	default:
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}
	return c.Do(method, url, headers, bytes.NewReader(body))
}

//########## Delete
func (c *HttpClient) Delete(url string, params ...interface{}) (*HttpResponse, int, error) {
	for _, p := range params {
		url = addParams(url, convertToUrlValues(p))
	}
	return c.Do(http.MethodDelete, url, nil, nil)
}

//########## Head
func (c *HttpClient) Head(url string) (*HttpResponse, int, error) {
	return c.Do(http.MethodHead, url, nil, nil)
}

//########## Put
func (c *HttpClient) Put(url string, params ...interface{}) (*HttpResponse, int, error) {
	for _, p := range params {
		url = addParams(url, convertToUrlValues(p))
	}
	return c.Do(http.MethodPut, url, nil, nil)
}

//func (c *HttpClient) Put(url string, body io.Reader) (*HttpResponse, int, error) {
//	return c.Do(http.MethodPut, url, nil, body)
//}

//########## Patch
func (c *HttpClient) Patch(url string, params ...map[string]string) (*HttpResponse, int, error) {
	for _, p := range params {
		url = addParams(url, convertToUrlValues(p))
	}
	return c.Do(http.MethodPatch, url, nil, nil)
}

//######################################################################################################################
//######################################## privated methods ############################################################
//######################################################################################################################
func addParams(_url string, params url.Values) string {
	if len(params) == 0 {
		return _url
	}

	if !strings.Contains(_url, "?") {
		_url += "?"
	}

	if strings.HasSuffix(_url, "?") || strings.HasSuffix(_url, "&") {
		_url += params.Encode()
	} else {
		_url += "&" + params.Encode()
	}

	return _url
}

func checkParamFile(params url.Values) bool {
	for k := range params {
		if k[0] == '@' {
			return true
		}
	}
	return false
}

func addFormFile(writer *multipart.Writer, name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
}

//支持多个header合并, 越后面优先级越高
func mergeHeaders(headers ...map[string]string) map[string]string {
	res := make(map[string]string)
	for _, singleHeader := range headers {
		for k, v := range singleHeader {
			res[k] = v
		}
	}
	return res
}

//参数标准化,最终类型url.Values
func convertToUrlValues(v interface{}) url.Values {
	switch vtype := v.(type) {
	case url.Values:
		return vtype
	case map[string][]string:
		return url.Values(vtype)
	case map[string]string:
		res := make(url.Values)
		for k, v := range vtype {
			res.Add(k, v)
		}
		return res
	case []interface{}:
		res := make(url.Values)
		for _, v2 := range v.([]interface{}) {
			res = convertToUrlValues(v2)
		}
		return res
	case nil:
		return make(url.Values)
	default:
		panic("Invalid Value")
	}
}

//将http.response转化成[]byte类型
func (c *HttpClient) readResponseToByte(resp *HttpResponse, err error) ([]byte, int, error) {
	if err != nil {
		return serverByteError(err)
	}
	bytesData, err := c.ReadResponse(resp)
	if err != nil {
		return serverByteError(err)
	}
	return bytesData, resp.StatusCode, err
}

//统一500错误抛出
func serverError(perr error) (response *HttpResponse, code int, err error) {
	return nil, http.StatusInternalServerError, perr
}

func serverByteError(err error) ([]byte, int, error) {
	return nil, http.StatusInternalServerError, err
}
