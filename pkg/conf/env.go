package conf

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
)

type (
	config struct {
		MaxProc    int    `yaml:"maxProc"`
		Port       int    `yaml:"port"`
		Domain     string `yaml:"domain",omitempty`
		Storage    string `yaml:"uploadDir"`
		MaxAge     int32  `yaml:"maxAge"`
		DefaultImg string `yaml:"defaultImg"`
		Type       string `yaml:"type"`
		JWT        jwt    `yaml:"jwt"`
		Redis      redis  `yaml:"redis"`
		Aliyun     aliyun `yaml:"aliyun"`
	}

	jwt struct {
		Secret    string `yaml:"secret"`
		Algorithm string `yaml:"algorithm",omitempty`
	}

	redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		PoolSize int    `yaml:"poolSize"`
		Expire   int32  `yaml:"expire"`
	}
	aliyun struct {
		Region          string `yaml:"region"`
		AccessKeyId     string `yaml:"accessKeyId"`
		AccessKeySecret string `yaml:"accessKeySecret"`
		Bucket          string `yaml:"bucket"`
		Secure          bool   `yaml:"secure"`
	}
)

var Config *config

func init() {
	currentDir, _ := os.Getwd()
	Config = load(currentDir + "/conf.yml")

	log.Print(Config)
}

func load(file string) *config {
	//default value
	var c = &config{
		Port:       3000,
		Storage:    "/data/images",
		MaxAge:     31536000,
		DefaultImg: "default.png",
		Type:       "local",
		JWT:        jwt{Secret: "123456", Algorithm: "HS512"},
		Redis: redis{
			Addr:     "127.0.0.1:6379",
			DB:       1,
			PoolSize: 100,
			Expire:   10800},
	}
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	// 设置go processor数量
	if c.MaxProc == 0 {
		c.MaxProc = runtime.NumCPU()
	}
	return c
}
