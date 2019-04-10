package controller

import (
	"html/template"
	"log"
)

var (
	templates map[string]*template.Template
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	templates = populateTemplates()
}
