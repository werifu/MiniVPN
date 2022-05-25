package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
)

// Hash returns passwdHash (not so safe but easy)
func Hash(passwd string, salt string) string {
	hash := sha256.Sum256([]byte(passwd + salt))
	return fmt.Sprintf("%x", hash[:])
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Cmd should be:\nadd_user your_username your_password")
	}
	username := os.Args[1]
	passwd := os.Args[2]
	fmt.Println("Generate a new user: " + username)
	fmt.Println("Password after hash: " + Hash(passwd, username))
}
