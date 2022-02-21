package config

import (
	"net/http"
	"html/template"
)

var tpl *template.Template

func ExecuteTemplate(w http.ResponseWriter, tplName string, data interface{}) {
	tpl.ExecuteTemplate(w, tplName, data)
}