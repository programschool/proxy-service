package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/programschool/proxy-service/config"
	"net"
	"net/http"
	"net/http/httputil"
)

var conf = config.Load()

func main() {
	go listen80()

	http.HandleFunc("/", Handle())
	address := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	fmt.Println(fmt.Sprintf("Listen: %s", address))
	server := &http.Server{Addr: address}
	server.SetKeepAlivesEnabled(false)
	_ = server.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
}

func listen80() {
	address := fmt.Sprintf("%s:%s", conf.Host, "80")
	fmt.Println(fmt.Sprintf("Listen: %s", address))
	server := &http.Server{Addr: address}
	server.SetKeepAlivesEnabled(false)
	_ = server.ListenAndServe()
}

func Handle() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		serverPort := "443"
		scheme := "https"
		domain, port, err := net.SplitHostPort(r.Host)
		if err != nil {
			serverPort = "80"
			scheme = "http"
			port = "80"
			domain = r.Host
		}

		//fmt.Println("scheme")
		//fmt.Println(scheme)
		//fmt.Println("serverPort")
		//fmt.Println(serverPort)
		//fmt.Println("port")
		//fmt.Println(port)
		//fmt.Println(domain)

		info := getFromRedis(domain)
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")

		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Set("container", fmt.Sprintf("%s:%s", info.containerIp, port))
			req.URL.Scheme = scheme
			req.URL.Host = fmt.Sprintf("%s:%s", info.dockerServer, serverPort)
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

type ContainerInfo struct {
	containerIp  string
	dockerServer string
}

func getFromRedis(domain string) ContainerInfo {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisServer,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var info ContainerInfo
	info.containerIp = rdb.HGet(ctx, domain, "container_ip").Val()
	info.dockerServer = rdb.HGet(ctx, domain, "docker_server").Val()

	defer rdb.Close()
	return info
}
