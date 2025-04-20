/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var cfgFilePath string

type Config struct {
	ServerConf serverConfig
}

type serverConfig struct {
	Default bool
	Host    string
	Port    uint
}

func (conf serverConfig) host() string {
	return fmt.Sprintf("%s:%d", conf.Host, conf.Port)
}

func createDefaultConfig() Config {
	// server config
	myServerConfig := serverConfig{
		Default: true,
		Host:    "",
		Port:    8080,
	}
	// Main Config
	mainConfig := Config{
		ServerConf: myServerConfig,
	}
	return mainConfig
}

func (conf Config) Host(configName string) (string, error) {
	if configName == "server" {
		return conf.ServerConf.host(), nil
	}
	return "", errors.New("unsupported config")
}

func (conf Config) updateServerConf() error {
	host := viper.GetString("host")
	if len(host) == 0 {
		conf.ServerConf.Host = host
	}
	port := viper.GetUint("port")
	if port > 0 {
		conf.ServerConf.Port = port
	}
	return nil
}

func ReadConfigFile(cfgFilePath string) (Config, error) {
	mainConfig := createDefaultConfig()
	viper.AddConfigPath(cfgFilePath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Using default!")
			return mainConfig, nil
		} else {
			// Config file was found but another error was produced
		}
	} else {
		mainConfig.updateServerConf()
	}

	return mainConfig, err
}
