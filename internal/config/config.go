package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct{
	Env string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct{
	Address string `yaml:"address" env-default:"localhost:8000"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`

}

func MustLoad()*Config{
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == ""{
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err){
		log.Fatal("config file does not exist: ", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil{
		log.Fatal("cannot read config: ", err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil{
		log.Fatal("cannot unmarshal config: ", err)
	}

	return &cfg
}