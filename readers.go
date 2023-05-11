package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readBotToken(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}

	return "", fmt.Errorf("no token found in %s", filename)
}

func encryptID(id int64) int64 {
    // XOR each byte of the ID with a key
    key := int64(42) // Choose a secret key
    encrypted := id ^ key

    return encrypted
}

func decryptID(encrypted int64) int64 {
    // XOR each byte of the encrypted ID with the same key used for encryption
    key := int64(42) // Secret key used during encryption
    decrypted := encrypted ^ key

    return decrypted
}
