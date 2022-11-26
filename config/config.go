package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type TGBotConfig struct {
	Token string `yaml:"token"`
	Key   string `yaml:"key"`
}

func ReadConfig(filename string) *TGBotConfig {
	fs, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}

	conf := TGBotConfig{}
	err = yaml.NewDecoder(fs).Decode(&conf)
	if err != nil {
		log.Panic(err)
	}

	return &conf
}
