module github.com/programschool/proxy-service/docker-service

go 1.16

require (
	github.com/go-redis/redis/v8 v8.8.0
	github.com/programschool/proxy-service/config v0.0.0-incompatible
)

replace github.com/programschool/proxy-service/config => ../config
