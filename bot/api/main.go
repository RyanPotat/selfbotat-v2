package api

import (
	"bytes"
	"net/http"
	"time"
)

var Make = &http.Client{
	Timeout: time.Second * 10,
}

func MakeRequest(
	method string, 
	url string, 
	headers map[string]string,
	body *bytes.Buffer,
) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	

	res, err := Make.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
