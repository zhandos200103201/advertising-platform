package database

import (
	"database/sql"
	"fmt"
)

var db *sql.DB
var err error

func ConnectToDB() *sql.DB {
	db, err = sql.Open("mysql", "zhandos:SAy#wm81j5AcM$Oy@/go")
	if err != nil {
		fmt.Println("Server could connect with database")
	}
	return db
}
