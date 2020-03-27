package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	urlSuffix = "/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy"
)

var (
	proxyUrl *url.URL
	jweToken string
)

func main() {
	var err error

	initViper()

	refreshJwe()
	go func() {
		t := time.NewTicker(10*time.Minute)
		defer t.Stop()
		for range t.C {
			refreshJwe()
		}
	}()

	proxyUrlString := viper.GetString("Proxy") + urlSuffix + "/"
	proxyUrl, err = url.Parse(proxyUrlString)
	if err != nil {
		log.Fatalf("Failed to parse proxy url: %s\nError: %v\n", proxyUrlString, err)
	}

	listen := viper.GetString("Listen")
	log.Println("Serving on", listen)
	http.HandleFunc("/", handleRequest)
	if err := http.ListenAndServe(listen, nil); err != nil {
		log.Fatalf("Error listening and serving\nError: %v\n", err)
	}
}
