package main

import (
	"database/sql"
	"html/template"
	"net/http"
)

var db *sql.DB
var tpl *template.Template

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:passwor")
}

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}