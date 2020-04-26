package util

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const (
	SystemWindows = iota
	SystemLinux
)

const TCLOUD_API_VERSION = "v4"
const TCLOUD_TARGET_SYSTEM = SystemWindows

type WebConfig struct {
	Domain  string `yaml:"domain"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}
type DBConfig struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DataBase string `yaml:"database"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Session  int    `yaml:"session"`
	Register int    `yaml:"register"`
	Upload   int    `yaml:"upload"`
	Expired  int    `yaml:"expired"`
}

type LogConfig struct {
	Path  string `yaml:"path"`
	Level int    `yaml:"level"`
}

type MailConfig struct {
	Disable  bool   `yaml:"disable"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Passowrd string `yaml:"password"`
}

type FileConfig struct {
	Path string `yaml:"path"`
}

type Config struct {
	Web      WebConfig   `yaml:"web"`
	Database DBConfig    `yaml:"database"`
	Redis    RedisConfig `yaml:"redis"`
	Log      LogConfig   `yaml:"log"`
	Mail     MailConfig  `yaml:"mail"`
	File     FileConfig  `yaml:"file"`
}

var p *Config = nil

func GetConfig() *Config {
	if p == nil {
		p = new(Config)
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file\n", err)
		}

		filePath := os.Getenv("CONFIG_PATH")
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Error loading %s\n%e", filePath, err)
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("Error read config file\n", err)
		}
		print(data)

		err = yaml.Unmarshal(data, p)
		if err != nil {
			log.Fatal("Error parse config file\n", err)
		}
	}
	return p
}
