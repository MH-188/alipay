package alipay

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func HttpDoRequest(method, url string, data []byte) ([]byte, error) {
	newRequest, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	newRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	client := http.DefaultClient
	resp, err := client.Do(newRequest)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
