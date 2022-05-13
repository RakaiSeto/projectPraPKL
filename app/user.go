package app

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

	// get uname
	c, _ := r.Cookie("session")
	s := dbSessions[c.Value].uname
	
	// delete customer with said uname from db
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

	// redirect to login
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	// get uname
	c, _ := r.Cookie("session")
	s := dbSessions[c.Value].uname

	// select all order with said uname and role == "active"
	rows, err := db.Query("SELECT * FROM fullOrder WHERE custid=$1 AND role=$2", s, "active")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		fmt.Println(err)
		return
	}
	defer rows.Close()
	
	// put those rows inside ords
	ords := make([]Order, 0)
	for rows.Next() {
		ord := Order{}
		err := rows.Scan(&ord.Id, &ord.Uname, &ord.Price, &ord.Created, &ord.Profit, &ord.Role)
		if err != nil {
			fmt.Println(err)
		}
		ords = append(ords, ord)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}
	
	// parse ords into template
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

	// get uname
	c, _ := r.Cookie("session")
	s := dbSessions[c.Value].uname
	 
	// increment id
	g = getNumber("foid")
	g++
	
	// update g in database
	_, err := db.Exec("UPDATE number SET value = $1 WHERE name=$2", g, "foid") 
	if err != nil {
  		panic(err)
	}

	// insert empty order into db
	_, err = db.Exec("INSERT INTO fullorder (id, custid, totalprice, profit, role) VALUES ($1, $2, $3, $4, $5)", g, s, 0, 0, "active")
	if err != nil {
  		panic(err)
	}

	// redirect to order list
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

	// put id into variable
	foid = r.FormValue("id")

	// select all product order with said foid and role == "active"
	rows, err := db.Query("SELECT * FROM productorder WHERE orderid=$1 AND role=$2", foid, "active")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	// put rows into pords
	pords := make([]Productorder, 0)	
	for rows.Next() {
		pord := Productorder{}
		err := rows.Scan(&pord.Id, &pord.Orderid, &pord.Procode, &pord.Qty, &pord.Discount, &pord.Poprice, &pord.Otherexp, &pord.Created, &pord.Otherdiscount, &pord.Role, &pord.Profit)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		pords = append(pords, pord)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}

	// parse pords into template
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

	// get id
	id := r.FormValue("id")

	// update role of order with said id
	_, err := db.Exec("UPDATE fullorder SET role = $1 WHERE id=$2", "deleted", id)
	if err != nil {
		panic(err)
	}
	
	// update role of any product order with said foid
	_, err = db.Exec("UPDATE productorder SET role=$1 WHERE orderid=$2", "deleted", id)
	if err != nil {
		panic(err)
	}
	
	// redirect to order list
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

	// get all product, sort it by code small to big
	rows, err := db.Query("SELECT * FROM product ORDER BY prodcode ASC")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		fmt.Println(err)
		return
	}
	defer rows.Close()

	// put rows into prods
	prods := make([]Product, 0)
	for rows.Next() {
		prod := Product{}
		err := rows.Scan(&prod.Prodcode, &prod.Name, &prod.Catprice)
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

	// parse prods into template
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

	// get form values
	prodcode := r.FormValue("prodcode")
	name := r.FormValue("name")
	catprice := r.FormValue("catprice")

	// insert product into database
	_, err := db.Exec("INSERT INTO product (prodcode, name, catprice) VALUES ($1, $2, $3)", prodcode, name, catprice)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	
	// redirect to product list
	http.Redirect(w, r, "/productList", http.StatusSeeOther)
}

