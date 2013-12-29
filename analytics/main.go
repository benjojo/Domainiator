package main

import (
	"fmt"
	"time"
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
	fmt.Printf("There are %d rows, I will start scanning though them 10,000 at a time\n", rowcount)
	// Now I need to test how long this is going to take by doing the most legit way of testing this...
	prestart := time.Now()

	// Do the complex query mid way though the result set.
	// This will give us a rough idea how long it will take...
	database.QueryRow("SELECT SUM(LENGTH(`Data`)) FROM `Results` WHERE id > (94605236/2) AND id < (94605236/2) + 1001")

	timetaken := time.Since(prestart)

	fmt.Printf("Highly optomistic estimation is %f mins or %f Hours\n", timetaken.Minutes()*float64(rowcount/1000), timetaken.Hours()*float64(rowcount/1000))

}
