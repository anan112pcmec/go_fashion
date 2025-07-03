package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"github.com/anan112pcmec/go_fashion/app"
)

func generateRandomKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// updateSecretKey memperbarui SECRET_KEY di file .env
func updateSecretKey(filename string) error {
	secretKey, err := generateRandomKey()
	if err != nil {
		return err
	}

	envMap, err := godotenv.Read(filename)
	if err != nil {
		return err
	}

	envMap["SECRET_KEY"] = secretKey // Update SECRET_KEY

	return godotenv.Write(envMap, filename) // Simpan perubahan ke file .env
}

func main() {

	err3 := updateSecretKey(".env")
	if err3 != nil {
		log.Fatal("Gagal memperbarui SECRET_KEY:", err3)
	}
	fmt.Println("SECRET_KEY berhasil diperbarui di .env")

	app.Jalan()
}
