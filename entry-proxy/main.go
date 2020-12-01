package main

import (
	"../config"
	"../proxy-middleware"
	"fmt"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	proxy_middleware.GetTarget = func(c echo.Context) string {
		req := c.Request()
		//res := c.Response()
		//fmt.Println(fmt.Sprintf("Proxy: %s", req.Host))

		// 查询子域名获得ip地址
		req.Header.Add("container", "http://172.17.0.2:8080")
		return "https://192.168.10.104"
	}

	e.Use(proxy_middleware.Proxy(proxy_middleware.NewRoundRobinBalancer()))

	conf := config.Load()
	// go p90(e, conf)
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, conf.Port), conf.CertFile, conf.KeyFile))
}

func p90(e *echo.Echo, conf config.Conf) {
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, "90"), conf.CertFile, conf.KeyFile))
}