func UpdateProductForm(w http.ResponseWriter, r *http.Request) {
	// save product id in var, parse template
	prodid = r.FormValue("code")
	tpl.ExecuteTemplate(w, "updateProduct.html", prodid)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get current product values
	curProd := Product{}
	row := db.QueryRow("SELECT * FROM product WHERE prodcode=$1", prodid)
	err := row.Scan(&curProd.Prodcode, &curProd.Name, &curProd.Catprice)
	if err == sql.ErrNoRows {
		http.Error(w, "No Product with that code", http.StatusBadRequest)
	}

	// get value from form
	name := r.FormValue("name")
	catprice := r.FormValue("catprice")

	// turn catprice to int
	catpriceConv, err := strconv.Atoi(catprice)
	if err != nil {
		panic(err)
	}

	// update values
	curProd.Name = name
	curProd.Catprice = catpriceConv

	// update in db
	_, err = db.Exec("UPDATE product SET name = $1, catprice = $2 WHERE prodcode = $3", curProd.Name, curProd.Catprice, curProd.Prodcode)
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

	// get code
	code := r.FormValue("code")

	// delete prooduct from database
	_, err := db.Exec("DELETE FROM fullorder WHERE prodcode=$1", code)
	if err != nil {
		panic(err)
	}

	// redirect to product list
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
	
	// increment id
	h = getNumber("poid")
	h++

	// update h in database
	_, err := db.Exec("UPDATE number SET value = $1 WHERE name=$2", h, "poid") 
	if err != nil {
  		panic(err)
	}

	// get product
	prodcode := r.FormValue("prodcode")
	prod := Product{}
	row := db.QueryRow("SELECT * FROM product WHERE prodcode=$1", prodcode)
	err = row.Scan(&prod.Prodcode, &prod.Name, &prod.Catprice)

	if err == sql.ErrNoRows {
		http.Error(w, "No Product with that code", http.StatusBadRequest)
		fmt.Println(err)
	}

	// get order id
	orderid := foid

	// convert qty to int from form
	qty := r.FormValue("qty")
	qtyConv, err := strconv.Atoi(qty)
	if err != nil {
		panic(err)
	}
	
	// convert curcat to int from form
	curcat := r.FormValue("curcat")
	curcatConv, err := strconv.Atoi(curcat)
	if err != nil {
		panic(err)
	}
	
	// count discount
	discount := (prod.Catprice - curcatConv)

	// get other discount and convert to int
	otherdisc := r.FormValue("otherdisc")
	otherdiscConv, err := strconv.Atoi(otherdisc)
	if err != nil {
		panic(err)
	}
	
	// count po price
	poprice := (prod.Catprice * qtyConv) - (discount * qtyConv) - otherdiscConv

	// get other expenses and convert to int
	otherexp := r.FormValue("otherexp")
	otherexpConv, err := strconv.Atoi(otherexp)
	if err != nil {
		panic(err)
	}
	
	// get current cpl price and convert to int
	currentcpl := r.FormValue("curcpl")
	currentcplConv, err := strconv.Atoi(currentcpl)
	if err != nil {
		panic(err)
	}

	// count profit
	profit := poprice - (currentcplConv * qtyConv) - otherexpConv

	// insert po into db 
	_, err = db.Exec("INSERT INTO productorder (id, orderid, prodcode, qty, discount, poprice, otherexp, otherdisc, role, profit) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", h,  orderid, prod.Prodcode, qty, discount, poprice, otherexp, otherdisc, "active", profit)
	if err != nil {
		panic(err)
	}

	// get fullorder with said foid
	fo := Order{}
	row = db.QueryRow("SELECT id, totalprice, profit FROM fullorder WHERE id=$1", foid)
	err = row.Scan(&fo.Id, &fo.Price, &fo.Profit)

	// update full order with said foid
	fo.Price += poprice
	fo.Profit += profit

	// reinsert to db
	_, err = db.Exec("UPDATE fullorder SET totalprice = $1, profit = $2 WHERE id=$3", fo.Price, fo.Profit, foid)


	// redirect
	url := "/seeOrder?id=" + foid

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func UpdateProductOrderForm(w http.ResponseWriter, r *http.Request) {
	// save poid in var, parse update template
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

	// get current product order number
	curPO := Productorder{}
	row := db.QueryRow("SELECT * FROM productorder WHERE id=$1", poid)
	err := row.Scan(&curPO.Id, &curPO.Orderid, &curPO.Procode, &curPO.Qty, &curPO.Discount, &curPO.Poprice, &curPO.Otherexp, &curPO.Created, &curPO.Otherdiscount, &curPO.Role, &curPO.Profit)
	if err != nil {
		fmt.Println(err)
	}
	
	// get current fullorder number
	curFO := Order{}
	row = db.QueryRow("SELECT * FROM fullorder WHERE id=$1", curPO.Orderid)
	err = row.Scan(&curFO.Id, &curFO.Uname, &curFO.Price, &curFO.Created, &curFO.Profit, &curFO.Role)
	if err != nil {
		fmt.Println(err)
	}

	// get product
	prodcode := r.FormValue("prodcode")
	prod := Product{}
	row = db.QueryRow("SELECT * FROM product WHERE prodcode=$1", prodcode)
	err = row.Scan(&prod.Prodcode, &prod.Name, &prod.Catprice)
	if err == sql.ErrNoRows {
		http.Error(w, "No Product with that code", http.StatusBadRequest)
		fmt.Println(err)
	}

	// convert qty to int from form
	qty := r.FormValue("qty")
	qtyConv, err := strconv.Atoi(qty)
	if err != nil {
		panic(err)
	}
	
	// convert curcat to int from form
	curcat := r.FormValue("curcat")
	curcatConv, err := strconv.Atoi(curcat)
	if err != nil {
		panic(err)
	}
	
	// count discount
	discount := (prod.Catprice - curcatConv)
	
	// get other discount and convert to int
	otherdisc := r.FormValue("otherdisc")
	otherdiscConv, err := strconv.Atoi(otherdisc)
	if err != nil {
		panic(err)
	}
	
	// count po price
	poprice := (prod.Catprice * qtyConv) - (discount * qtyConv) - otherdiscConv

	// get other expenses and convert to int
	otherexp := r.FormValue("otherexp")
	otherexpConv, err := strconv.Atoi(otherexp)
	if err != nil {
		panic(err)
	}
	
	// get current cpl price and convert to int
	currentcpl := r.FormValue("curcpl")
	currentcplConv, err := strconv.Atoi(currentcpl)
	if err != nil {
		panic(err)
	}
	
	// count profit
	profit := poprice - (currentcplConv * qtyConv) - otherexpConv

	// update product order in db
	_, err = db.Exec("UPDATE productorder SET prodcode = $1, qty = $2, discount = $3, poprice = $4, otherexp = $5, otherdisc = $6 WHERE id=$7", prodcode, qty, discount, poprice, otherexp, otherdisc, poid)
	if err != nil {
		panic(err)
	}
	
	// subtract current full order number with current product order number so it's back to before current product order is added
	curFO.Price -= curPO.Poprice
	curFO.Profit -= curPO.Profit
	
	// update full order number to after updated
	curFO.Price += poprice
	curFO.Profit += profit

	// update those number to FO database
	_, err = db.Exec("UPDATE fullorder SET totalprice = $1, profit = $2 WHERE id = $3", curFO.Price, curFO.Profit, curFO.Id)
	if err != nil {
		panic(err)
	}

	// redirect
	url := "/seeOrder?id=" + foid

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func DeleteProductOrder(w http.ResponseWriter, r *http.Request){
	if !IsAlreadyLogin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	code := r.FormValue("id")

	// get current product order number
	curPO := Productorder{}
	row := db.QueryRow("SELECT * FROM productorder WHERE id=$1", poid)
	err := row.Scan(&curPO.Id, &curPO.Orderid, &curPO.Procode, &curPO.Qty, &curPO.Discount, &curPO.Poprice, &curPO.Otherexp, &curPO.Created, &curPO.Otherdiscount, &curPO.Role, &curPO.Profit)
	if err != nil {
		fmt.Println(err)
	}
	
	// get current fullorder number
	curFO := Order{}
	row = db.QueryRow("SELECT * FROM fullorder WHERE id=$1", curPO.Orderid)
	err = row.Scan(&curFO.Id, &curFO.Uname, &curFO.Price, &curFO.Created, &curFO.Profit, &curFO.Role)
	if err != nil {
		fmt.Println(err)
	}

	// update current fullorder number
	curFO.Price -= curPO.Poprice
	curFO.Profit -= curPO.Profit

	// update current PO role
	_, err = db.Exec("UPDATE productorder SET role = $1 WHERE id=$2", "deleted", code)
	if err != nil {
		panic(err)
	}

	// update current FO in db
	_, err = db.Exec("UPDATE fullorder SET totalprice = $1, profit = $2 WHERE id=$3", curFO.Price, curFO.Profit, curPO.Orderid)
	if err != nil {
		panic(err)
	}

	// redirect
	url := "/seeOrder?id=" + fmt.Sprint(curPO.Orderid)

	http.Redirect(w, r, url, http.StatusSeeOther)
}