package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Attemping to connect to DB")
	database, _ := GetDB()
	// Okay so we are gonna grab a few really basic stats here.
	start := time.Now()
	// So the first bit we are going to get is the total done today.
	// SELECT COUNT(*) FROM Results WHERE `Timestamp` > timestampadd(hour, -24, now())
	var Total int

	database.QueryRow("SELECT COUNT(*) FROM Results WHERE `Data` != 'f'").Scan(&Total)
	// This one takes approx 10 mins :eek:

	var TotalFailed int
	database.QueryRow("SELECT COUNT(*) FROM Results WHERE `Data` = 'f'").Scan(&TotalFailed)

	// TODO: make it process the headers and get the avg content length.

	database.Exec("INSERT INTO `Domaniator`.`CachedResults` (`RequestCount`, `FailedCount`) VALUES (?, ?);", Total, TotalFailed)

	database.Close()
	elapsed := time.Since(start)
	fmt.Println("Done in %s", elapsed)
}
