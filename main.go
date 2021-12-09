package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func main() {
	proofOfWork()
	// blockchain.Blcokchain()
	// cli.Start()
}

func proofOfWork() {
	difficulty := 2
	nonce := 1
	target := strings.Repeat("0", difficulty)
	for {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte("Test"+fmt.Sprint(nonce))))
		fmt.Printf("Target : %s\n Hash : %s\n Nonce : %d\n--------\n", target, hash, nonce)
		if strings.HasPrefix(hash, target) {
			break
		} else {
			nonce++
		}
	}
}
