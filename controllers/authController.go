package controllers

import (
	"Goland/database"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"time"
)

var tpl *template.Template
var db *sql.DB
var err error

const SecretKey = "Hello"

func Signup(res http.ResponseWriter, req *http.Request) {
	db = database.ConnectToDB()
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

var Logerror string

func Login(res http.ResponseWriter, req *http.Request) {
	db = database.ConnectToDB()

	if req.Method != "POST" && Logerror != "" {
		tpl.ExecuteTemplate(res, "login.html", Logerror)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", 301)
		Logerror = "Dont have any user"
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
		Logerror = "Password is incorrect"
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    databaseUsername,
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(), //1 day
	})
	token, err := claims.SignedString([]byte(SecretKey))

	newCookie := http.Cookie{
		Name:     "jar",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		MaxAge:   99,
	}

	http.SetCookie(res, &newCookie)
	http.Redirect(res, req, "/", 301)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	newCookie := http.Cookie{
		Name:     "jar",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(res, &newCookie)
	http.Redirect(res, req, "/", 301)
}

func Home(res http.ResponseWriter, req *http.Request) {
	tpl, _ = tpl.ParseGlob("templates/*.html")

	username, _ := GetUser(req)
	tpl.ExecuteTemplate(res, "index.html", username)

}

func GetUser(req *http.Request) (string, error) {
	cookie, err := req.Cookie("jar")
	if err != nil {
		fmt.Println("Something was wrong")
		return "", err
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		fmt.Println("Something was wrong")
		return "", err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	fmt.Println(claims.Issuer)

	return claims.Issuer, nil
}
