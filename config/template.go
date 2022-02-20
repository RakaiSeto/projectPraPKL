package config

import (
	"net/http"
	"html/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("../templates/*"))
}

func executeTemplate(w http.ResponseWriter, tplName string, data interface{}) {
	tpl.ExecuteTemplate(w, tplName, data)
}