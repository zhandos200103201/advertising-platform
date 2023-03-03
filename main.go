package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"strings"
)
import _ "github.com/go-sql-driver/mysql"

import "golang.org/x/crypto/bcrypt"

import "net/http"

var db *sql.DB
var err error
var tpl *template.Template

func addProduct(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "templates/add_product.html")
		return
	}

	description := req.FormValue("description")
	price := req.FormValue("price")
	quantity := req.FormValue("quantity")
	name := req.FormValue("name")

	db, err = sql.Open("mysql", "zhandos:SAy#wm81j5AcM$Oy@/go")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO products(name, description, price, quantity) VALUES(?, ?, ?, ?)", name, description, price, quantity)
	if err != nil {
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	}
	http.Redirect(res, req, "/show_products", 301)
}

type Product struct {
	Name        string
	Description string
	Price       int
	Quantity    int
}

func getProducts() []Product {
	db, err = sql.Open("mysql", "zhandos:SAy#wm81j5AcM$Oy@/go")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT name, description, price, quantity FROM products")
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
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

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		tpl.ExecuteTemplate(res, "signup.html", nil)
		return
	}

	username := req.FormValue("username")
	password1 := req.FormValue("password1")
	password2 := req.FormValue("password2")
	email := req.FormValue("email")
	name := req.FormValue("name")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	if password1 == password2 {
		switch {
		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.", 500)
				return
			}

			_, err = db.Exec("INSERT INTO users(username, password, email, name) VALUES(?, ?, ?, ?)", username, hashedPassword, email, name)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.", 500)
				return
			}

			res.Write([]byte("User created!"))
			return
		case err != nil:
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		default:
			http.Redirect(res, req, "/", 301)
		}
	} else {
		http.Error(res, "Password doesn't match. Both passwords should be same!", 500)
		return
	}
}

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		tpl.ExecuteTemplate(res, "login.html", nil)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseusername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseusername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		fmt.Println("Dont have any user")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		fmt.Println("Password is incorrect")
		return
	}

	res.Write([]byte("Hello " + databaseusername))

}

func homePage(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "index.html", getProducts())
	return
}

func main() {
	tpl, _ = tpl.ParseGlob("templates/*.html")
	db, err = sql.Open("mysql", "zhandos:SAy#wm81j5AcM$Oy@/go")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/products", showProducts)
	http.HandleFunc("/search", getProduct)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/add_product", addProduct)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8081", nil)
}

func showProducts(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "showProducts.html", getProducts())
	return

}

func getProduct(res http.ResponseWriter, req *http.Request) {
	db, err = sql.Open("mysql", "zhandos:SAy#wm81j5AcM$Oy@/go")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	name := req.FormValue("Target")

	fmt.Println(name)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT name, description, price, quantity FROM products")
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Name, &product.Description, &product.Price, &product.Quantity); err != nil {
			fmt.Println("error in rows")
		}
		products = append(products, product)
	}

	var result []Product
	for _, i := range products {
		if strings.Contains(i.Name, name) || i.Name == name {
			result = append(result, i)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	tpl.ExecuteTemplate(res, "showProducts.html", result)
}
