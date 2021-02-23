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
		//c.Logger().Print(parseHost[0])
		// 查询子域名获得ip地址
		c.Logger().Print(fmt.Sprintf("http://%s:2090", info.containerIp))
		c.Logger().Print(info.dockerServer)
		req.Header.Add("container", fmt.Sprintf("http://%s:2090", info.containerIp))
		return fmt.Sprintf("https://%s", info.dockerServer)
	}

	e.Use(proxy_middleware.Proxy(proxy_middleware.NewRoundRobinBalancer()))

	// go p90(e, conf)
	e.Logger.Print("Entry Proxy For Node Router")
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, conf.Port), conf.CertFile, conf.KeyFile))
}

func p90(e *echo.Echo, conf config.Conf) {
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, "90"), conf.CertFile, conf.KeyFile))
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
