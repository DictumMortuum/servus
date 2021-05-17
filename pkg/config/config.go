package config

import (
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Timezone string `yaml:"timezone"`
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Server   string `yaml:"server"`
		Port     string `yaml:"port"`
	} `yaml:"database"`
	Zerotier struct {
		Token   string `yaml:"token"`
		Network string `yaml:"network"`
	} `yaml:"zerotier"`
	Mpd struct {
		Server string `yaml:"server"`
		Port   string `yaml:"port"`
	}
	Telegram struct {
		Enabled bool     `yaml:"enabled" default:"false"`
		Token   string   `yaml:"token"`
		Users   []string `yaml:"users"`
	}
	Calendar struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Server   string `yaml:"server"`
	} `yaml:"calendar"`
	Atlas struct {
		ClientId string `yaml:"client_id"`
	} `yaml:"atlas"`
	Cache struct {
		Expiration int    `yaml:"expiration"`
		Cleanup    int    `yaml:"cleanup"`
		Path       string `yaml:"path"`
	} `yaml:"cache"`
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

func (c *Config) GetMariaDBConnection(db string) string {
	rs := []string{
		c.Database.Username,
		":",
		c.Database.Password,
		"@tcp(",
		c.Database.Server,
		":",
		c.Database.Port,
		")/",
		db,
		"?parseTime=true",
		"&loc=" + url.QueryEscape(c.Timezone),
	}

	return strings.Join(rs, "")
}

func (c *Config) GetMPDConnection() string {
	return c.Mpd.Server + ":" + c.Mpd.Port
}
