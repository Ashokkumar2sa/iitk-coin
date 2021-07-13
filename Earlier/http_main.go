package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func handleRequests() {

	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/secretpage", secretpage)
	http.HandleFunc("/award", award)
	http.HandleFunc("/balance", balance)
	http.HandleFunc("/transfer", transfer)
	http.ListenAndServe(":9000", nil)

}

func main() {
	handleRequests()
}
