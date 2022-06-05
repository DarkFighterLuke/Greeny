package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

var configPath string

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Schema   string `yaml:"schema"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Email struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"email"`
	Jwt struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

// SetConfigPath Set a configuration path and validate it
func SetConfigPath(path string) error {
	configPath = path
	err := ValidateConfigPath(configPath)
	return err
}

//GetConfig Return a new Config
func GetConfig() (*Config, error) {
	conf, err := NewConfig(configPath)
	return conf, err
}

//NewConfig Given a configuration path read a new Config
func NewConfig(congigPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(congigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

//ValidateConfigPath Check if path is valid or not
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
