package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var Config config

type config struct {
	Server struct {
		DefaultPorts []string `yaml:"defaultPorts"`
		Host         string   `yaml:"host"`
		ExpireTime   int      `yaml:"expireTime"`
	}
}

func init() {
	bytes, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
		return
	}
	if err = yaml.Unmarshal(bytes, &Config); err != nil {
		panic(err)
		return
	}

}
