package main

import (
	"../middleware"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer()))

	e.Logger.Fatal(e.StartTLS("0.0.0.0:443", "../ssl/boxlayer.com/fullchain.pem", "../ssl/boxlayer.com/privkey.pem"))
}
