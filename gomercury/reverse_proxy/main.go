package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	url, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	err = http.ListenAndServe(":3001", proxy)
	if err != nil {
		panic(err)
	}
}
