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
	var rowcount int
	database.QueryRow("SELECT COUNT(*) FROM `Results` WHERE `Data` != 'f'").Scan(&rowcount)
	/*
		Infact hold on,I might just grab the top ID count and work with that in batches of 10k
		that way I won't be using LIMIT (Known to be nearly the worse thing added to MySQL) and
		I will (hopefully) be rolling in performance
	*/
}
