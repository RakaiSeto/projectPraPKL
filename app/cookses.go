package config

import (
	"fmt"
	"net/http"
	"time"
	_ "github.com/lib/pq"

	uuid "github.com/satori/go.uuid"
)

type session struct {
	uname string
	lastActivity time.Time
}

var dbSessions = map[string]session{} //"uname" : sesID + last Act
var foid string
var poid string
var dbSessionsCleaned time.Time
const sessionLength int = 300
var g int //keep track of full order id, will always increment
var h int //keep track of product order id, will always increment


func init() {dbSessionsCleaned = time.Now()}

func IsAlreadyLogin(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	s, ok := dbSessions[c.Value]
	if ok {
		s.lastActivity = time.Now()
		dbSessions[c.Value] = s
	}
	
	// refresh session
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	return ok
}

func createSession(w http.ResponseWriter, r *http.Request, uname string) {
	sID, _ := uuid.NewV4()
	c := &http.Cookie{
		Name: "session",
		Value: sID.String(),
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	// insert to sesDB
	dbSessions[c.Value] = session{uname, time.Now()}
}

func UpdateLastActivity(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	s := dbSessions[c.Value]
	s.lastActivity = time.Now()
	dbSessions[c.Value] = s
}

func cleanSessions() {
	for k, v := range dbSessions {
		if time.Since(v.lastActivity) > (time.Second * 30) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now()
}

func getNumber(variable string) (value int) {
	var temp int
	row := db.QueryRow("SELECT value FROM number WHERE name=$1", variable)
	err := row.Scan(temp)
	if err != nil{
		panic(err)
	}

	return temp
}