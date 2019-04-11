package controller

import (
	"html/template"
	"log"

	"github.com/spf13/viper"
)

var (
	templates map[string]*template.Template
	addr      string
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	addr = viper.Get("server.address").(string) + ":" + viper.Get("server.port").(string)
}
