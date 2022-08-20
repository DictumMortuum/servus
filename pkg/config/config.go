package config

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	PathTemplates string
	Timezone      string `yaml:"timezone"`
	Database      struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Server   string `yaml:"server"`
		Port     string `yaml:"port"`
	} `yaml:"database"`
	Databases map[string]string
	Zerotier  struct {
		Token   string `yaml:"token"`
		Network string `yaml:"network"`
	} `yaml:"zerotier"`
	Mpd struct {
		Server string `yaml:"server"`
		Port   string `yaml:"port"`
	}
	Mealie struct {
		Token string `yaml:"token"`
	} `yaml:"mealie"`
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
	Deta struct {
		ProjectKey string `yaml:"project_key"`
	} `yaml:"deta"`
}

var (
	App *Config
)

func Read(path string) error {
	mode := os.Getenv("GIN_MODE")
	path_templates := "templates/*"
	path_cfg := "servusrc"

	if path != "" {
		path_cfg = path
	}

	if mode == "release" {
		gin.DisableConsoleColor()
		f, _ := os.Create("/var/log/servus.log")
		gin.DefaultWriter = io.MultiWriter(f)
		path_templates = "/usr/share/webapps/servus/*"
		path_cfg = "/etc/servusrc"
	}

	file, err := os.Open(path_cfg)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	err = d.Decode(&App)
	if err != nil {
		return err
	}

	App.PathTemplates = path_templates

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
