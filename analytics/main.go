package main

import (
	"fmt"
)

func main() {
	fmt.Println("DB analyitcs started -- Connecting to Database")
	database, e := GetDB()
	if e != nil {
		fmt.Errorf("Oh dear, I was unable to connect to the database for this reason %s\n I will exit now.", e.Error())
		return
	}
	fmt.Println("DB connected. Getting the top of the ID stack")

	var rowcount int
	database.QueryRow("SELECT `id` FROM `Results` ORDER BY `id` DESC LIMIT 1").Scan(&rowcount)
	/*
		Infact hold on,I might just grab the top ID count and work with that in batches of 10k
		that way I won't be using LIMIT (Known to be nearly the worse thing added to MySQL) and
		I will (hopefully) be rolling in performance
	*/

	fmt.Printf("There are %d rows, I will start scanning though them 10,000 at a time", rowcount)
}
