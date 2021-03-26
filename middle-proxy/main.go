package main

import (
	"../config"
	"../proxy-middleware"
	"fmt"
	"github.com/labstack/echo"
)

var conf = config.Load()

func main() {
	e := echo.New()
	proxy := proxy_middleware.NewProxy{}

	proxy.GetTarget = func(c echo.Context) string {
		req := c.Request()
		// res := c.Response()
		req.Header.Set("Cache-Control", "no-cache, private, max-age=0")
		req.Header.Set("Pragma", "no-cache")
		c.Logger().Print(fmt.Sprintf("Proxy: %s", req.Host))
		c.Logger().Print(req.Header.Get("container"))
		return req.Header.Get("container")
	}

	go listen80()
	e.Use(proxy.Proxy(proxy.NewRoundRobinBalancer()))
	e.Logger.Print("Middle Proxy For Container Node")
	e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%s", conf.Host, conf.Port), conf.CertFile, conf.KeyFile))
}

func listen80() {
	e := echo.New()
	proxy := proxy_middleware.NewProxy{}

	proxy.GetTarget = func(c echo.Context) string {
		req := c.Request()
		// res := c.Response()
		c.Logger().Print(fmt.Sprintf("Proxy: %s", req.Host))
		c.Logger().Print(req.Header.Get("container"))
		req.Header.Set("Cache-Control", "no-cache, private, max-age=0")
		req.Header.Set("Pragma", "no-cache")
		return req.Header.Get("container")
	}

	e.Use(proxy.Proxy(proxy.NewRoundRobinBalancer()))

	e.Logger.Print("Middle Proxy For Container Node")
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", conf.Host, "80")))
}
