package main

import (
	"net/http"
	"github.com/RakaiSeto/projectPraPKL/config"
	"html/template"
)

var tpl *template.Template

func init() {tpl = template.Must(template.ParseGlob("templates/*"))}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/loginForm", loginForm)
	http.HandleFunc("/loginProcess", config.LoginProcess)
	http.HandleFunc("/signupForm", signupForm)
	http.HandleFunc("/signupProcess", config.CreateUserProcess)
	http.HandleFunc("/userHome", config.UserHome)
	http.HandleFunc("/logout", config.Logout)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	l := config.IsAlreadyLogin(w, r)
	if l{
		http.Redirect(w, r, "/userHome", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/loginForm", http.StatusSeeOther)
	}
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func signupForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}