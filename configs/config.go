package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire string `mapstructure:"expire"`
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "myapp.db")

	// viper.ReadInConfig()： 从磁盘读取配置文件，并把内容加载进 Viper 的内部存储
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %v", err)
		log.Println("Using default values and environment variables")
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func Load() (*Config, error) {

	// 配置文件名称、类型与路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "myapp.db")

	// 从磁盘加载配置文件到viper中，不存在则使用默认值
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %v", err)
		log.Println("Using default values and environment variables")
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// 从viper解析配置
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}
