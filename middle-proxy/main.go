package main

import (
	"crypto/tls"
	"fmt"
	"github.com/programschool/proxy-service/config"
	"net/http"
	"net/http/httputil"
	"time"
)

var conf = config.Load()

func main() {
	go listen80()

	http.HandleFunc("/", Handle())
	address := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	fmt.Println(fmt.Sprintf("Listen: %s", address))
	server := &http.Server{
		Addr: address,
		//ReadHeaderTimeout: 1 * time.Minute,
		//IdleTimeout:       1 * time.Minute,
		//ReadTimeout:       1 * time.Minute,
	}
	server.SetKeepAlivesEnabled(false)
	//_ = server.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
	_ = server.ListenAndServe()
}

func listen80() {
	address := fmt.Sprintf("%s:%s", conf.Host, "8000")
	fmt.Println(fmt.Sprintf("Listen: %s", address))
	server := &http.Server{
		Addr: address,
		//ReadHeaderTimeout: 1 * time.Minute,
		//IdleTimeout:       1 * time.Minute,
		//ReadTimeout:       1 * time.Minute,
	}
	server.SetKeepAlivesEnabled(false)
	_ = server.ListenAndServe()
}

func Handle() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")

		//fmt.Println("container")
		//fmt.Println(r.Header.Get("container"))

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.URL.Scheme = "http"
			req.URL.Host = r.Header.Get("container")
		}

		proxy := &httputil.ReverseProxy{
			Director: director,
			Transport: &http.Transport{
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				IdleConnTimeout:       30 * time.Second,
				MaxIdleConnsPerHost:   32, // seems about optimal, see #2805
				ResponseHeaderTimeout: 2 * time.Minute,
				DisableKeepAlives:     true,
			},
		}

		proxy.ServeHTTP(w, r)
	}
}
