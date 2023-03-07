package main

import (
	"Goland/controllers"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
)

var tpl *template.Template

func main() {
	tpl, _ = tpl.ParseGlob("templates/*.html")

	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/products", controllers.ShowProducts)
	http.HandleFunc("/search", controllers.GetProduct)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/add_product", controllers.AddProduct)
	http.HandleFunc("/", controllers.Home)
	http.ListenAndServe(":8081", nil)
}
