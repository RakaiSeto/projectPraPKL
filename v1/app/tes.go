package app

import (
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
)

type AppuserJSON struct {
	Id    int `json:"id"`
	Uname    string `json:"uname"`
	Email    string `json:"email"`
	Role 	string `json:"role"`
}

type Error struct {
	Code int `json:"code"`
	Message error `json:"message"`

}
type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func GetAllAppusers(c *gin.Context) {
	rows, err := db.Query("SELECT id, uname, email, role FROM testable")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	users := make([]AppuserJSON, 0)
	for rows.Next() {
		user := AppuserJSON{}
		err := rows.Scan(&user.Id ,&user.Uname, &user.Email, &user.Role)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}
	c.IndentedJSON(http.StatusOK, users)
}

func GetAppuserByName(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow("SELECT uname, role FROM testable where uname=$1", id)

	user := AppuserJSON{}
	err := row.Scan(&user.Uname, &user.Role, &user.Id)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, user)
}

func PostAppuser(c *gin.Context) {
	var user Appuser

	if err := c.BindJSON(&user); err != nil {
        return
    }

	row := db.QueryRow("SELECT uname FROM testable WHERE uname=$1", user.Uname)
	usr := Appuser{}
	_ = row.Scan(&usr.Uname)
	if usr.Uname != "" {
		ResponseHandler(c, "User already exist, please use /update", http.StatusConflict)
		return
	}

	_, err := db.Exec("INSERT INTO testable (uname, password, role) VALUES ($1, $2, $3)", user.Uname, string(user.Password), user.Role)
	if err != nil {
		ErrorHandler(c, err, 500)
	}

	ResponseHandler(c, "success", 200)
}

func PatchAppuser(c *gin.Context) {
	uname := c.Param("uname")
	appuser1 := Appuser{}
	row := db.QueryRow("SELECT * from testable where uname =$1", uname)
	err := row.Scan(&appuser1.Uname, &appuser1.Password, &appuser1.Role)
	if err != nil {
		ErrorHandler(c, err, 500)
		panic(err)
	}

	var updateAppuser Appuser
	if err := c.ShouldBindJSON(&updateAppuser); err != nil {
		ErrorHandler(c, err, 500)
		return
	  }

	if updateAppuser.Uname != "" {appuser1.Uname = updateAppuser.Uname}
	if updateAppuser.Password != "" {appuser1.Password = updateAppuser.Password}
	if updateAppuser.Role != "" {appuser1.Role = updateAppuser.Role}

	_, err = db.Exec("UPDATE testable SET uname=$2, password=$3, role=$4 WHERE uname=$1", uname, appuser1.Uname, appuser1.Password, appuser1.Role)
	if err != nil {
		ErrorHandler(c, err, 500)
	}

	ResponseHandler(c, "success", 200)
}

func DeleteAppuser(c *gin.Context) {
	uname := c.Param("uname")
	_, err := db.Exec("DELETE FROM testable WHERE uname = $1", uname)
	if err != nil {
		ErrorHandler(c, err, 500)
	}
	ResponseHandler(c, "success", 200)
}

func ErrorHandler(c *gin.Context, message error, code int) {
	err := Error{}
	err.Code = code
	err.Message = message

	c.IndentedJSON(http.StatusInternalServerError, err)
}

func ResponseHandler(c *gin.Context, message string, code int) {
	response := Response{
		Code: code,
		Message: message,
	}
	c.IndentedJSON(http.StatusOK, response)
}