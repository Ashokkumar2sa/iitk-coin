package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var SigningKey = []byte("secretkeyfortoken")

func GenerateJWT(rollno string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = rollno
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	tokenString, err := token.SignedString(SigningKey)

	if err != nil {
		fmt.Println("Something Went Wrong ")
		return "", err
	}

	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func signup(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Running  ...")
	database, _ :=
		sql.Open("sqlite3", "./user_data.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, password TEXT)")
	statement.Exec()

	statement, _ =
		database.Prepare("INSERT INTO people (firstname, lastname, password) VALUES (?, ?, ?)")

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
			statement.Exec(name, rollno, hash)
			fmt.Println("Stored user and data")
		} else {
			fmt.Println("User rollno already registered.. Procede to login or if forgot password the reset it ")
		}
	}
	rows, _ :=
		database.Query("SELECT id, firstname, lastname, password FROM people")
	var id int
	var firstname string
	var lastname string
	var password string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname, &password)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname + " " + password)
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
				}
				fmt.Println(tokenString)
				f = 1
				break
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

func handleRequests() {

	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/secretpage", secretpage)

	http.ListenAndServe(":9000", nil)

}

func main() {
	handleRequests()
}
