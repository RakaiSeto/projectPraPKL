package config

import (
	"database/sql"
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"html/template"
)

type user struct {
	Uname    string
	Password []byte
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("app/templates/*"))
}

func CreateUserProcess(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}

	// is taken?
	uname := req.FormValue("uname")
	if uname == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM customer WHERE uname = $1", uname)

	usr := user{}
	err := row.Scan(&usr.Uname, &usr.Password)
	if err == nil {
		http.Error(w, "Username already exist", http.StatusBadRequest)
		return
	} else if err != nil {
		fmt.Println(err)
	}

	// get form values
	password := req.FormValue("password")


	// validate
	if uname == "" || password == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// turn password to hash
	bs, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	fmt.Println(password)

	// insert user to DB
	_, err = db.Exec("INSERT INTO Customer (uname, password) VALUES ($1, $2)", uname, bs)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	
	fmt.Println("succesful")
	
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func LoginProcess(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// get form value
	uname := req.FormValue("uname")
	password := req.FormValue("password")
	fmt.Println(password)

	// validate
	if uname == "" || password == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// is username even there?
	row := db.QueryRow("SELECT * FROM customer WHERE uname = $1", uname)
	usr := user{}
	err := row.Scan(&usr.Uname, &usr.Password)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, req)
		http.Error(w, "Username do not match", http.StatusForbidden)
		return
	case err != nil:
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	fmt.Println(usr.Password)

	// compare
	err = bcrypt.CompareHashAndPassword(usr.Password, []byte(password))
	if err != nil {
		http.Error(w, "Password do not match", http.StatusForbidden)
		return
	}

	// create session
	createSession(w, req, uname)

	fmt.Println("login successful")

	http.Redirect(w, req, "/userHome", http.StatusSeeOther)
}

func UserHome(w http.ResponseWriter, r *http.Request) {
	l := IsAlreadyLogin(w, r)
	if !l {
		http.Redirect(w, r, "/loginForm", http.StatusSeeOther)
		return
	}
	c, err := r.Cookie("session")
	if err != nil {
		panic(err)
	}
	s := dbSessions[c.Value]
	UpdateLastActivity(w, r)
	tpl.ExecuteTemplate(w, "userHome.html", s.uname)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	c, _ := r.Cookie("session")
	// delete session from db
	delete(dbSessions, c.Value)
	// set cookie
	c = &http.Cookie{
		Name: "session",
		Value: "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean dbSessions
	cleanSessions()

	http.Redirect(w, r, "/loginForm", http.StatusSeeOther)
}