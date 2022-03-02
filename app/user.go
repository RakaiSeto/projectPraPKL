package config

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
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
		return
	}

	// insert user to DB
	_, err = db.Exec("INSERT INTO customer (uname, password) VALUES ($1, $2)", uname, bs)
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	c, _ := r.Cookie("session")

	s := dbSessions[c.Value].uname
	statement := `DELETE FROM customer WHERE uname =$1`
	_, err := db.Exec(statement, s)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, `we have encountered an error 
		%v
		please try again later`, err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

}

func OrderList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c, _ := r.Cookie("session")
	s := dbSessions[c.Value].uname

	rows, err := db.Query("SELECT * FROM fullOrder WHERE custid=$1", s)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		fmt.Println(err)
		return
	}
	defer rows.Close()
	
	ords := make([]Order, 0)
	for rows.Next() {
		ord := Order{}
		err := rows.Scan(&ord.Id, &ord.Uname, &ord.Price)
		if err != nil {
			fmt.Println(err)
		}
		prof := ((ord.Price / 100) * 23)
		ord.Profit = prof
		ords = append(ords, ord)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}
	
	tpl.ExecuteTemplate(w, "orderList.html", ords)
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	c, _ := r.Cookie("session")
	s := dbSessions[c.Value].uname
	 
	if g == 0 {
		g = 1 
	} else {
		g++
	}

	sqlStatement := `
	INSERT INTO fullorder (id, custid, totalprice)
	VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, g, s, 0)
	if err != nil {
  		panic(err)
	}

	http.Redirect(w, r, "/orderList", http.StatusSeeOther)
}

func SeeOrder(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	foid = r.FormValue("id")

	rows, err := db.Query("SELECT * FROM productorder WHERE orderid=$1", foid)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	pords := make([]Productorder, 0)	
	for rows.Next() {
		pord := Productorder{}
		err := rows.Scan(&pord.Id, &pord.Orderid, &pord.Prodid, &pord.Qty, &pord.Otherdiscount, &pord.Poprice)
		if err != nil {
			fmt.Println(err)
			return
		}
		beDisc := pord.Poprice + pord.Otherdiscount
		pord.Beforediscount = beDisc
		pords = append(pords, pord)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}

	tpl.ExecuteTemplate(w, "productOrderList.html", pords)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	_, err := db.Exec("DELETE FROM fullorder WHERE id=$1", id)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/orderList", http.StatusSeeOther)
}

// PRODUCT SECTION

func ProductList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rows, err := db.Query("SELECT * FROM product")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		fmt.Println(err)
		return
	}
	defer rows.Close()

	prods := make([]Product, 0)
	for rows.Next() {
		prod := Product{}
		err := rows.Scan(&prod.Prodcode, &prod.Name, &prod.Catprice, &prod.Memprice)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		prods = append(prods, prod)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpl.ExecuteTemplate(w, "productList.html", prods)
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	prodcode := r.FormValue("prodcode")
	name := r.FormValue("name")
	catprice := r.FormValue("catprice")
	memprice := r.FormValue("memprice")

	_, err := db.Exec("INSERT INTO product (prodcode, name, catprice, memprice) VALUES ($1, $2, $3, $4)", prodcode, name, catprice, memprice)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	
	fmt.Println("succesful")
	http.Redirect(w, r, "/productList", http.StatusSeeOther)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	code := r.FormValue("code")
	_, err := db.Exec("DELETE FROM fullorder WHERE prodcode=$1", code)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/productList", http.StatusSeeOther)
}

func AddProductOrder(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	
	if h == 0 {
		h = 1
	} else {
		h++
	}

	prodcode := r.FormValue("prodcode")
	prod := Product{}
	row := db.QueryRow("SELECT * FROM product WHERE prodcode=$1", prodcode)
	err := row.Scan(&prod.Prodcode, &prod.Name, &prod.Catprice, &prod.Memprice)

	if err == sql.ErrNoRows {
		http.Error(w, "No Product with that code", http.StatusBadRequest)
		fmt.Println(err)
	}

	orderid := foid
	qty := r.FormValue("qty")
	t, err := strconv.Atoi(qty)
	if err != nil {
		panic(err)
	}
	beforediscount := (prod.Catprice * t)
	discount := r.FormValue("otherdisc")
	t, err = strconv.Atoi(discount)
	if err != nil {
		panic(err)
	}
	poprice := beforediscount - t

	_, err = db.Exec("INSERT INTO productorder (id, orderid, procode, qty, discount, poprice) VALUES ($1, $2, $3, $4, $5, $6)", h, orderid, prodcode, qty, discount, poprice)
	if err != nil {
		panic(err)
	}

	row = db.QueryRow("SELECT * FROM fullorder WHERE id=$1", orderid)

	fo := Order{}
	err = row.Scan(&fo.Id, &fo.Uname, &fo.Price, &fo.Profit)
	if err != nil {
		panic(err)
	}

	fo.Price = fo.Price + poprice

	_, err = db.Exec("UPDATE fullorder SET totalprice=$1 WHERE id=$2", fo.Price, fo.Id)
	if err != nil {
		panic(err)
	}

	te := foid
	re := `/seeOrder?id=` + te

	http.Redirect(w, r, re, http.StatusSeeOther)
}

func UpdateProductOrderForm(w http.ResponseWriter, r *http.Request) {
	poid = r.FormValue("id")
	tpl.ExecuteTemplate(w, "updateProductOrder.html", nil)
}

func UpdateProductOrder(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}


	id := poid
	prodcode := r.FormValue("prodcode")
	prod := Product{}
	row := db.QueryRow("SELECT * FROM product WHERE prodcode=$1", prodcode)
	err := row.Scan(&prod.Prodcode, &prod.Name, &prod.Catprice, &prod.Memprice)

	if err == sql.ErrNoRows {
		http.Error(w, "No Product with that code", http.StatusBadRequest)
		fmt.Println(err)
	}

	orderid := foid
	qty := r.FormValue("qty")
	t, err := strconv.Atoi(qty)
	if err != nil {
		panic(err)
	}
	beforediscount := (prod.Catprice * t)
	discount := r.FormValue("otherdisc")
	t, err = strconv.Atoi(discount)
	if err != nil {
		panic(err)
	}
	poprice := beforediscount - t

	_, err = db.Exec("UPDATE productorder SET orderid = $1, prodcode = $2, qty = $3, discount = $4, poprice = $5 WHERE id=$6", orderid, prodcode, qty, discount, poprice, id)
	if err != nil {
		panic(err)
	}

	te := foid
	re := `/seeOrder?id=` + te

	http.Redirect(w, r, re, http.StatusSeeOther)
}

func DeleteProductOrder(w http.ResponseWriter, r *http.Request){
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	code := r.FormValue("id")

	_, err := db.Exec("DELETE FROM productorder WHERE id=$1", code)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/productOrderList", http.StatusSeeOther)
}