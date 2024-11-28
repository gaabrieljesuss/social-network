package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Porta     = 0
	SecretKey []byte
)

func Carregar() {
	var erro error

	if erro = godotenv.Load(".env"); erro != nil {
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("SERVER_PORT"))
	if erro != nil {
		Porta = 8000
	}

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
