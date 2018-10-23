package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	if err := startHttpServer(); err != nil {
		fmt.Printf("failed to start http server: %v\n", err)
		os.Exit(1)
	}
}

func startHttpServer() error {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Ping")
	})
	http.HandleFunc("/callService", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "You must specify a target url in the \"url\" query parameter", http.StatusBadRequest)
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get url: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read response body: %v", err.Error()), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(b)
	})
	return http.ListenAndServe(":8080", nil)
}
