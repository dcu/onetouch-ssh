package main

import (
	"os"
)

func main() {
	authyID := os.Args[1]

	verification := NewVerification(authyID)
	verification.perform()
}
