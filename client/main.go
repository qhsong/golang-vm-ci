package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/qhsong/golang-vm-ci/common"
)

type Config struct {
	ControllURL string
}

func getConfig() *Config {
	content, err := ioutil.ReadFile("/etc/client.conf")
	var conf Config
	if err != nil {
		log.Println("Error to read config file")
		conf.ControllURL = "192.168.122.1"
	} else {
		if _, err := toml.Decode(string(content), &conf); err != nil {
			log.Println("Error to Decode config file")
			conf.ControllURL = "192.168.122.1"
		}
	}
	return &conf
}

func main() {
	conf := getConfig()

	serverURL := fmt.Sprintf("http://%s/", conf.ControllURL)
	resp, err := http.Get(serverURL + "client/")
	if err != nil {
		log.Fatal("Unable to connect host server")
	}
	defer resp.Body.Close()

	jd := json.NewDecoder(resp.Body)
	var task common.Task
	err = jd.Decode(&task)
	if err != nil {
		log.Fatal("Can not decode json string")
	}
}
