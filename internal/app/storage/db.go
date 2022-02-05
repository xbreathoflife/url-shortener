package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/xbreathoflife/url-shortener/internal/app/entities"
	"log"
)

const (
	createTableQuery = "CREATE TABLE IF NOT EXISTS url(" +
		"    id           SERIAL PRIMARY KEY," +
		"    original_url TEXT NOT NULL UNIQUE," +
		"    short_url    TEXT NOT NULL," +
		"    uuid         TEXT NOT NULL" +
		")"

	countURLQuery = "SELECT COUNT(*) as count FROM url"

	insertURLQuery = "INSERT INTO url(id, original_url, short_url, uuid) VALUES ($1, $2, $3, $4)"

	getURLByIDQuery = "SELECT original_url FROM url WHERE id = $1"

	getURLsByUserQuery = "SELECT original_url, short_url, uuid FROM url WHERE uuid = $1"

	getExistingURL = "SELECT short_url FROM url WHERE original_url = $1"
)
type DBStorage struct {
	ConnString string
	BaseURL    string
}

func NewDBStorage(connString string, baseURL string) *DBStorage {
	storage := &DBStorage{ConnString: connString, BaseURL: baseURL}
	return storage
}

func (s *DBStorage) Init(ctx context.Context) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, createTableQuery)

	return err
}

func (s *DBStorage) CheckConnect(_ context.Context) error {
	conn, err := sql.Open("pgx", s.ConnString)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	defer conn.Close()

	return nil
}

func (s *DBStorage) connect(_ context.Context) (*sql.DB, error) {
	if s.ConnString == "" {
		log.Fatal("Connection string is empty\n")
	}
	conn, err := sql.Open("pgx", s.ConnString)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	return conn, nil
}

func (s *DBStorage) InsertNewURL(ctx context.Context, id int, baseURL string, shortenedURL string, uuid string) error {
	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, insertURLQuery, id, baseURL, shortenedURL, uuid)
	return err
}

func (s *DBStorage) GetURLByID(ctx context.Context, id int) (string, error) {
	conn, err := s.connect(ctx)
	if err != nil {
		return "", err
	}

	defer conn.Close()
	var u string
	row := conn.QueryRowContext(ctx, getURLByIDQuery, id)
	err = row.Scan(&u)
	if err != nil {
		fmt.Printf("Unable to get row count: %v\n", err)
		return "", err
	}
	return u, nil
}

func (s *DBStorage) GetUserURLs(ctx context.Context, uuid string) ([]entities.URL, error) {
	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	rows, err := conn.QueryContext(ctx, getURLsByUserQuery, uuid)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []entities.URL
	for rows.Next() {
		var u entities.URL
		if err := rows.Scan(&u.BaseURL, &u.ShortenedURL, &u.UserID); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	return urls, nil
}

func (s *DBStorage) GetNextID(ctx context.Context) (int, error) {
	conn, err := s.connect(ctx)
	if err != nil {
		return 0, err
	}

	defer conn.Close()
	var id int
	row := conn.QueryRowContext(ctx, countURLQuery)
	err = row.Scan(&id)
	if err != nil {
		fmt.Printf("Unable to get row count: %v\n", err)
		return 0, err
	}

	return id, nil
}

func (s *DBStorage) GetURLIfExist(ctx context.Context, url string) (string, error) {
	conn, err := s.connect(ctx)
	if err != nil {
		return "", err
	}

	defer conn.Close()
	var str sql.NullString
	err = conn.QueryRowContext(ctx, getExistingURL, url).Scan(&str)

	if err != nil && err.Error() != "no rows in result set" {
		return "", err
	}

	if str.Valid {
		return str.String, nil
	}
	return "", nil
}

func (s *DBStorage) GetBaseURL() string {
	return s.BaseURL
}

func (s *DBStorage) InsertBatch(ctx context.Context, records []entities.Record) error {
	db, err := s.connect(ctx)
	defer db.Close()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO url(id, original_url, short_url, uuid) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	for _, r := range records {
		if _, err = stmt.Exec(r.ID, r.BaseURL, r.ShortenedURL, r.UserID); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}