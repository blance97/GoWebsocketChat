package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
		SessionID TEXT,
		CurrentRoom TEXT,
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
func updateCurrentRoom(Username string, Roomname string) {
	sql_stmt := `UPDATE Users SET CurrentRoom = $1 WHERE Username = $2`
	if _, err := db.Exec(sql_stmt, Roomname, Username); err != nil {
		log.Println("Error in Updadint current room: ", err)
		return
	}
	log.Println("Updated CurrentRoom")
	return
}
func listUsersinRoom(Roomname string) []string {
	sql_stmt := `SELECT Username FROM	Users WHERE CurrentRoom=$1`
	dbMu.Lock()
	rows, err := db.Query(sql_stmt, Roomname)
	dbMu.Unlock()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var a string
		err2 := rows.Scan(&a)
		if err2 != nil {
			log.Println("Error scanning username")
		}
		result = append(result, a)
	}
	log.Println(result)
	return result
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
	if _, err := db.Exec(sql_stmt, sid, Username); err != nil {
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

func createRoomTable() {
	sql_table := `
	CREATE TABLE IF NOT EXISTS Rooms(
		Owner TEXT ,
		Roomname TEXT PRIMARY KEY,
		Private TEXT,
		Password TEXT,
		DateCreated
	);
	`
	_, err := db.Exec(sql_table)
	if err != nil {
		log.Print(err)
	}
}
func StoreRoomInfo(Owner string, Roomname string, Private string, Password string) {
	sql_stmt := `
	INSERT OR REPLACE INTO Rooms(
		Owner,
		Roomname,
		Private,
		Password,
    DateCreated
	)values(?, ?, ?, ?, ?)
	`
	stmt, err := db.Prepare(sql_stmt)
	if err != nil {
		log.Print(err)
	}
	c := Rooms{
		Owner:       Owner,
		Roomname:    Roomname,
		Private:     Private,
		Password:    Password,
		DateCreated: time.Now().Unix(),
	}
	if _, err := stmt.Exec(c.Owner, c.Roomname, c.Private, c.Password, c.DateCreated); err != nil {
		log.Println(err)
	}
	log.Println("Store New Room Info")
}
func RoomExist(Roomname string) bool {
	sql_stmt := `SELECT * FROM Rooms WHERE Rooms WHERE Roomname = $1`
	err := db.QueryRow(sql_stmt, Roomname).Scan()
	if err == sql.ErrNoRows {
		return false
	}
	return true
}
func RemoveRoom(Roomname string) {
	sql_stmt := `DELETE FROM Rooms WHERE Roomname=$1`
	if _, err := db.Exec(sql_stmt, Roomname); err != nil {
		log.Println("Error Deleting Room From Database", err)
	}
}
func PrivateRoomChecker(Roomname string) bool {
	var priv string
	sql_stmt := `SELECT Private FROM Rooms WHERE Roomname=$1`
	if err := db.QueryRow(sql_stmt, Roomname).Scan(&priv); err != nil {
		log.Println("Error query for private room checker", err)
	}
	if priv == "false" {
		return false
	}
	return true
}
func GetPrivateRoomPass(Roomname string) string {
	var pass string
	sql_stmt := `SELECT Password FROM Rooms WHERE Roomname=$1`
	if err := db.QueryRow(sql_stmt, Roomname).Scan(&pass); err != nil {
		log.Println("Error query for private room checker", err)
	}
	return pass
}
