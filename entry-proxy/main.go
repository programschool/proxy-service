package main

import (
	"../config"
	"../proxy-middleware"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	proxy_middleware.GetTarget = func(c echo.Context) string {
		req := c.Request()
		info := getFromRedis(req.Host)

		// 查询子域名获得ip地址

		fmt.Println(fmt.Sprintf("http://%s:8080", info.container_ip))

		req.Header.Add("container", fmt.Sprintf("http://%s:8080", info.container_ip))
		return fmt.Sprintf("https://%s", info.docker_server)
	}

	e.Use(proxy_middleware.Proxy(proxy_middleware.NewRoundRobinBalancer()))

	conf := config.Load()
	// go p90(e, conf)
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
		Addr:     "192.168.10.102:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var info ContainerInfo
	info.container_ip = rdb.HGet(ctx, domain, "container_ip").Val()
	info.docker_server = rdb.HGet(ctx, domain, "docker_server").Val()

	return info
	// defer rdb.Close()
}
