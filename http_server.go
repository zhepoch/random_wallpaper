package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func ListenHttpServer() error {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/change_query_key", ChangeQueryKeyHandler)

	addr := fmt.Sprintf("127.0.0.1:%d", *ListenPort)
	return http.ListenAndServe(addr, nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("only support method: GET"))
		return
	}

	log.Infoln("Get: /")
	reader, err := GetChangeIndex()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, reader)
	return
}

func ChangeQueryKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("only support method: GET"))
		return
	}

	value := r.FormValue("key")
	log.Infoln("POST: /change_query_key,", r.Form.Encode())
	//if value != "" {
	//	encodeKey := url.QueryEscape(value)
	//	*PhotoQueryKey = encodeKey
	//}

	encodeKey := url.QueryEscape(value)
	*PhotoQueryKey = encodeKey

	UserChange <- struct{}{}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Update Success: "))
	_, _ = w.Write([]byte(*PhotoQueryKey))
	_, _ = w.Write([]byte(`<br/> <a href="/">go back</a>`))
	return
}
