package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

const configFilePath string = "etc/config.yaml"

var config *Config

type Config struct {
	Version   string `yaml:"version" mapstructure:"version"`
	JwtSecret string `yaml:"jwtSecret" mapstructure:"jwtSecret"`
	*Server   `yaml:"server" mapstructure:"server"`
	*Mysql    `yaml:"mysql" mapstructure:"mysql"`
}

type Server struct {
	Addr string `yaml:"addr" mapstructure:"addr"`
}

type Mysql struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	Port     string `yaml:"port" mapstructure:"port"`
	Username string `yaml:"username" mapstructure:"username"`
	Database string `yaml:"database" mapstructure:"database"`
	Password string `yaml:"password" mapstructure:"password"`
}

func GetConfig() *Config {
	if config == nil {
		Init()
	}
	return config
}

// 使用viper去获取配置信息
func Init() (err error) {
	config = new(Config)
	//viper.SetConfigFile("./config.yaml")
	viper.SetConfigName("config") //指定配置文件名称（不需要带后缀）
	viper.SetConfigType("yaml")   //指定配置文件类型
	viper.AddConfigPath(".")      //指定查找配置文件的路径（这里使用相对路径）

	err = viper.ReadInConfig() //读取配置文件信息
	if err != nil {
		log.Printf("viper readinconfig failed err :", err)
		return
	}
	//把信息读到结构体中去，反序列化进去
	if err = viper.Unmarshal(config); err != nil {
		log.Println("viper unmarshal failed ", err)

	}
	viper.WatchConfig() //检测配置信息是否改变，有改变会重新渲染
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err = viper.Unmarshal(config); err != nil {
			log.Println("viper unmarshal failed ", err)
		}
		log.Println("config is changed!")
	})
	return
}
