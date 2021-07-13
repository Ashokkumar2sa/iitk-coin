package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func InitDB() {
	userdb, err := sql.Open("sqlite3", "./user_data.db")
	checkErr(err)
	Walletdb, err := sql.Open("sqlite3", "./wallet_data.db")
	checkErr(err)
	TransHistorydb, err := sql.Open("sqlite3", "./trans_data.db")
	checkErr(err)
	defer userdb.Close()
	defer Walletdb.Close()
	defer TransHistorydb.Close()

	statement, err := userdb.Prepare("CREATE TABLE IF NOT EXISTS User (rollno TEXT, name TEXT, password TEXT)")
	checkErr(err)
	log.Println("User Database created ")
	statement.Exec()

	statement, err = Walletdb.Prepare("CREATE TABLE IF NOT EXISTS Wallet (rollno TEXT, coins INTEGER)")
	checkErr(err)
	log.Println("Wallet Database created ")
	statement.Exec()

	statement, err = TransHistorydb.Prepare("CREATE TABLE IF NOT EXISTS TransactionHistory (sender TEXT, receiver TEXT, coins INTEGER, remarks TEXT)")
	checkErr(err)
	log.Println("Transaction History Database ")
	statement.Exec()
}

func AddUser(user User) bool {
	userdb, err := sql.Open("sqlite3", "./user_data.db")
	checkErr(err)
	Walletdb, err := sql.Open("sqlite3", "./wallet_data.db")
	checkErr(err)
	defer userdb.Close()
	defer Walletdb.Close()
	if !UserExists(user) {
		statement, err := userdb.Prepare("INSERT INTO User (rollno, name, password) VALUES (?, ?, ?)")
		checkErr(err)
		statement.Exec(user.Rollno, user.Name, HashPwd(user.Password))

		statement, err = Walletdb.Prepare("INSERT INTO Wallet (rollno, coins) VALUES (?, ?)")
		checkErr(err)
		statement.Exec(user.Rollno, 0)

		log.Printf("New user details : rollno = %s, name = %s \n ", user.Rollno, user.Name)
		log.Printf("Wallet for user initiated\n")
		return true
	} else {
		log.Println("User with same roll no. already exists!")
		return false
	}
}

func UserValid(user LoginRequest) bool {
	userdb, err := sql.Open("sqlite3", "./user_data.db")
	checkErr(err)
	defer userdb.Close()
	rows, err := userdb.Query("SELECT * from User")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var rollno string
		var name string
		var password string
		err = rows.Scan(&rollno, &name, &password)
		checkErr(err)
		if user.Rollno == rollno && CheckPasswords(password, user.Password) {
			return true
		}
	}

	return false
}

func UserExists(user User) bool {
	userdb, err := sql.Open("sqlite3", "./user_data.db")
	checkErr(err)
	defer userdb.Close()
	err = userdb.QueryRow("SELECT rollno FROM User WHERE rollno = ?", user.Rollno).Scan(&user.Rollno)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return true
}

func HashPwd(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}

func CheckPasswords(hashedPwd string, pwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(pwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func ReturnBalance(rollno string) int64 {
	Walletdb, err := sql.Open("sqlite3", "./wallet_data.db")
	checkErr(err)
	defer Walletdb.Close()
	var coins int64
	err = Walletdb.QueryRow("SELECT coins FROM Wallet WHERE rollno = ?", rollno).Scan(&coins)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return -1
	}
	return coins
}

func RewardMoney(user RewardPayload) bool {
	Walletdb, err := sql.Open("sqlite3", "./wallet_data.db")
	checkErr(err)
	defer Walletdb.Close()
	log.Printf("Updating the balance!")
	statement, err := Walletdb.Prepare("UPDATE Wallet SET coins = coins + ? WHERE rollno = ?")
	if err != nil {
		return false
	}
	statement.Exec(user.Coins, user.Rollno)

	TransHistorydb, err := sql.Open("sqlite3", "./trans_data.db")
	checkErr(err)
	defer TransHistorydb.Close()
	statement, err = TransHistorydb.Prepare("INSERT INTO TransactionHistory (sender, receiver, coins, remarks) VALUES (?, ?, ?, ?)")
	checkErr(err)
	statement.Exec("000007", user.Rollno, user.Coins, "Reward")

	return true
}

func TransferCoins(user TransferPayload) bool {
	Walletdb, err := sql.Open("sqlite3", "./wallet_data.db")
	checkErr(err)
	defer Walletdb.Close()

	ctx := context.Background()
	tx, err := Walletdb.BeginTx(ctx, nil)
	checkErr(err)

	user.Coins = DeductTax(user)

	res, err := tx.ExecContext(ctx, "UPDATE Wallet SET coins = coins - ? WHERE rollno=? AND coins - ? >= 0", user.Coins, user.SenderRollno, user.Coins)
	checkErr(err)
	rows_affected, err := res.RowsAffected()
	checkErr(err)

	if rows_affected != 1 {
		tx.Rollback()
		return false
	}

	res, err = tx.ExecContext(ctx, "UPDATE Wallet SET coins = coins + ? WHERE rollno=?", user.Coins, user.ReceiverRollno, user.Coins)
	checkErr(err)
	rows_affected, err = res.RowsAffected()
	checkErr(err)

	if rows_affected != 1 {
		tx.Rollback()
		return false
	}

	err = tx.Commit()
	checkErr(err)

	TransHistorydb, err := sql.Open("sqlite3", "./trans_data.db")
	checkErr(err)
	defer TransHistorydb.Close()
	statement, err := TransHistorydb.Prepare("INSERT INTO TransactionHistory (sender, receiver, coins, remarks) VALUES (?, ?, ?, ?)")
	checkErr(err)
	statement.Exec(user.SenderRollno, user.ReceiverRollno, user.Coins, "Transfer")
	return true
}

func DeductTax(user TransferPayload) int64 {
	if (user.SenderRollno[0:2] == user.ReceiverRollno[0:2]) && (len(user.SenderRollno) == len(user.ReceiverRollno)) {
		return int64(float64(user.Coins) * 0.98)
	} else {
		return int64(float64(user.Coins) * 0.67)
	}
}
