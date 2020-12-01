// +build go1.11
package proxy_middleware

import (
	"crypto/tls"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httputil"
)

func proxyHTTP(tgt *ProxyTarget, c echo.Context, config ProxyConfig) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(tgt.URL)
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		desc := tgt.URL.String()
		if tgt.Name != "" {
			desc = fmt.Sprintf("%s(%s)", tgt.Name, tgt.URL.String())
		}
		c.Set("_error", echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("remote %s unreachable, could not forward: %v", desc, err)))
	}
	//proxy.Transport = config.Transport
	// 跳过ssl验证
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	proxy.ModifyResponse = config.ModifyResponse

	return proxy
}
