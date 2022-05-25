package app

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

var dbSessions = map[string]session{} //sesID : uname + last Act
var foid string
var poid string
var prodid string
const sessionLength int = 1800
var g int //keep track of full order id, will always increment
var h int //keep track of product order id, will always increment

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

func createSession(w http.ResponseWriter, uname string, role string) {
	sID := uuid.NewV4()
	// if err != nil {
	// 	panic(err)
	// }
	c := &http.Cookie{
		Name: "session",
		Value: sID.String(),
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	// put role in redis
	err := rdb.HSet(ctx, "roledb", uname, role).Err()
	if err != nil {
		fmt.Println(err)
	}

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

func checkRole(uname string) string {
	query, err := rdb.HGet(ctx, "roledb", uname).Result()
	if err != nil {
		panic(err)
	}
	return query
}

func deleteRedisSession(uname string) {
	_ = rdb.HDel(ctx, "roledb", uname)
}

func cleanSessions() {
	for k, v := range dbSessions {
		if time.Since(v.lastActivity) > (time.Second * 1800) {
			delete(dbSessions, k)
		}
	}
}

func getNumber(variable string) (value int) {
	var temp int
	row := db.QueryRow("SELECT value FROM number WHERE type=$1", variable)
	err := row.Scan(&temp)
	if err != nil{
		panic(err)
	}

	return temp
}

func createCookie(w http.ResponseWriter, name string, value string, expire int) {
	c := &http.Cookie{
		Name: name,
		Value: value,
		Expires: time.Now().Add(time.Duration(expire) * time.Second),
	}
	http.SetCookie(w, c)
}