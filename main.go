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

	/*// need response from user
	resp, err := http.Get("https://localhost/8000")
	if err != nil {
		log.Fatalln(err)
	}*/

	statement, _ =
		database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	statement.Exec("Ashok", "190195")
	statement.Exec("User", "roll no")

	rows, _ :=
		database.Query("SELECT id, firstname, lastname FROM people")
	var id int
	var firstname string
	var lastname string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	}
}
