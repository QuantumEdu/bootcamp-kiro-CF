package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	pins := map[string]string{"1234": "Admin", "1235": "Maria Cajera"}
	for pin, name := range pins {
		hash, _ := bcrypt.GenerateFromPassword([]byte(pin), 10)
		fmt.Printf("%s (PIN %s): %s\n", name, pin, string(hash))
	}
}
