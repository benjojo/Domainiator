package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB() (con *sql.DB, e error) {
	con, err := sql.Open("mysql", *databasestring)
	if err != nil {
		fmt.Println("[DB] An error happened in the setup of a SQL connection")
		return con, err
	}
	con.Ping()
	con.Exec("SET NAMES UTF8")
	return con, err
}
