package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
)

type DbClient struct {
	db *leveldb.DB
}

func InitDblient() (*DbClient, error) {
	//Go user home dir
	homDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbDir := filepath.Join(homDir, ".tcp-client-server/server/db")
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %v", err)
	}
	d, err := leveldb.OpenFile(dbDir, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %v", err)
	}

	log.Printf("db is created at: %v\n", dbDir)

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
