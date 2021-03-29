package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/programschool/proxy-service/config"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var conf = config.Load()

func main() {
	go listen80()

	// e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, conf.Port), conf.CertFile, conf.KeyFile))

	http.HandleFunc("/", Handle())
	_ = http.ListenAndServeTLS(fmt.Sprintf("%s:%s", conf.Host, conf.Port), conf.CertFile, conf.KeyFile, nil)
}

func listen80() {
	http.HandleFunc("/", Handle())
	_ = http.ListenAndServeTLS(fmt.Sprintf("%s:%s", conf.Host, "80"), conf.CertFile, conf.KeyFile, nil)
}

func Handle() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//log.Println("request:", r.RemoteAddr, "want", r.RequestURI)
		//Many webservers are configured to not serve pages if a request doesn’t appear from the same host.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		//auth := r.Header.Get("Docker-Auth")
		//w.Header().Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(auth)))

		_, port, _ := net.SplitHostPort(r.Host)
		p := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   r.Header.Get("Docker-Ip") + ":" + port,
		})
		//log.Printf("respond ip %s", r.Header.Get("Docker-Ip"))

		p.ServeHTTP(w, r)
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

	log.Print(info.containerIp)
	log.Print(info.dockerServer)

	defer rdb.Close()
	return info
}
