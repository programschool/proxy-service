package main

import (
	"../config"
	"../proxy-middleware"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"log"
	"strings"
)

var conf = config.Load()

func main() {
	e := echo.New()
	proxy := proxy_middleware.NewProxy{}

	proxy.GetTarget = func(c echo.Context) string {
		c.Logger().Print("listen 2090")
		req := c.Request()

		parseHost := strings.Split(req.Host, ":")
		info := getFromRedis(parseHost[0])
		c.Logger().Print(parseHost[0])
		//c.Logger().Print(parseHost[1])
		// 查询子域名获得ip地址
		c.Logger().Print(fmt.Sprintf("container http://%s:2090", info.containerIp))
		c.Logger().Print(info.dockerServer)
		req.Header.Add("container", fmt.Sprintf("http://%s:2090", info.containerIp))
		return fmt.Sprintf("https://%s", info.dockerServer)
	}

	go listen80()

	e.Use(proxy.Proxy(proxy.NewRoundRobinBalancer()))
	e.Logger.Print("Entry Proxy For Node Router")
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, conf.Port), conf.CertFile, conf.KeyFile))
}

func listen80() {
	e := echo.New()
	proxy := proxy_middleware.NewProxy{}

	proxy.GetTarget = func(c echo.Context) string {
		c.Logger().Print("listen 80")
		req := c.Request()

		parseHost := strings.Split(req.Host, ":")
		info := getFromRedis(parseHost[0])
		c.Logger().Print(parseHost[0])
		//c.Logger().Print(parseHost[1])
		// 查询子域名获得ip地址
		c.Logger().Print(fmt.Sprintf("container http://%s:80", info.containerIp))
		c.Logger().Print(info.dockerServer)
		req.Header.Add("container", fmt.Sprintf("http://%s:80", info.containerIp))
		return fmt.Sprintf("http://%s", info.dockerServer)
	}

	e.Use(proxy.Proxy(proxy.NewRoundRobinBalancer()))
	e.Logger.Print("Entry Proxy For Node Preview Router")
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", conf.Host, "80")))
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
