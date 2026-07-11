package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	size := 32

	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("failed to read random key: %v", err)
	}

	secret := base64.RawURLEncoding.EncodeToString(key)
	fmt.Printf("Add this to .env:\nJWT_SECRET=%s\n", secret)
}
