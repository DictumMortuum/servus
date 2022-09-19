package config

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"io"
	"os"
)

type Config struct {
	PathTemplates string
	Timezone      string `config:"timezone"`
	Database      struct {
		Username string `config:"username"`
		Password string `config:"password"`
		Server   string `config:"server"`
		Port     string `config:"port"`
	} `config:"database"`
	Databases map[string]string
	Zerotier  struct {
		Token   string `config:"token"`
		Network string `config:"network"`
	} `config:"zerotier"`
	Mpd struct {
		Server string `config:"server"`
	} `config:"mpd"`
	Mealie struct {
		Token string `config:"token"`
	} `config:"mealie"`
	Telegram struct {
		Enabled bool     `config:"enabled" default:"false"`
		Token   string   `config:"token"`
		Users   []string `config:"users"`
	} `config:"telegram"`
	Calendar struct {
		Username string `config:"username"`
		Password string `config:"password"`
		Server   string `config:"server"`
	} `config:"calendar"`
	Atlas struct {
		ClientId string `config:"client_id"`
	} `config:"atlas"`
	Cache struct {
		Expiration int    `config:"expiration"`
		Cleanup    int    `config:"cleanup"`
		Path       string `config:"path"`
	} `config:"cache"`
	Deta struct {
		ProjectKey string `config:"project_key"`
	} `config:"deta"`
}

var (
	App Config
)

func Read(path string) error {
	mode := os.Getenv("GIN_MODE")
	path_templates := "templates/*"
	path_cfg := "servusrc.yml"

	if path != "" {
		path_cfg = path
	}

	if mode == "release" {
		gin.DisableConsoleColor()
		f, _ := os.Create("/var/log/servus.log")
		gin.DefaultWriter = io.MultiWriter(f)
		path_templates = "/usr/share/webapps/servus/*"
		path_cfg = "/etc/conf.d/servusrc.yml"
	}

	loader := confita.NewLoader(
		file.NewBackend(path_cfg),
	)

	err := loader.Load(context.Background(), &App)
	if err != nil {
		return err
	}

	App.PathTemplates = path_templates

	return nil
}
