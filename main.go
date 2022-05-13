package main

import (
	"net/http"
	"github.com/RakaiSeto/projectPraPKL/app"
	"html/template"
)

var tpl *template.Template

func init() {tpl = template.Must(template.ParseGlob("app/templates/*.html"))}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/loginForm", loginForm)
	http.HandleFunc("/loginProcess", app.LoginProcess)
	http.HandleFunc("/signupForm", signupForm)
	http.HandleFunc("/signupProcess", app.CreateUserProcess)
	http.HandleFunc("/userHome", app.UserHome)
	http.HandleFunc("/orderList", app.OrderList)
	http.HandleFunc("/addOrder", app.AddOrder)
	http.HandleFunc("/seeOrder", app.SeeOrder)
	http.HandleFunc("/deleteOrder", app.DeleteOrder)
	http.HandleFunc("/addProductOrderForm", addProductOrderForm)
	http.HandleFunc("/addProductOrder", app.AddProductOrder)
	http.HandleFunc("/updateProductOrder", app.UpdateProductOrderForm)
	http.HandleFunc("/updateProductOrderProcess", app.UpdateProductOrder)
	http.HandleFunc("/deleteProductOrder", app.DeleteProductOrder)
	http.HandleFunc("/productList", app.ProductList)
	http.HandleFunc("/addProductForm", addProduct)
	http.HandleFunc("/addProduct", app.AddProduct)
	http.HandleFunc("/updateProduct", app.UpdateProductForm)
	http.HandleFunc("/updateProductProcess", app.UpdateProduct)
	http.HandleFunc("/deleteProduct", app.DeleteProduct)
	http.HandleFunc("/logout", app.Logout)
	http.HandleFunc("/deleteUser", app.DeleteUser)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	l := app.IsAlreadyLogin(w, r)
	if l{
		http.Redirect(w, r, "/userHome", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/loginForm", http.StatusSeeOther)
	}
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func signupForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

func addProductOrderForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "addProductOrder.html", nil)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "addProduct.html", nil)
}

