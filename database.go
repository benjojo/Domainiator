package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB() (con *sql.DB, e error) {
	con, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/Domaniator")
	if err != nil {
		fmt.Println("[DB] An error happened in the setup of a SQL connection")
	}
	con.Ping()
	con.Exec("SET NAMES UTF8")
	return con, err
}
