module github.com/programschool/proxy-service/docker-service

go 1.16

require (
	github.com/go-redis/redis/v8 v8.8.0
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/labstack/echo/v4 v4.2.1
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/programschool/proxy-service/config v0.0.0-incompatible
	github.com/programschool/proxy-service/library/dockertools v0.0.0-00010101000000-000000000000
	gotest.tools/v3 v3.0.3 // indirect
)

replace (
	github.com/programschool/proxy-service/config => ../config
	github.com/programschool/proxy-service/library/dockertools => ../library/dockertools
)
