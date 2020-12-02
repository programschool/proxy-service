package main

import (
	"../config"
	"../library/dockertools"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

var conf config.Conf

func main() {
	// Echo instance
	e := echo.New()
	conf = config.Load()
	e.POST("/", handle)
	// Start server
	address := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	e.Logger.Fatal(e.StartTLS(address, conf.CertFile, conf.KeyFile))
}

func handle(c echo.Context) error {
	post := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&post)
	if err != nil {
		return err
	}

	w := c.Response()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// service is container node ip addres
	server := post["server"].(string)
	dockerHost := fmt.Sprintf("%s:%s", server, conf.DockerPort)

	dock := dockertools.Dock{}.New(dockerHost)
	defer dock.Close()

	var data map[string]interface{}

	switch post["action"] {
	case "create":
		image := fmt.Sprintf("%s/%s", post["registry-host"], post["image"].(string))
		dockerID := dock.CreateAndStart(
			image,
			int64(post["memory"].(float64)),
			post["size"].(string))
		// dock.ExecCommand(dockerID, []string{"bash", "/etc/rc.local"})
		dockerIP := dock.IpAddress(dockerID)
		// fmt.Println(dockerID)
		// fmt.Println(dockerIP)
		// 获取返回结果
		data = map[string]interface{}{
			"docker_id": dockerID,
			"docker_ip": dockerIP,
			"server":    server,
		}
		domain := post["domain"].(string)
		saveToRedis(domain, "docker_ip", dockerIP)
		saveToRedis(domain, "docker_server", server)
	case "start":
		dockerID := post["docker-id"].(string)
		isRun := dock.IsRunning(dockerID)
		flag := 0
		if !isRun {
			dock.Start(dockerID)
			userDomain := post["domain"].(string)
			dockerIP := dock.IpAddress(dockerID)
			saveToRedis(userDomain, "docker_ip", dockerIP)
			server := post["server"].(string)
			saveToRedis(userDomain, "docker_server", server)
			saveToRedis(userDomain, "docker_id", dockerID)
			flag = 1
		}
		//fmt.Printf("start docker, id %s, isrun %t\n", dockerID, isRun)
		// dock.ExecCommand(dockerID, []string{"bash", "/etc/rc.local"})
		data = map[string]interface{}{
			"error": 0,
			"flag":  flag,
		}
	case "reStart":
		dockerID := post["docker-id"].(string)
		dock.ReStart(dockerID)
		userDomain := post["domain"].(string)
		dockerIP := dock.IpAddress(dockerID)
		saveToRedis(userDomain, "docker_ip", dockerIP)
		server := post["server"].(string)
		saveToRedis(userDomain, "docker_server", server)
		saveToRedis(userDomain, "docker_id", dockerID)
		// dock.ExecCommand(dockerID, []string{"bash", "/etc/rc.local"})
		data = map[string]interface{}{
			"error": 0,
		}
	case "stop":
		dockerID := post["docker-id"].(string)
		dock.Stop(dockerID)
		data = map[string]interface{}{
			"error": 0,
		}
	case "remove":
		dockerID := post["docker-id"].(string)
		dock.Stop(dockerID)
		dock.Remove(dockerID)
		data = map[string]interface{}{
			"error": 0,
		}
	case "bash":
		dockerID := post["docker-id"].(string)
		command := strings.Split(post["command"].(string), " ")
		userDomain := post["domain"].(string)
		bash := []string{"bash"}
		//fmt.Println("exec bash ", command[0], command[1])
		if len(command) >= 2 {
			if strings.TrimSpace(command[0]) == "startOnlineEditor" {
				saveToRedis(userDomain, "auth", strings.TrimSpace(command[1]))
			}
		}
		inspect, err := dock.ExecCommand(dockerID, append(bash, command...))
		if err == nil {
			data = map[string]interface{}{
				"error":        0,
				"action":       post["action"],
				"command":      command,
				"container-id": inspect.ContainerID,
			}
		} else {
			data = map[string]interface{}{
				"error":        1,
				"action":       post["action"],
				"command":      command,
				"container-id": inspect.ContainerID,
			}
		}
		//fmt.Printf("exec bash finish %v\n", inspect)
	default:
		//
		data = map[string]interface{}{
			"error":   2,
			"message": "未知操作，请检查命令",
		}
	}

	encodeData, _ := json.Marshal(data)
	c.String(200, string(encodeData))

	return nil
}

func saveToRedis(domain string, key string, val string) {
	var ctx = context.Background()

	//fmt.Println(domain)
	//fmt.Println(key)
	//fmt.Println(val)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.10.102:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.HMSet(ctx, domain, key, val)

	if err != nil {
		fmt.Println(err)
	} else {
		defer rdb.Close()
	}
}
