package main

import (
	"html/template"
	"net/http"

	"github.com/RakaiSeto/projectPraPKL/app"
	"github.com/gin-gonic/gin"
)

var tpl *template.Template

func init() {tpl = template.Must(template.ParseGlob("app/templates/*.html"))}

func main() {
	router := gin.Default()
	router.GET("/user", app.GetAllAppusers)
	router.GET("/user/:uname", app.GetAppuserByName)
	router.POST("/user", app.PostAppuser)
	router.PATCH("/user/:uname", app.PatchAppuser)
	router.DELETE("/user/:uname", app.DeleteAppuser)
	http.HandleFunc("/loginForm", loginForm)
	http.HandleFunc("/loginProcess", app.LoginProcess)
	http.HandleFunc("/signupForm", signupForm)
	http.HandleFunc("/signupProcess", app.CreateUserProcess)
	http.HandleFunc("/adminHome", app.AdminHome)
	http.HandleFunc("/customerHome", app.CustomerHome)
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
	router.Run("localhost:8080")
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

