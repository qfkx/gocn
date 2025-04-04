package http

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

// ResultModel 通用返回值的结构体
type C返回值结构 struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

func F_提交JSON(url string, data interface{}) (string, error) {
	jsons, err1 := json.Marshal(data)
	if err1 != nil {
		return "", err1
	}
	req := &fasthttp.Request{}
	req.SetRequestURI(url)
	req.SetBody(jsons)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")
	resp := &fasthttp.Response{}
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		return "", err
	}
	b := resp.Body()
	return string(b), nil
}
