package user

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RakaiSeto/projectPraPKL/v2/db"
	proto "github.com/RakaiSeto/projectPraPKL/v2/proto"
)

var dbconn *sql.DB
var varError error

func init() {
	dbconn = db.Db
}

func AllUser() ([]*proto.User, error) {
	rows, err := dbconn.Query("SELECT id, uname, email, role FROM public.user")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	users := make([]*proto.User, 0)
	for rows.Next() {
		var user proto.User
		err := rows.Scan(&user.Id, &user.Uname, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func OneUser(id int) (*proto.User, error) {
	row := dbconn.QueryRow("SELECT id, uname, email, role FROM public.user where id=$1", id)

	user := proto.User{}
	err := row.Scan(&user.Id, &user.Uname, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func AddUser(user *proto.User) (*proto.AddUserStatus, error) {
	row := dbconn.QueryRow("SELECT id FROM public.user WHERE uname=$1", user.GetUname())
	var i int
	user.Role = "customer"
	err := row.Scan(&i)
	if i != 0 {
		return nil, err
	}
	
	_, err = dbconn.Exec("INSERT INTO public.user (uname, email, password, role) VALUES ($1, $2, $3, $4)", user.GetUname(), user.GetEmail(), user.GetPassword(), "customer")
	if err != nil {
		errString := err.Error()
		resp := proto.AddUserStatus{Response: "failed", Error: &errString}
		return &resp, sql.ErrConnDone
	}

	row = dbconn.QueryRow("SELECT id FROM public.user WHERE uname=$1", user.GetUname())
	err = row.Scan(&user.Id)
	if err != nil {
		errString := err.Error()
		resp := proto.AddUserStatus{Response: "failed", Error: &errString}
		return &resp, sql.ErrConnDone
	}

	resp := proto.AddUserStatus{Response: "success", User: user}

	return &resp, nil
}

func UpdateUser(user *proto.User) (*proto.ResponseStatus, error){
	QueryUser := proto.User{}
	row := dbconn.QueryRow("SELECT * from public.user where id = $1", user.Id)
	err := row.Scan(&QueryUser.Id, &QueryUser.Uname, &QueryUser.Email, &QueryUser.Password, &QueryUser.Role)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if user.Uname != "" {QueryUser.Uname = user.Uname}
	if user.Email != "" {QueryUser.Email = user.Email}
	if *user.Password != "" {QueryUser.Password = user.Password}

	_, err = dbconn.Exec("UPDATE public.user SET uname=$2, email=$3, password=$4 WHERE id=$1", QueryUser.Id, QueryUser.Uname, QueryUser.Email, QueryUser.Password)
	if err != nil {
		varError = err
		return nil, varError
	}

	return &proto.ResponseStatus{Response: "Success"}, nil
}

func DeleteUser(id int) (*proto.ResponseStatus, error) {
	row := dbconn.QueryRow("SELECT uname FROM public.user where id=$1", id)

	var name string
	
	err := row.Scan(&name)
	if err != nil {
		return nil, err
	} else if name == "" {
		varError = fmt.Errorf("user not found")
		return nil, varError
	}

	_, err = dbconn.Exec("DELETE FROM public.user WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	return &proto.ResponseStatus{Response: "Success"}, nil
}