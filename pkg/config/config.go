package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DebugEnabled bool          `yaml:"debug_enabled" default:"true" usage:"Show debugs"`
	PluginName   string        `yaml:"plugin_name" default:"docker" usage:"The name of the plugin"`
	ChannelName  string        `yaml:"channel_name" default:"channel" usage:"The name of the channel"`
	Port         int           `yaml:"port" default:"8080" usage:"The port to listen on"`
	Channel      ChannelConfig `yaml:"channel" default:"" usage:"The name of the channel"`
}

type ChannelConfig struct {
	Host        string `yaml:"host" default:"localhost:8080" usage:"The remote host"`
	Port        int    `yaml:"port" default:"8080" usage:"The port to listen on"`
	ChannelName string `yaml:"channel_name" default:"channel" usage:"The name of the channel"`
}

func Load(cfgFile string) (*Config, error) {
	var config Config
	yamlFile, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
