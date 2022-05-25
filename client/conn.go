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

type Response struct {
	key       string
	value     string
	err       error
	method    string
	createdAt time.Time
}

func (r Response) Error() error {
	return r.err
}

func (r Response) Value() string {
	return r.value
}

func (r Response) Key() string {
	return r.key
}

func (r Response) Timestamp() time.Time {
	return r.createdAt
}

func (r Response) Method() string {
	return r.method
}

// 连接接口 http
// todo: 处理这个返回信息的格式，只保留value值
func doGET(req connRequest) (string, error) {
	url := fmt.Sprintf(`http://%s/key/%s`, req.address, req.key)

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
	url := fmt.Sprintf(`http://%s/key`, req.address)

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
// func buildURL(req connRequest) (string, error) {
// 	method := req.method
// 	switch method {
// 	case METHOD_GET:
// 		return fmt.Sprintf(`http://%s/key/%s`, req.address, req.key), nil
// 	case METHOD_POST:
// 		return fmt.Sprintf(`http://%s/key`, req.address), nil
// 	case METHOD_DELETE:
// 		return fmt.Sprintf(`http://%s/key/%s`, req.address, req.key), nil
// 	default:
// 		return "", fmt.Errorf("invalid request method: %s", req.method)
// 	}
// }

func invoke(req connRequest) Response {
	method := req.method
	switch method {
	case METHOD_GET:
		value, err := doGET(req)
		return Response{
			key:   req.key,
			value: value,
			err:   err,
		}
	case METHOD_POST:
		err := doPOST(req)
		return Response{
			key:   req.key,
			value: req.value,
			err:   err,
		}
	case METHOD_DELETE:
		err := doDELETE(req)
		return Response{
			key:   req.key,
			value: req.value,
			err:   err,
		}
	default:
		return Response{
			key:   "",
			value: "",
			err:   fmt.Errorf("invalid request method: $s", req.method),
		}
	}
}
