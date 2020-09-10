package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

var (
	App *Config
)

func Read(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	err = d.Decode(&App)
	if err != nil {
		return err
	}

	return nil
}
