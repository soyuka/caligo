package storage

import (
	"context"
	"time"

	"go.etcd.io/etcd/v3/clientv3"
)

type Config struct {
	Etcd              clientv3.Config
	ShortenerHostname string
}

func Read(config clientv3.Config, key string) (string, error) {
	cli, err := clientv3.New(config)

	if err != nil {
		return "", err
	}

	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()

	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", nil
	}

	return string(resp.Kvs[0].Value), nil
}

func Write(config clientv3.Config, key string, value string) error {
	cli, err := clientv3.New(config)

	if err != nil {
		return err
	}

	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = cli.Put(ctx, key, value)
	cancel()

	if err != nil {
		return err
	}

	return nil
}
