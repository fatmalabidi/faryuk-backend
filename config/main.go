package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/jinzhu/configor"
)

func MakeConfig() (*Config, error) {
	var configFilePath string

	newConfig := configor.New(&configor.Config{})

	getConfigFile := func() string {
		 fmt.Println("ENVIRONNEMENT",getEnvironment())
		switch getEnvironment() {
		case "test":
			configFilePath = "./config.test.yml"
		default:
			configFilePath = "./config.yml"
		}
		_, filename, _, _ := runtime.Caller(0)
		if strings.Contains(strings.ToLower(os.Args[0]), "test") {
			return path.Join(path.Dir(filename), "./config.test.yml")
		}
		return path.Join(path.Dir(filename), configFilePath)
	}

	conf := new(Config)
	err := newConfig.Load(conf, getConfigFile())
	fmt.Println("CONFIG LOADED FOR ENV: \nCONFIGOR_ENV=", os.Getenv("CONFIGOR_ENV"), "\nconf=",&conf)
	fmt.Println("ERROR:", err)
	return conf, err
}

func getEnvironment() string {
	if env := os.Getenv("CONFIGOR_ENV"); env != "" {
		return env
	}

	return "dev"
}

type Config struct {
	Server struct {
		Port int    `yaml:"Port"`
		Host string `yaml:"Host"`
	} `yaml:"Server"`

	Database struct {
		DbType string `yaml:"DbType"`
		Port   string `yaml:"Port"`
		Host   string `yaml:"Host"`
		DbName string `yaml:"DbName"`
	} `yaml:"Database"`
}
