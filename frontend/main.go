package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/pmylund/go-cache"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Test if the database works.
	fmt.Println("Dominiator FrontEnd Server. Attempting DB connection")
	Database, e := GetDB()
	if e != nil {
		panic(e)
	}
	// Make a cache that all the general stats will be put in.
	cacheobj := cache.New(60*time.Minute, 1*time.Minute)

	fmt.Println("DB connection possible")
	Database.Exec("SHOW TABLES")
	// Okay so now we have a database connection.
	m := martini.Classic()
	m.Map(cacheobj) // ensure that the cache obj is delivered to each request
	m.Get("/api/search/:q", SearchForDomains)
	m.Get("/api/stats/", GetOverviewStats)
	m.Run()
}

func API2JSON(res http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.RequestURI, "/api") { // This causes anything with a /api prefix to have the content type of json.
		res.Header().Set("Content-Type", "application/json")
	}
}
