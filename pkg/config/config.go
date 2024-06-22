package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
	"time"
)

var Data Config
var once sync.Once
var err error

type Db struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type App struct {
	ListeningPort  int           `mapstructure:"http_port"`
	DrainingPeriod time.Duration `mapstructure:"draining_period"`
	LogLevel       string        `mapstructure:"log_level"`
}

type Config struct {
	Db  Db  `mapstructure:"db"`
	App App `mapstructure:"app"`
}

func CreateConfig() (Config, error) {
	viper.SetConfigName("config")  // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("config/") // path to look for the config file in

	// Set environment variable prefix
	viper.SetEnvPrefix("APP")
	// Replace dots with underscores in environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("log_level", "INFO")
	viper.SetDefault("http_port", 8080)
	viper.SetDefault("draining_period", 30)

	// Read the config file
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("config file not found, skipping")
		} else {
			log.Fatal("something went wrong while config reading")
			return Config{}, err
		}
	}

	var C Config

	// Unmarshal the config into the struct
	err = viper.Unmarshal(&C)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
		return C, err
	}
	return C, nil
}

func GetConfig() Config {
	once.Do(func() {
		Data, err = CreateConfig()
	})
	if err != nil {
		log.Fatal("Server Shutdown: unable to read config", err)
	}
	return Data
}
