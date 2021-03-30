package main

import (
	"crypto/tls"
	"fmt"
	"github.com/programschool/proxy-service/config"
	"net/http"
	"net/http/httputil"
)

var conf = config.Load()

func main() {
	go listen80()

	http.HandleFunc("/", Handle())
	address := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	fmt.Println(fmt.Sprintf("Listen: %s", address))
	_ = http.ListenAndServeTLS(address, conf.CertFile, conf.KeyFile, nil)
}

func listen80() {
	address := fmt.Sprintf("%s:%s", conf.Host, "80")
	fmt.Println(fmt.Sprintf("Listen: %s", address))
	_ = http.ListenAndServe(address, nil)
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
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		proxy.ServeHTTP(w, r)
	}
}
