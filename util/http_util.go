package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
	log.Printf("http error: %s", err)
}

func BuildGETRequest(url string, headers map[string]string) (*http.Request, error) {
	fmt.Printf("building get request\n url:%s, headers count:%v\n", url, len(headers))
	req, err := BuildRequest("GET", url, "", headers)
	if err != nil {
		return &http.Request{}, err
	}
	return req, nil
}

func BuildDELETERequest(url string, headers map[string]string) (*http.Request, error) {
	fmt.Printf("building delete request\n url:%s, headers count:%v\n", url, len(headers))
	req, err := BuildRequest("DELETE", url, "", headers)
	if err != nil {
		return &http.Request{}, err
	}
	return req, nil
}

func BuildPOSTRequest(url, body string, headers map[string]string) (*http.Request, error) {
	fmt.Printf("building post request\n url:%s, headers count:%v\n", url, len(headers))
	req, err := BuildRequest("POST", url, body, headers)
	if err != nil {
		return &http.Request{}, err
	}
	return req, nil
}

func BuildPUTRequest(url, body string, headers map[string]string) (*http.Request, error) {
	fmt.Printf("building put request\n url:%s, headers count:%v\n", url, len(headers))
	req, err := BuildRequest("PUT", url, body, headers)
	if err != nil {
		return &http.Request{}, err
	}
	return req, nil
}

func BuildRequest(method, url, body string, headers map[string]string) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)
	if body != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
		headers["Content-Type"] = "application/json"
	}
	if err != nil {
		return req, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	return req, nil
}

func ReadResponse(r *http.Response) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body), err
}

func ReadRequestBody(r *http.Request) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body), err
}

func LogRequest(r *http.Request) {
	log.Printf("http request: from:%s, url:%s\nheaders:%v\n", r.RemoteAddr, r.URL, r.Header)
}

func LogResponse(status int, msg string) {
	log.Printf("http response: %v, msg: %s", status, msg)
}
