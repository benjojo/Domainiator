package main

import (
	"encoding/json"
	// "fmt"
	"github.com/codegangsta/martini"
	"github.com/pmylund/go-cache"
	"net/http"
	// "strings"
)

func SearchForDomains(res http.ResponseWriter, req *http.Request, cache *cache.Cache, prams martini.Params) string {
	database, _ := GetDB()
	defer database.Close()
	if prams["q"] == "" {
		http.Error(res, "No search query", http.StatusBadRequest)
		return ""
	}
	rows, _ := database.Query("SELECT Domain FROM `Domaniator`.`Results` WHERE Domain LIKE ? AND `Data` != 'f' LIMIT 10", prams["q"]+"%")
	resultsArray := make([]string, 0)
	defer rows.Close() // Ensure we don't leak connectctions
	for rows.Next() {
		var databack string
		err := rows.Scan(&databack)
		if err != nil {
			http.Error(res, "Error reading from database", 500)
		}
		resultsArray = append(resultsArray, databack)
	}
	b, _ := json.Marshal(resultsArray)
	return string(b)
}

type StatsResponce struct {
	RequestCount   int64
	FailedCount    int64
	TopHeaders     string
	AvgContentSize int64
}

func GetOverviewStats(res http.ResponseWriter, req *http.Request, cache *cache.Cache, prams martini.Params) string {
	database, _ := GetDB()
	defer database.Close()

	var RequestCount int64
	var FailedCount int64
	var TopHeaders string
	var AvgContentSize int64
	row := database.QueryRow("SELECT  `RequestCount`,  `FailedCount`,  `TopHeaders`,  `AvgContentSize` FROM `Domaniator`.`CachedResults` ORDER BY `Day` DESC LIMIT 1")
	row.Scan(&RequestCount, &FailedCount, &TopHeaders, &AvgContentSize)

	Result := StatsResponce{
		RequestCount:   RequestCount,
		FailedCount:    FailedCount,
		TopHeaders:     TopHeaders,
		AvgContentSize: AvgContentSize,
	}
	b, _ := json.Marshal(Result)
	return string(b)
}
