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
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return &http.Request{}, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	return req, nil
}

func BuildPOSTRequest(url, body string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return &http.Request{}, err
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
