package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Conf struct {
	Host        string
	Port        string
	DockerPort  string
	CertFile    string
	KeyFile     string
	Api         string
	RedisServer string
}

func Load() Conf {
	var config Conf
	f, err := os.Open(filepath.ToSlash("./config.json"))
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		log.Println(err.Error())
	}

	return config
}
