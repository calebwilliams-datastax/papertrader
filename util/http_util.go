package util

import "net/http"

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
