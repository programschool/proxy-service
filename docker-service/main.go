package main

import (
	"../config"
	"../library/dockertools"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"log"
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
	fmt.Println("Docker Service API")
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
		containerID := dock.CreateAndStart(
			image,
			int64(post["memory"].(float64)),
			post["size"].(string))
		// dock.ExecCommand(containerID, []string{"bash", "/etc/rc.local"})
		dockerIP := dock.IpAddress(containerID)
		// fmt.Println(containerID)
		// fmt.Println(dockerIP)
		// 获取返回结果
		data = map[string]interface{}{
			"container_id": containerID,
			"container_ip": dockerIP,
			"server":       server,
		}
		domain := post["domain"].(string)
		saveToRedis(domain, "container_ip", dockerIP)
		saveToRedis(domain, "docker_server", server)
	case "start":
		containerID := post["container-id"].(string)
		isRun := dock.IsRunning(containerID)
		flag := 0
		if !isRun {
			dock.Start(containerID)
			userDomain := post["domain"].(string)
			dockerIP := dock.IpAddress(containerID)
			saveToRedis(userDomain, "container_ip", dockerIP)
			server := post["server"].(string)
			saveToRedis(userDomain, "docker_server", server)
			saveToRedis(userDomain, "container_id", containerID)
			flag = 1
		}
		//fmt.Printf("start docker, id %s, isrun %t\n", containerID, isRun)
		// dock.ExecCommand(containerID, []string{"bash", "/etc/rc.local"})
		data = map[string]interface{}{
			"error": 0,
			"flag":  flag,
		}
	case "reStart":
		containerID := post["container-id"].(string)
		dock.ReStart(containerID)
		userDomain := post["domain"].(string)
		dockerIP := dock.IpAddress(containerID)
		saveToRedis(userDomain, "container_ip", dockerIP)
		server := post["server"].(string)
		saveToRedis(userDomain, "docker_server", server)
		saveToRedis(userDomain, "container_id", containerID)
		// dock.ExecCommand(containerID, []string{"bash", "/etc/rc.local"})
		data = map[string]interface{}{
			"error": 0,
		}
	case "stop":
		containerID := post["container-id"].(string)
		dock.Stop(containerID)
		data = map[string]interface{}{
			"error": 0,
		}
	case "remove":
		containerID := post["container-id"].(string)
		dock.Stop(containerID)
		err := dock.ContainerRemove(containerID)
		if err != nil {
			c.Logger().Errorf("containerID %s: 删除失败", containerID)
		}
		data = map[string]interface{}{
			"error": 0,
		}
	case "save":
		containerID := post["container-id"].(string)
		imageName := post["image"].(string)
		dock.Stop(containerID)
		img, _ := dock.Commit(containerID, imageName)
		errPush := dock.Push(imageName, post["username"].(string), post["password"].(string))
		if errPush != nil {
			c.Logger().Errorf("%s: PUSH 失败", imageName)
		}
		errCR := dock.ContainerRemove(containerID)
		if errCR != nil {
			c.Logger().Errorf("containerID %s: 删除失败", containerID)
		}
		dock.ImageRemove(img.ID)
		data = map[string]interface{}{
			"error": 0,
		}
	case "bash":
		containerID := post["container-id"].(string)
		command := strings.Split(post["command"].(string), " ")
		userDomain := post["domain"].(string)
		bash := []string{"bash"}
		//fmt.Println("exec bash ", command[0], command[1])
		if len(command) >= 2 {
			if strings.TrimSpace(command[0]) == "startOnlineEditor" {
				saveToRedis(userDomain, "auth", strings.TrimSpace(command[1]))
			}
		}
		inspect, err := dock.ExecCommand(containerID, append(bash, command...))
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

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisServer,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.HMSet(ctx, domain, key, val)

	if err != nil {
		log.Println("Error:")
		log.Println(err)
		log.Println(domain)
		log.Println(key)
		log.Println(val)
	} else {
		defer rdb.Close()
	}
}
