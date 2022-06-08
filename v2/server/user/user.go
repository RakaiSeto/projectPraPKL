package user

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/RakaiSeto/projectPraPKL/v2/db"
	proto "github.com/RakaiSeto/projectPraPKL/v2/server/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var dbconn *sql.DB

func init() {
	dbconn = db.Db
}

func AllUser() ([]*proto.User, error) {
	rows, err := dbconn.Query("SELECT id, uname, email, role FROM public.user")
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}
	defer rows.Close()

	users := make([]*proto.User, 0)
	for rows.Next() {
		var user proto.User
		err := rows.Scan(&user.Id, &user.Uname, &user.Email, &user.Role)
		if err != nil {
			return nil, status.Error(codes.Code(2), err.Error())
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
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}
	return &user, nil
}

func AddUser(user *proto.User) (*proto.AddUserStatus, error) {
	row := dbconn.QueryRow("SELECT id FROM public.user WHERE uname=$1", user.GetUname())
	var i int
	user.Role = "customer"
	err := row.Scan(&i)
	if i != 0 {
		return nil, status.Error(codes.Code(6), "user already exists, please change its uname")
	}
	
	_, err = dbconn.Exec("INSERT INTO public.user (uname, email, password, role) VALUES ($1, $2, $3, $4)", user.GetUname(), user.GetEmail(), user.GetPassword(), "customer")
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	row = dbconn.QueryRow("SELECT id FROM public.user WHERE uname=$1", user.GetUname())
	err = row.Scan(&user.Id)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	resp := proto.AddUserStatus{Response: "success", User: user}

	return &resp, nil
}

func UpdateUser(user *proto.User) (*proto.ResponseStatus, error){
	QueryUser := proto.User{}
	row := dbconn.QueryRow("SELECT * from public.user where id = $1", user.Id)
	err := row.Scan(&QueryUser.Id, &QueryUser.Uname, &QueryUser.Email, &QueryUser.Password, &QueryUser.Role)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set"){
			return nil, status.Error(codes.Code(5), "user not found")
		}
		return nil, status.Error(codes.Code(2), err.Error())
	}

	if user.Uname != "" {QueryUser.Uname = user.Uname}
	if user.Email != "" {QueryUser.Email = user.Email}
	if *user.Password != "" {QueryUser.Password = user.Password}

	_, err = dbconn.Exec("UPDATE public.user SET uname=$2, email=$3, password=$4 WHERE id=$1", QueryUser.Id, QueryUser.Uname, QueryUser.Email, QueryUser.Password)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	return &proto.ResponseStatus{Response: "Success"}, nil
}

func DeleteUser(inputUser *proto.User) (*proto.ResponseStatus, error) {
	row := dbconn.QueryRow("SELECT uname, password FROM public.user where id=$1", inputUser.Id)

	var name string
	var password string
	
	err := row.Scan(&name, &password)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	if name == ""{
		return nil, status.Error(codes.Code(5), "user not found")
	}
	
	if inputUser.GetPassword() != password {
		if inputUser.GetPassword() == "" {
			return nil, errors.New("please include password in request")
		}
		return nil, errors.New("wrong password for user")
	}

	_, err = dbconn.Exec("DELETE FROM public.user WHERE id=$1", inputUser.Id)
	if err != nil {
		return nil, status.Error(codes.Code(2), err.Error())
	}

	return &proto.ResponseStatus{Response: "Success"}, nil
}


