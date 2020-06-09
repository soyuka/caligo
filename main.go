package main

import (
	"log"
	"net/http"

	c "github.com/soyuka/caligo/config"
	"github.com/soyuka/caligo/handlers"

	bolt "go.etcd.io/bbolt"
)

func main() {
	config := c.GetConfig()

	db, err := bolt.Open(config.DBPath, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	env := &handlers.Env{
		DB:     db,
		Config: config,
	}

	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(config.DBBucketName))
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/favicon.ico", handlers.Handler{env, handlers.Favicon})
	http.Handle("/", handlers.Handler{env, handlers.GetIndex})

	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
