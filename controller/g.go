package controller

import (
	"html/template"
	"path/filepath"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var (
	templates map[string]*template.Template
	addr      string
	logger    = log.New()
)

func init() {
	setLogger()
	addr = viper.Get("server.address").(string) + ":" + viper.Get("server.port").(string)
}

func setLogger() {
	pathMap := lfshook.PathMap{
		logrus.DebugLevel: filepath.Join(basePath + "debug.log"),
		logrus.InfoLevel:  filepath.Join(basePath + "info.log"),
		logrus.WarnLevel:  filepath.Join(basePath + "warn.log"),
	}
	logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
}
