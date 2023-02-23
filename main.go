package main

import (
	"database/sql"
	"fmt"
	"html/template"
)
import _ "github.com/go-sql-driver/mysql"

import "golang.org/x/crypto/bcrypt"

import "net/http"

var db *sql.DB
var err error
var tpl *template.Template

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
	tpl.ExecuteTemplate(res, "index.html", nil)
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
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8081", nil)
}
