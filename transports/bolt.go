package transports

import (
	"fmt"
	"net/url"

	bolt "go.etcd.io/bbolt"
)

// BoltTransport implements the TransportInterface using the Bolt database.
type BoltTransport struct {
	db     *bolt.DB
	bucketName string
}

const defaultBoltBucketName = "caligo"

// NewBoltTransport create a new BoltTransport.
func NewBoltTransport(u *url.URL) (*BoltTransport, error) {
	var err error
	q := u.Query()
	bucketName := defaultBoltBucketName
	if q.Get("bucket_name") != "" {
		bucketName = q.Get("bucket_name")
	}

	path := u.Path // absolute path (bolt:///path.db)
	if path == "" {
		path = u.Host // relative path (bolt://path.db)
	}
	if path == "" {
		return nil, fmt.Errorf(`%q: missing path: %w`, u, ErrInvalidTransportDSN)
	}

	db, err := bolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf(`%q: %s: %w`, u, err, ErrInvalidTransportDSN)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(bucketName))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf(`%q: bucket could not be created: %w`, u, err)
	}


	return &BoltTransport{
		db:               db,
		bucketName:       bucketName,
	}, nil
}

func (b *BoltTransport) Put(id string, url string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		err := b.Put([]byte(id), []byte(url))
		return err
	})
}

func (b *BoltTransport) Get(id string) (string, error) {
	var url string
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucketName))
		url = string(b.Get([]byte(id)))
		return nil
	})

	return url, err
}
