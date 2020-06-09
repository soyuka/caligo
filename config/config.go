package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBPath            string
	DBBucketName      string
	ShortenerHostname string
	IdAlphabet        string
	IdLength          int
	Port              string
}

func GetConfig() Config {
	dbPath := os.Getenv("CALIGO_DB_PATH")

	if dbPath == "" {
		dbPath = "data.bolt"
	}

	shortenerHostname := os.Getenv("CALIGO_HOSTNAME")

	if shortenerHostname == "" {
		shortenerHostname = "http://localhost:5376"
	}

	port := os.Getenv("CALIGO_PORT")
	if port == "" {
		port = "5376"
	}

	idLength, err := strconv.ParseInt(os.Getenv("CALIGO_ID_LENGTH"), 10, 32)

	if idLength == 0 || err != nil {
		idLength = 12
	}

	idAlphabet := os.Getenv("CALIGO_ID_ALPHABET")

	if idAlphabet == "" {
		idAlphabet = "0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"
	}

	// todo: log config
	log.Println("DB Path", dbPath)
	log.Println("Hostname", shortenerHostname)

	return Config{
		ShortenerHostname: shortenerHostname,
		IdLength:          int(idLength),
		IdAlphabet:        idAlphabet,
		Port:              port,
		DBPath:            dbPath,
		DBBucketName:      "caligo",
	}
}
