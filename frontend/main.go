package main

import (
	// "fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// Test if the database works.
	Database, e := GetDB()
	if e != nil {
		panic(e)
	}
	// Okay so now we have a database connection.
	m := martini.Classic()
	m.Run()
}
