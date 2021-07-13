package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("AshokKumar")

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Home Page  : direct to /signup, /login, /secretpage\n")
}

func secretPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, " You can view secret page as youa e logged\n")
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		http.NotFound(w, r)
		return
	}
	fmt.Println(r.URL.Path)
	switch r.Method {
	case "GET":
		w.Write([]byte("Signup Page!\nSend a POST request to signup\n"))

	case "POST":
		var newUser User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newUser)
		if err != nil {
			log.Printf("error decoding sakura response: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
		}
		log.Println(newUser)
		res := AddUser(newUser)
		w.Header().Set("Content-Type", "application/json")
		if res {
			json.NewEncoder(w).Encode("Successfully signuped")
		} else {
			json.NewEncoder(w).Encode("Failed signup")
		}

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}

}

func logIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}

	fmt.Println(r.URL.Path)
	switch r.Method {
	case "GET":
		w.Write([]byte("Login Page!\nPOST request to login \n"))

	case "POST":
		var loginRequest LoginRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&loginRequest)
		checkErr(err)
		log.Println(loginRequest)
		log.Println("User valid : ", UserValid(loginRequest))
		if UserValid(loginRequest) {
			expirationTime := time.Now().Add(15 * time.Minute)
			claims := &CustomClaims{
				Rollno: loginRequest.Rollno,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			checkErr(err)
			log.Println("Token is: ", tokenString)
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})
			w.Write([]byte("Successfully logged in!\n"))
		} else {
			w.Write([]byte("Invalid user credentials \n"))
		}

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}

}

func isLogin(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := cookie.Value
		claims := &CustomClaims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		endpoint(w, r)
	})
}

func balance(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/balance" {
		http.NotFound(w, r)
		return
	}

	fmt.Println(r.URL.Path)

	switch r.Method {
	case "GET":
		rollnos, ok := r.URL.Query()["rollno"]

		if !ok || len(rollnos[0]) < 1 {
			log.Println("Url Param 'rollno' is missing")
			return
		}

		rollno := rollnos[0]

		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := cookie.Value
		jwtClaims, _ := extractClaims(tknStr)
		if jwtClaims["rollno"] == rollno {
			coins := ReturnBalance(rollno)
			if coins >= 0 {
				w.Write([]byte("Rollno : " + rollno + "\n Balance : " + strconv.Itoa(int(coins)) + " coins\n"))
			} else {
				w.Write([]byte("User does not exist!\n"))
			}
		}
	case "POST":
		w.Write([]byte("Try Get Request.\n"))
	}
}

func reward(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/reward" {
		http.NotFound(w, r)
		return
	}

	fmt.Println(r.URL.Path)

	switch r.Method {
	case "GET":
		w.Write([]byte("Welcome to Reward Page!\nSend a POST request to award coins to user.\n"))

	case "POST":
		var rewardPayload RewardPayload
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&rewardPayload)
		checkErr(err)
		log.Println(rewardPayload)
		res := RewardMoney(rewardPayload)
		w.Header().Set("Content-Type", "application/json")
		if res {
			log.Printf("Coins awarded to rollno = %s , amounting = %d", rewardPayload.Rollno, rewardPayload.Coins)
			json.NewEncoder(w).Encode("Reward Success")
		} else {
			log.Printf("Reward coins failed")
			json.NewEncoder(w).Encode("Reward Failed")
		}

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

func transfer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/transfer" {
		http.NotFound(w, r)
		return
	}

	fmt.Println(r.URL.Path)

	switch r.Method {
	case "GET":
		w.Write([]byte("Welcome to Transfer Page!\nSend a POST request to tranfer coins peer to peer (P2P).\n"))

	case "POST":
		var transferPayload TransferPayload
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&transferPayload)
		checkErr(err)
		log.Println(transferPayload)
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tknStr := cookie.Value
		jwtClaims, _ := extractClaims(tknStr)
		if jwtClaims["rollno"] == transferPayload.SenderRollno {
			res := TransferCoins(transferPayload)
			w.Header().Set("Content-Type", "application/json")
			if res {
				log.Printf("Coins transfered from rollno = %s to rollno = %s amounting = %d", transferPayload.SenderRollno, transferPayload.ReceiverRollno, transferPayload.Coins)
				json.NewEncoder(w).Encode("Transfer Success")
			} else {
				log.Printf("Transaction failed!")
				json.NewEncoder(w).Encode("Transfer Failed")
			}
		}

	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
