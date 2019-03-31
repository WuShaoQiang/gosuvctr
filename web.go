package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func postFormWithAuth(url string, data url.Values, key string) (r JSONResponse, err error) {
	resp, err := postWithAuth(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), key)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return r, fmt.Errorf("POST %v %v", strconv.Quote(url), string(body))
	}
	return r, nil
}

func postWithAuth(url, contentType string, body io.Reader, key string) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", key)
	return http.DefaultClient.Do(req)
}

func getWithAuth(url string, key string) (resp *http.Response, err error) {
	// url := remoteURL + pathname
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", key)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
