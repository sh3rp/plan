package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
)

func main() {
	fmt.Printf("Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		fmt.Printf("Error getting password: %v\n")
		os.Exit(1)
	}
	hasher := sha1.New()
	hasher.Write(pass)
	hash := hasher.Sum(nil)
	fmt.Printf("Password = %s\n", hex.EncodeToString(hash))
}
