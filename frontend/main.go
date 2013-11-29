package main

import (
	// "fmt"
	"github.com/codegangsta/martini"
	// "net/http"
)

func main() {
	// Test if the database works.
	Database, e := GetDB()
	if e != nil {
		panic(e)
	}
	Database.Exec("SHOW TABLES")
	// Okay so now we have a database connection.
	m := martini.Classic()
	m.Run()
}
