package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type DBStorage struct {
	ConnString string
}

func NewDBStorage(connString string) *DBStorage {
	storage := &DBStorage{ConnString: connString}
	return storage
}

func (s *DBStorage) Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), s.ConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	return conn, nil
}
