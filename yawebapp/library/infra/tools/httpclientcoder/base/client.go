package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"yawebapp/library/infra/app"
	"yawebapp/library/infra/helper"
)

type ClientConf struct {
	Endpoint string `toml:"endpoint"`
}

type AuthHandlerFunc func(*HTTPRequest) error

func DoNothingAuthHandler(r *HTTPRequest) error { return nil }

type PayloadI interface {
	Json() string
	URLEncoded() string
}

type HTTPRequest struct {
	Ctx      context.Context
	AuthFunc AuthHandlerFunc

	Method string

	Endpoint  string
	Path      string
	PathParam map[string]interface{}

	Headers map[string]string

	Body PayloadI
}

func (c HTTPRequest) BindResponse(val interface{}) error {
	// auth request, will modiy r as needed.
	err := c.AuthFunc(&c)
	if err != nil {
		return err
	}

	// make request
	req, err := c.convertRequest()
	if err != nil {
		return err
	}

	// send request
	begin := time.Now()
	app.Logger.Debug(c.Ctx, "before http request:", fmt.Sprint(req), err)
	data, err := c.sendRequest(req)
	elapsed := fmt.Sprintf("%s[%vms]%s", helper.ColorYellowBold, float64(time.Since(begin).Nanoseconds())/1e6, helper.ColorReset)
	app.Logger.Debug(c.Ctx, "after http request:", elapsed, string(data), err)
	if err != nil {
		return err
	}

	// adapt result
	data = c.adaptResult(data)

	// bind response
	return json.Unmarshal(data, val)
}

func (c HTTPRequest) adaptResult(data []byte) []byte {
	var res map[string]interface{}
	json.Unmarshal(data, &res)

	if _, ok := res["code"]; !ok {
		if v, ok := res["error_code"]; ok {
			res["code"] = v
		}
	}

	if _, ok := res["message"]; !ok {
		if v, ok := res["error_message"]; ok {
			res["message"] = v
		}
	}

	if v, ok := res["data"]; ok && v == "" {
		res["data"] = nil
	}

	j, _ := json.Marshal(res)
	return j
}

func (c HTTPRequest) sendRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%v, %v", resp.Status, string(data))
	}

	return data, nil
}

func (c HTTPRequest) convertRequest() (req *http.Request, err error) {
	reqMethod := strings.ToUpper(c.Method)

	reqHeader := make(http.Header)
	for key, val := range c.Headers {
		reqHeader.Add(key, val)
	}
	traceID := c.Ctx.Value("trace_id")
	if traceID != nil {
		reqHeader.Add("trace_id", fmt.Sprint(traceID))
	}

	var reqUrl *url.URL
	var reqPayload io.ReadCloser
	switch reqMethod {
	case "GET":
		reqUrl, err = url.Parse(fmt.Sprintf("%v/%v?%v", c.endpoint(), c.path(), c.Body.URLEncoded()))
	case "POST":
		reqUrl, err = url.Parse(fmt.Sprintf("%v/%v", c.endpoint(), c.path()))
		reqPayload = c.payload()
	}
	if err != nil {
		return nil, err
	}

	return &http.Request{
		Method: reqMethod,
		URL:    reqUrl,
		Header: reqHeader,
		Body:   reqPayload,
	}, nil
}

func (c HTTPRequest) endpoint() string {
	return strings.TrimRight(c.Endpoint, "/")
}

func (c HTTPRequest) path() string {
	p := ""
	pathArr := strings.Split(strings.Trim(c.Path, "/"), "/")
	for _, key := range pathArr {
		key := strings.Trim(key, "}{")
		if val, ok := c.PathParam[key]; ok {
			p = fmt.Sprintf("%v/%v", p, val)
		} else {
			p = fmt.Sprintf("%v/%v", p, key)
		}

	}
	return strings.Trim(p, "/")
}

func (c HTTPRequest) payload() io.ReadCloser {
	var ct string
	var ok bool
	if ct, ok = c.Headers["content-type"]; !ok {
		ct = "application/json"
	}

	var data string
	switch ct {
	case "application/json":
		data = c.Body.Json()
	default:
		data = c.Body.URLEncoded()
	}

	app.Logger.Debug(c.Ctx, "http payload:", data)
	return ioutil.NopCloser(strings.NewReader(data))
}
