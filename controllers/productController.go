package controllers

import (
	"Goland/database"
	"Goland/models"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func ShowProducts(res http.ResponseWriter, req *http.Request) {
	tpl, _ = tpl.ParseGlob("templates/*.html")
	tpl.ExecuteTemplate(res, "showProducts.html", GetProducts())
	return
}

func GetProducts() []models.Product {
	db = database.ConnectToDB()

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT name, description, price, quantity FROM products")
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Name, &product.Description, &product.Price, &product.Quantity); err != nil {
			return nil
		}
		products = append(products, product)
	}

	if err != nil {
		log.Fatal(err)
	}
	return products
}

func GetProduct(res http.ResponseWriter, req *http.Request) {
	db = database.ConnectToDB()

	name := req.FormValue("Target")

	fmt.Println(name)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT name, description, price, quantity FROM products")
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Name, &product.Description, &product.Price, &product.Quantity); err != nil {
			fmt.Println("error in rows")
		}
		products = append(products, product)
	}

	var result []models.Product
	for _, i := range products {
		if strings.Contains(strings.ToLower(i.Name), strings.ToLower(name)) || strings.ToLower(i.Name) == strings.ToLower(name) {
			result = append(result, i)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	tpl.ExecuteTemplate(res, "showProducts.html", result)
}

func AddProduct(res http.ResponseWriter, req *http.Request) {
	db = database.ConnectToDB()

	if req.Method != "POST" {
		http.ServeFile(res, req, "templates/add_product.html")
		return
	}

	description := req.FormValue("description")
	price := req.FormValue("price")
	quantity := req.FormValue("quantity")
	name := req.FormValue("name")

	_, err = db.Exec("INSERT INTO products(name, description, price, quantity) VALUES(?, ?, ?, ?)", name, description, price, quantity)
	if err != nil {
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	}
	http.Redirect(res, req, "/show_products", 301)
}
