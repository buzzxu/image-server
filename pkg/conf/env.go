package conf

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"gopkg.in/gographics/imagick.v3/imagick"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

type (
	config struct {
		MaxProc    int        `yaml:"maxProc"`
		Port       int        `yaml:"port"`
		Domain     string     `yaml:"domain",omitempty`
		Storage    string     `yaml:"uploadDir"`
		MaxAge     int        `yaml:"maxAge"`
		BodyLimit  string     `yaml:"bodyLimit"`
		SizeLimit  string     `yaml:"sizeLimit"`
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
		Addr       string        `yaml:"addr"`
		Password   string        `yaml:"password"`
		DB         int           `yaml:"db"`
		PoolSize   int           `yaml:"poolSize"`
		Expire     int32         `yaml:"expire"`
		Expiration time.Duration `yaml:"-"`
	}
	aliyun struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKeyId     string `yaml:"accessKeyId"`
		AccessKeySecret string `yaml:"accessKeySecret"`
		Bucket          string `yaml:"bucket"`
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
	println(fmt.Sprintf("Port:%x,Domain:%s,Type:%s,Storage:%s", Config.Port, Config.Domain, Config.Type, Config.Storage))
}

func load(file string) *config {
	//default value
	var c = &config{
		Port:       3000,
		Domain:     "",
		Storage:    "/data/images",
		MaxAge:     31536000,
		BodyLimit:  "5M",
		SizeLimit:  "500K",
		DefaultImg: "default.png",
		Type:       "local",
		JWT:        &jwt{Secret: "123456", Algorithm: "HS512"},
		Redis: &redis{
			Addr:     "127.0.0.1:6379",
			DB:       1,
			PoolSize: runtime.NumCPU() * 20,
			Password: "",
			Expire:   10800},
		WaterMark: &watermark{
			Enable:    false,
			Gravity:   "southeast", //center,northwest, northeast,southwest, southeast
			Font:      "cochin.ttc",
			PointSize: 36,
			Color:     "white",
		},
	}
	if isConfExsits(file) {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("读取配置文件内容失败,%v ", err)
		}
		if yaml.Unmarshal(yamlFile, c); err != nil {
			log.Fatalf("解析配置文件失败: %v", err)
		}
	} else {
		log.Warnf("未发现配置文件:%s,使用默认配置", file)
	}

	// 设置go processor数量
	if c.MaxProc == 0 {
		c.MaxProc = runtime.NumCPU()
	}
	if c.WaterMark.Enable {
		c.WaterMark.convertGravityType()
	}
	if c.Type == "local" {
		c.Redis.env()
	}
	if c.Domain == "none" {
		c.Domain = ""
	}
	return c
}

func (e watermark) GravityType() imagick.GravityType {
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

func (r *redis) env() {
	r.Expiration = time.Duration(r.Expire) * time.Second
	if r.PoolSize == 0 {
		r.PoolSize = runtime.NumCPU() * 20
	}
	if r.Password == "none" {
		r.Password = ""
	}
}

func isConfExsits(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}
