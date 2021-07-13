package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()
	http.HandleFunc("/", homePage)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/login", logIn)
	http.Handle("/secretpage", isLogin(secretPage))

	http.HandleFunc("/reward", reward)
	http.Handle("/transfer", isLogin(transfer))
	http.Handle("/balance", isLogin(balance))

	log.Println("Listen and Serve at 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
