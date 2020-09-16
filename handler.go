package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func handleRequest(res http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

	req.Header.Del("Accept-Encoding")
	req.Header.Set("jweToken", jweToken)
	req.AddCookie(&http.Cookie{
		Name:       "jweToken",
		Value:      url.QueryEscape(jweToken),
	})
	req.Host = proxyUrl.Host

	proxy.ServeHTTP(res, req)
}

