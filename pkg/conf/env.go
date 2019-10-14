package conf

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type (
	config struct {
		MaxProc    int        `yaml:"maxProc"`
		Port       int        `yaml:"port"`
		Domain     string     `yaml:"domain",omitempty`
		Storage    string     `yaml:"uploadDir"`
		MaxAge     int        `yaml:"maxAge"`
		DefaultImg string     `yaml:"defaultImg"`
		Type       string     `yaml:"type"`
		JWT        *jwt       `yaml:"jwt"`
		Redis      *redis     `yaml:"redis"`
		Aliyun     *aliyun    `yaml:"aliyun"`
		WaterMark  *watermark `yarml:"watermark"`
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
	watermark struct {
		Enable      bool    `yaml:"enable"`
		Text        string  `yaml:"text"`
		Font        string  `yaml:"font"`
		Gravity     string  `yaml:"gravity"`
		PointSize   float64 `yaml:"pointsize"`
		Color       string  `yarml:"color"`
		gravityType imagick.GravityType
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
		JWT:        &jwt{Secret: "123456", Algorithm: "HS512"},
		Redis: &redis{
			Addr:     "127.0.0.1:6379",
			DB:       1,
			PoolSize: 100,
			Expire:   10800},
		WaterMark: &watermark{
			Enable:    false,
			Gravity:   "southeast", //center,northwest, northeast,southwest, southeast
			Font:      "cochin.ttc",
			PointSize: 36,
			Color:     "white",
		},
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
	if c.WaterMark.Enable {
		c.WaterMark.convertGravityType()
	}
	return c
}

func (e *watermark) GravityType() imagick.GravityType {
	return e.gravityType
}
func (e *watermark) convertGravityType() {
	gravity := strings.ToLower(e.Gravity)
	switch gravity {
	case "center":
		e.gravityType = imagick.GRAVITY_CENTER
		break
	case "northwest":
		e.gravityType = imagick.GRAVITY_NORTH_WEST
		break
	case "northeast":
		e.gravityType = imagick.GRAVITY_NORTH_EAST
		break
	case "southwest":
		e.gravityType = imagick.GRAVITY_SOUTH_WEST
		break
	case "southeast":
		e.gravityType = imagick.GRAVITY_SOUTH_EAST
		break
	}
}
