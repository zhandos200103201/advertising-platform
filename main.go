package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
)

var tpl *template.Template
var db *sql.DB

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")

	http.HandleFunc("/add_user", getAddUser)
	http.HandleFunc("/", getIndex)
	http.ListenAndServe("localhost:8081", nil)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
	return
}

func getAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "add_user.html", nil)
		r.ParseForm()
	} else if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		db, err := sql.Open("mysql", "zhandos:SAy#wm81j5AcM$Oy@tcp(localhost:3306)/go")
		if err != nil {
			fmt.Println("Error sql.Open")
			panic(err.Error())
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `users`(`id`, `name`, `email`, `phone`) VALUES ('%s', '%s', '%s', '%s')", id, name, email, phone))
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()
		http.Redirect(w, r, "/", 301)
	}
	return

}
