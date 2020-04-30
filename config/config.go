package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/etcd/v3/clientv3"
)

type Config struct {
	Etcd              clientv3.Config
	ShortenerHostname string
	IdAlphabet        string
	IdLength          int
	Port              string
}

func GetConfig() Config {

	etcdUrl := os.Getenv("ETCD_URL")
	endpoints := []string{"localhost:2379"}

	if etcdUrl != "" {
		endpoints = strings.Split(etcdUrl, ",")
	}

	dialTimeout, _ := time.ParseDuration(os.Getenv("ETCD_DIAL_TIMEOUT"))

	if dialTimeout == 0 {
		dialTimeout = 5 * time.Second
	}

	shortenerHostname := os.Getenv("CALIGO_HOSTNAME")

	if shortenerHostname == "" {
		shortenerHostname = "http://localhost:8080"
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
	log.Println("ETCD endpoints", endpoints)
	log.Println("Hostname", shortenerHostname)

	return Config{
		ShortenerHostname: shortenerHostname,
		IdLength:          int(idLength),
		IdAlphabet:        idAlphabet,
		Port:              port,
		Etcd: clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: dialTimeout,
		},
	}
}
