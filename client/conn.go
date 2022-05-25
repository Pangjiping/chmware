package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	METHOD_GET    = "get"
	METHOD_POST   = "post"
	METHOD_DELETE = "delete"
)

type connRequest struct {
	address string
	method  string
	key     string
	value   string
}

type connResponse struct {
	address   string
	method    string
	key       string
	value     string
	timestamp time.Time // 本次操作的时间
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// todo 连接接口 http
func doGET(req connRequest) (string, error) {
	url, err := buildURL(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func doPOST(req connRequest) error {
	url, err := buildURL(req)
	if err != nil {
		return nil
	}

	b, err := json.Marshal(map[string]string{req.key: req.value})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func doDELETE(req connRequest) error {
	ru, err := url.Parse(fmt.Sprintf("%s/key/%s", req.address, req.key))
	if err != nil {
		return err
	}

	request := &http.Request{
		Method: "DELETE",
		URL:    ru,
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// todo 连接池

// todo 连接配置信息

// 构建请求链接
func buildURL(req connRequest) (string, error) {
	method := req.method
	switch method {
	case METHOD_GET:
		return fmt.Sprintf(`http://%s/key/%s`, req.address, req.key), nil
	case METHOD_POST:
		return fmt.Sprintf(`http://%s/key`, req.address), nil
	case METHOD_DELETE:
		return fmt.Sprintf(`http://%s/key/%s`, req.address, req.key), nil
	default:
		return "", fmt.Errorf("invalid request method: %s", req.method)
	}
}
