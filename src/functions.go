package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"os"
		"log"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func initRooms() {
	mu.Lock()
	defer mu.Unlock()
	files, _ := ioutil.ReadDir("log/")
	if len(files) == 0 {
		_, err := os.OpenFile("log/"+"room1", os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
		}
		StoreRoomInfo("Admin", "room1", "false", "")
		server := NewServer("/entry/" + "room1") // start server
		go server.Listen("room1")
	} else {
		for _, file := range files {
			server := NewServer("/entry/" + file.Name())
			go server.Listen(file.Name())
		}
	}
}
