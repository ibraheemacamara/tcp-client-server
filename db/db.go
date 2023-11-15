package db

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type DbClient struct {
	db *leveldb.DB
}

func InitDblient() (*DbClient, error) {
	d, err := leveldb.OpenFile("path/to/db", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %v", err)
	}

	return &DbClient{
		db: d,
	}, nil
}

func (c *DbClient) Get(key []byte) ([]byte, error) {
	return c.db.Get(key, nil)
}

func (c *DbClient) Put(key []byte, data []byte) error {
	return c.db.Put(key, data, nil)
}
