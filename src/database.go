package main

import (
	"database/sql"
_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"time"
)

var dbMu sync.Mutex

/**
Initialize db tables with file path
@param String filepath
@Returns *sql.DB
*/
func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	log.Println(filepath)
	if err != nil {
		log.Print(err)
	}
	if db == nil {
		log.Print("db nil")
	}
	log.Println("Successfull opened db")
	return db
}

/**
Creates the user table with the following paramters:
IP: String
Name: string
*/
func CreateUserTable() {
	sql_table := `
	CREATE TABLE IF NOT EXISTS Users(
		IP TEXT ,
		Username TEXT PRIMARY KEY,
		Pass TEXT,
		SessionID TEXT ,
    DateCreated
	);
	`
	_, err := db.Exec(sql_table)
	if err != nil {
		log.Print(err)
	}
}

func CheckValidSessionToken(sessionToken string) bool {
	sql_stmt := `SELECT SessionID FROM Users`
	dbMu.Lock()
	rows, err := db.Query(sql_stmt)
	dbMu.Unlock()
	if err != nil {
		log.Println(" No Results in database", err.Error())
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var SessionToken string
		if err := rows.Scan(&SessionToken); err != nil {
			log.Println("Error scanning databasse for sessiond id", err)
			return false
		}
		if sessionToken == SessionToken {
			return false
		}
	}
	log.Println("Valid Session Unique")
	return true
}

func StoreUserInfo(socketClientIP string, Username string, Password string, SessionID string) {
	sql_stmt := `
	INSERT OR REPLACE INTO Users(
		IP,
		Username,
		Pass,
		SessionID,
    DateCreated
	)values(?, ?, ?, ?, ?)
	`
	stmt, err := db.Prepare(sql_stmt)
	if err != nil {
		log.Print(err)
	}
	c := User{
		IP:          socketClientIP,
		Username:    Username,
		Password:    Password,
		SessionID:   SessionID,
		DateCreated: time.Now().Unix(),
	}
	if _, err := stmt.Exec(c.IP, c.Username, c.Password, c.SessionID, c.DateCreated); err != nil {
		log.Println(err)
	}
	log.Println("Store New User Info")
}

func getUserInfo(socketClientIP string) (string, error) {
	var ip string
	sql_stmt := "SELECT Username FROM Users WHERE IP = $1"
	if err := db.QueryRow(sql_stmt, socketClientIP).Scan(&ip); err != nil {
		return "", err
	}
	return ip, nil
}

func getUsername(sessionToken string) (string, error) {
	var Username string
	sql_stmt := `SELECT Username FROM Users WHERE SessionID = $1`
	if err := db.QueryRow(sql_stmt, sessionToken).Scan(&Username); err != nil {
		return "", err
	}
	return Username, nil
}
func storeNewSessionToken(sid string, Username string) {
	sql_stmt := `UPDATE Users SET SessionID = $1 WHERE Username = $2`
		dbMu.Lock()
	if _, err := db.Exec(sql_stmt, sid, Username); err != nil {
			dbMu.Unlock()
		log.Println("Error in storing sessionToken: ", err)
		return
	}
	log.Println("Stored SessionToken")
	return
}
func getUserPassword(Username string) (string, error) {
	var Password string
	log.Println(Username)
	sql_stmt := "SELECT Pass FROM Users WHERE Username = $1"
	if err := db.QueryRow(sql_stmt, Username).Scan(&Password); err != nil {
		return "", err
	}
	return Password, nil
}
