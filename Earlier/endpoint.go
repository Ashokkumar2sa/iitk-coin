package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"

	_ "github.com/mattn/go-sqlite3"
)

//var passAdmin string = "awardpassword"
func balance(w http.ResponseWriter, req *http.Request) {
	database, _ :=
		sql.Open("sqlite3", "./user_data.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, num_coin INTEGER, password TEXT)")
	statement.Exec()

	fmt.Println("running...\n You have to be logged int to your accout \n enter the jwt token")
	var your_jwt_token string
	fmt.Scanln(&your_jwt_token)

	token, err := jwt.Parse(your_jwt_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return SigningKey, nil
	})

	if err != nil {
		fmt.Println("Please sign in with correct credentials")
	}

	if token.Valid {
		var rollno string
		fmt.Println("Provide your rollno again...")
		fmt.Scanln(&rollno)
		temp, _ :=
			database.Query("SELECT lastname , num_coin FROM people")
		var num_coin int
		var lastname string
		for temp.Next() {
			temp.Scan(&lastname, &num_coin)
			if rollno == lastname {
				fmt.Println(strconv.Itoa(num_coin) + " " + "coins")
				break
			}
		}
	}
}
func transfer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("running...\n You have to be logged int to your accout so enter the jwt token") //  enter the roll of user user to enter the password ")
	var your_jwt_token string
	fmt.Scanln(&your_jwt_token)

	token, err := jwt.Parse(your_jwt_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return SigningKey, nil
	})

	if err != nil {
		fmt.Println("Please sign in with correct credentials")
	}

	if token.Valid {
		database, _ :=
			sql.Open("sqlite3", "./user_data.db")
		statement, _ :=
			database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, num_coin INTEGER, password TEXT)")
		statement.Exec()

		var a1 int = 0
		var a2 int = 0
		fmt.Println("Enter your roll no ")
		var roll1 string
		fmt.Scanln(&roll1)
		fmt.Println("Enter the receiver roll no ")
		var roll2 string
		fmt.Scanln(&roll2)
		fmt.Println("Enter the amount make it sure less than availabe")
		var num int
		fmt.Scanln(&num)
		temp, _ :=
			database.Query("SELECT id, firstname, lastname ,num_coin , password FROM people")
		var num_coin int
		var id int
		var firstname string
		var lastname string
		var password string
		for temp.Next() {
			temp.Scan(&id, &firstname, &lastname, &num_coin, &password)
			if roll1 == lastname && a1 == 0 {
				a1 = 1
			}
			if roll2 == lastname && a2 == 0 {
				a2 = 1
			}
		}
		if a1 == 0 || a2 == 0 {
			fmt.Println("Please provide correct users... ")
		} else {
			a1 = 0
			a2 = 0
			temp, _ :=
				database.Query("SELECT  lastname ,num_coin  FROM people")
			var amu int
			var num_coin int
			var lastname string
			for temp.Next() {
				temp.Scan(&lastname, &num_coin)
				if roll1 == lastname && a1 == 0 {
					amu = num_coin - num
					statement, _ =
						database.Prepare("Update people set num_coin=? where lastname=? ")

					statement.Exec(amu, roll1)
					fmt.Println("Updated the coins for " + roll1)
				}

				if roll2 == lastname && a2 == 0 {
					amu = num_coin + num
					statement, _ =
						database.Prepare("Update people set num_coin=? where lastname=? ")

					statement.Exec(amu, roll2)
					fmt.Println("Updated the coins for " + roll2)
				}
			}
		}
	}
}

