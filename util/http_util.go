package util

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func BuildGETRequest(url string, headers map[string]string) (*http.Request, error) {
	req, err := BuildRequest("GET", url, "", headers)
	if err != nil {
		return &http.Request{}, err
	}
	return req, nil
}

func BuildDELETERequest(url string, headers map[string]string) (*http.Request, error) {
	req, err := BuildRequest("DELETE", url, "", headers)
	if err != nil {
		return &http.Request{}, err
	}
	return req, nil
}

func BuildPOSTRequest(url, body string, headers map[string]string) (*http.Request, error) {
	req, err := BuildRequest("POST", url, body, headers)
	if err != nil {
		return &http.Request{}, err
	}
	return req, nil
}

func BuildPUTRequest(url, body string, headers map[string]string) (*http.Request, error) {
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
