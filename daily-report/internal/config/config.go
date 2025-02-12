package config

import (
	"fmt"
	"log"
	"sync"

	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置结构
type Config struct {
	Database struct {
		Driver    string `yaml:"driver"`
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Name      string `yaml:"name"`
		Charset   string `yaml:"charset"`
		ParseTime bool   `yaml:"parseTime"`
		Loc       string `yaml:"loc"`
	} `yaml:"database"`

	Server struct {
		Port int    `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
}

var (
	config *Config
	once   sync.Once
)

// GetConfig 返回配置实例
func GetConfig() *Config {
	once.Do(func() {
		config = &Config{}

		// 读取配置文件
		data, err := os.ReadFile("config/config.yaml")
		if err != nil {
			log.Fatal("Failed to read config file:", err)
		}

		// 解析配置文件
		err = yaml.Unmarshal(data, config)
		if err != nil {
			log.Fatal("Failed to parse config file:", err)
		}
	})
	return config
}

// GetDSN 返回数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.Charset,
		c.Database.Loc,
	)
}
