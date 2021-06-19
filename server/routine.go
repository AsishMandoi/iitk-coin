package server

import (
	"net/http"
)

func Routine(endpoint string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc("/"+endpoint, handlerFunc)
}
