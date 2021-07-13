package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ :=
		sql.Open("sqlite3", "./Data.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()


	statement, _ =
		database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
        
	//taking input form a user 
	fmt.Println("Enter Your Name: ")
	var name string
	fmt.Scanln(&name)

	fmt.Println("Enter Your Roll no: ")
	var rollno string
	fmt.Scanln(&rollno)

	statement.Exec(name, rollno)

	rows, _ :=
		database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var c_name string
	var c_rollno string
	for rows.Next() {
		rows.Scan(&id, &c_name, &c_rollno)
		fmt.Println(strconv.Itoa(id) + ": " + c_name + " " + c_rollno)
	}
}