func award(w http.ResponseWriter, req *http.Request) {
	fmt.Println("running...\n You have to be logged int to your accout so enter the jwt token") //  enter the roll of user user to enter the password ")
	var your_jwt_token string
	fmt.Scanln(&your_jwt_token)

	token, err := jwt.Parse(your_jwt_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return SigningKey, nil
	})

	if err != nil {
		fmt.Println("Please sign in with correct credentials")
	}

	if token.Valid {
		database, _ :=
			sql.Open("sqlite3", "./user_data.db")
		statement, _ :=
			database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, num_coin INTEGER, password TEXT)")
		statement.Exec()

		var already int
		fmt.Println("Enter roll no of user....And the amout of coin")

		var roll string
		fmt.Scanln(&roll)
		fmt.Println("Enter the amount .. ")
		var count int
		fmt.Scanln(&count)
		temp, _ :=
			database.Query("SELECT id, firstname, lastname ,num_coin , password FROM people")
		var amo int
		var id int
		var num_coin int
		var firstname string
		var lastname string
		var password string
		for temp.Next() {
			temp.Scan(&id, &firstname, &lastname, &num_coin, &password)
			if roll == lastname {
				amo = num_coin + count
				statement, _ =
					database.Prepare("Update people set num_coin=? where lastname=? ")

				statement.Exec(amo, roll)
				fmt.Println("Updated the coins int this user account")
				already = 1

				rows, _ :=
					database.Query("SELECT id, firstname, lastname ,num_coin , password FROM people")
				var id int
				var amo int
				var firstname string
				var lastname string
				var password string
				for rows.Next() {
					rows.Scan(&id, &firstname, &lastname, &amo, &password)
					fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname + " " + strconv.Itoa(amo) + " " + " " + password)
				}
				break
			}
		}
		if already == 0 {
			fmt.Println(" User is not registered so not able to award coin... ")
		}
	}
}
func signup(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Running  ...")
	database, _ :=
		sql.Open("sqlite3", "./user_data.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, num_coin INTEGER, password TEXT)")
	statement.Exec()

	statement, _ =
		database.Prepare("INSERT INTO people (firstname, lastname, num_coin, password) VALUES (?, ?, ?, ?)")

	var coin_amu int = 0
	fmt.Println("Enter Your Name: ")
	var name string
	fmt.Scanln(&name)

	fmt.Println("Enter Your Roll no: ")
	var rollno string
	fmt.Scanln(&rollno)

	fmt.Println("Enter password")
	var pass string
	fmt.Scanln(&pass)
	if name == "" || rollno == "" || pass == "" {
		fmt.Println("Please enter a valid string")
	} else {

		var hash, _ = hashPassword(pass)

		var already int
		already = 0

		temp, _ :=
			database.Query("SELECT lastname FROM people")

		var lastname string
		for temp.Next() {
			temp.Scan(&lastname)
			if rollno == lastname {
				already = 1
				break
			}
		}
		if already == 0 {
			statement.Exec(name, rollno, coin_amu, hash)
			fmt.Println("Stored user and data")
		} else {
			fmt.Println("User rollno already registered.. Procede to login or if forgot password the reset it ")
		}
	}
	rows, _ :=
		database.Query("SELECT id, firstname, lastname ,num_coin , password FROM people")
	var id int
	var amo int
	var firstname string
	var lastname string
	var password string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname, &amo, &password)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname + " " + strconv.Itoa(amo) + " " + " " + password)
	}
}
func login(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Enter Roll No:")

	var rollno string
	fmt.Scanln(&rollno)
	fmt.Println("Enter Password")
	var pass string
	fmt.Scanln(&pass)

	database, _ :=
		sql.Open("sqlite3", "./user_data.db")
	rows, _ :=
		database.Query("SELECT lastname, password FROM people")

	var password string
	var lastname string
	var f int
	f = 0
	for rows.Next() {
		rows.Scan(&lastname, &password)
		if lastname == rollno {
			if doPasswordsMatch(password, pass) {
				fmt.Println("Logged in... Generating token")
				var tokenString string
				tokenString, err := GenerateJWT(rollno)
				if err != nil {
					fmt.Println("Failed to generate token")
				} else {
					fmt.Println(tokenString)
					f = 1
					break
				}
			}
		}
	}
	if f == 0 {
		fmt.Println("Invalid Credentials")
	}
}
func secretpage(w http.ResponseWriter, r *http.Request) {
	//if user is authorised that is claim is set to true  then we display his details
	fmt.Println("Enter the generated jwt token: ")
	var your_jwt_token string
	fmt.Scanln(&your_jwt_token)

	token, err := jwt.Parse(your_jwt_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return SigningKey, nil
	})

	if err != nil {
		fmt.Println("")
	}

	if token.Valid {
		fmt.Println("Super Secret Information displayed as correct jwt token is provided")

	}
}
