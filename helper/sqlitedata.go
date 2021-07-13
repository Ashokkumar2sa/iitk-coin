package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var MyDB *sql.DB
var WalletDB *sql.DB
var TransactionHistoryDB *sql.DB

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func InitialiseDB() {
	MyDB, err := sql.Open("sqlite3", "./user-database-190349.db")
	checkErr(err)
	WalletDB, err := sql.Open("sqlite3", "./wallet-database-190349.db")
	checkErr(err)
	TransactionHistoryDB, err := sql.Open("sqlite3", "./transaction-database-190349.db")
	checkErr(err)
	defer MyDB.Close()
	defer WalletDB.Close()
	defer TransactionHistoryDB.Close()

	statement, err := MyDB.Prepare("CREATE TABLE IF NOT EXISTS User (rollno TEXT, name TEXT, password TEXT)")
	checkErr(err)
	log.Println("User Database opened and table created (if not existed) successfully!")
	statement.Exec()

	statement, err = WalletDB.Prepare("CREATE TABLE IF NOT EXISTS Wallet (rollno TEXT, coins INTEGER)")
	checkErr(err)
	log.Println("Wallet Database opened and table created (if not existed) successfully!")
	statement.Exec()

	statement, err = TransactionHistoryDB.Prepare("CREATE TABLE IF NOT EXISTS TransactionHistory (sender TEXT, receiver TEXT, coins INTEGER, remarks TEXT)")
	checkErr(err)
	log.Println("Wallet Database opened and table created (if not existed) successfully!")
	statement.Exec()
}

func AddUser(user models.User) bool {
	MyDB, err := sql.Open("sqlite3", "./user-database-190349.db")
	checkErr(err)
	WalletDB, err := sql.Open("sqlite3", "./wallet-database-190349.db")
	checkErr(err)
	defer MyDB.Close()
	defer WalletDB.Close()
	if !UserExists(user) {
		statement, err := MyDB.Prepare("INSERT INTO User (rollno, name, password) VALUES (?, ?, ?)")
		checkErr(err)
		statement.Exec(user.Rollno, user.Name, HashPwd(user.Password))

		statement, err = WalletDB.Prepare("INSERT INTO Wallet (rollno, coins) VALUES (?, ?)")
		checkErr(err)
		statement.Exec(user.Rollno, 0)

		log.Printf("New user details : rollno = %s, name = %s added in database user-database-190349.db\n", user.Rollno, user.Name)
		log.Printf("Wallet for user initiated in database wallet-database-190349.db\n")
		return true
	} else {
		log.Println("User with same roll no. already exists!")
		return false
	}
}

func UserValid(user models.LoginRequest) bool {
	MyDB, err := sql.Open("sqlite3", "./user-database-190349.db")
	checkErr(err)
	defer MyDB.Close()
	rows, err := MyDB.Query("SELECT * from User")
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

func UserExists(user models.User) bool {
	MyDB, err := sql.Open("sqlite3", "./user-database-190349.db")
	checkErr(err)
	defer MyDB.Close()
	err = MyDB.QueryRow("SELECT rollno FROM User WHERE rollno = ?", user.Rollno).Scan(&user.Rollno)
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
	WalletDB, err := sql.Open("sqlite3", "./wallet-database-190349.db")
	checkErr(err)
	defer WalletDB.Close()
	var coins int64
	err = WalletDB.QueryRow("SELECT coins FROM Wallet WHERE rollno = ?", rollno).Scan(&coins)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return -1
	}
	return coins
}

func RewardMoney(user models.RewardPayload) bool {
	WalletDB, err := sql.Open("sqlite3", "./wallet-database-190349.db")
	checkErr(err)
	defer WalletDB.Close()
	log.Printf("Updating the balance!")
	statement, err := WalletDB.Prepare("UPDATE Wallet SET coins = coins + ? WHERE rollno = ?")
	if err != nil {
		return false
	}
	statement.Exec(user.Coins, user.Rollno)

	TransactionHistoryDB, err := sql.Open("sqlite3", "./transaction-database-190349.db")
	checkErr(err)
	defer TransactionHistoryDB.Close()
	statement, err = TransactionHistoryDB.Prepare("INSERT INTO TransactionHistory (sender, receiver, coins, remarks) VALUES (?, ?, ?, ?)")
	checkErr(err)
	statement.Exec("000007", user.Rollno, user.Coins, "Reward")

	return true
}

func TransferCoins(user models.TransferPayload) bool {
	WalletDB, err := sql.Open("sqlite3", "./wallet-database-190349.db")
	checkErr(err)
	defer WalletDB.Close()

	ctx := context.Background()
	tx, err := WalletDB.BeginTx(ctx, nil)
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

	TransactionHistoryDB, err := sql.Open("sqlite3", "./transaction-database-190349.db")
	checkErr(err)
	defer TransactionHistoryDB.Close()
	statement, err := TransactionHistoryDB.Prepare("INSERT INTO TransactionHistory (sender, receiver, coins, remarks) VALUES (?, ?, ?, ?)")
	checkErr(err)
	statement.Exec(user.SenderRollno, user.ReceiverRollno, user.Coins, "Transfer")
	return true
}

func DeductTax(user models.TransferPayload) int64 {
	if (user.SenderRollno[0:2] == user.ReceiverRollno[0:2]) && (len(user.SenderRollno) == len(user.ReceiverRollno)) {
		return int64(float64(user.Coins) * 0.98)
	} else {
		return int64(float64(user.Coins) * 0.67)
	}
}

// func UpdateTransactionHistory() {
// 	TransactionHistoryDB, err := sql.Open("sqlite3", "./transaction-database-190349.db")
// 	checkErr(err)
// 	defer TransactionHistoryDB.Close()

// 	statement, err := TransactionHistoryDB.Prepare("INSERT INTO Transaction (sender, receiver, amount, status) VALUES (?, ?, ?, ?)")
// 	checkErr(err)
// 	statement.Exec(000007, user.Rollno, user.Coins, 0)
// }
