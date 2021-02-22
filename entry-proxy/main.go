package main

import (
	"../config"
	"../proxy-middleware"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"strings"
)

var conf = config.Load()

func main() {
	e := echo.New()

	proxy_middleware.GetTarget = func(c echo.Context) string {
		req := c.Request()

		parseHost := strings.Split(req.Host, ":")
		info := getFromRedis(parseHost[0])
		//fmt.Println(parseHost[0])
		// 查询子域名获得ip地址
		//fmt.Println(fmt.Sprintf("http://%s:2090", info.container_ip))
		//fmt.Println(info.docker_server)
		req.Header.Add("container", fmt.Sprintf("http://%s:2090", info.container_ip))
		return fmt.Sprintf("https://%s", info.docker_server)
	}

	e.Use(proxy_middleware.Proxy(proxy_middleware.NewRoundRobinBalancer()))

	// go p90(e, conf)
	fmt.Println("Entry Proxy For Node Router")
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, conf.Port), conf.CertFile, conf.KeyFile))
}

func p90(e *echo.Echo, conf config.Conf) {
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, "90"), conf.CertFile, conf.KeyFile))
}

type ContainerInfo struct {
	container_ip  string
	docker_server string
}

func getFromRedis(domain string) ContainerInfo {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisServer,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var info ContainerInfo
	info.container_ip = rdb.HGet(ctx, domain, "container_ip").Val()
	info.docker_server = rdb.HGet(ctx, domain, "docker_server").Val()

	defer rdb.Close()
	return info
}
