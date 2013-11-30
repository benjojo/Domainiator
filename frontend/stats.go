package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/pmylund/go-cache"
	"net/http"
	"strings"
)

func SearchForDomains(res http.ResponseWriter, req *http.Request, cache *cache.Cache, prams martini.Params) string {
	database, _ := GetDB()
	defer database.Close()
	var query string
	if prams["q"] != nil {
		http.Error(res, "No search query", http.StatusBadRequest)
		return ""
	}
	rows, _ := database.Query("SELECT Domain FROM `Domaniator`.`Results` WHERE Domain LIKE ? AND `Data` != 'f' LIMIT 10", prams["q"]+"%")
	resultsArray := make([]string, 0)
	defer rows.Close() // Ensure we don't leak connectctions
	for rows.Next() {
		var databack string
		err := rows.Scan(&databack)
		append(resultsArray, databack)
	}
	b, _ := json.Marshal(databack)
	return string(b)
}
