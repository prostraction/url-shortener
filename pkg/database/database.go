package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"urlshort/pkg/hashfunc"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	con *pgxpool.Pool
}

func (db *DB) ConnectDB() (err error) {
	db.con, err = pgxpool.Connect(context.Background(), os.Getenv("DB_URL"))
	return err
}

func (db *DB) CreateTableDB() (err error) {
	err = db.con.QueryRow(context.Background(),
		`CREATE TABLE IF NOT EXISTS url (id SERIAL, full_link text, short_link text);`).Scan()
	if err != nil && err.Error() != "no rows in result set" {
		return
	}
	err = db.con.QueryRow(context.Background(),
		`CREATE INDEX IF NOT EXISTS ind ON url USING HASH (full_link);`).Scan()
	if err != nil && err.Error() != "no rows in result set" {
		return
	} else {
		return nil
	}
}

func (db *DB) FromHash(hash string) (string, error) {
	var url string
	err := db.con.QueryRow(context.Background(),
		`SELECT full_link FROM url WHERE short_link = $1`, hash).Scan(&url)
	return url, err
}

func (db *DB) ToHash(url string) (string, error) {
	hash := hashfunc.GetBaseEnc(url)[:10]
	for i := 10; i < 32; i++ {
		if value, exists := db.FromHash(hash[i-10 : i]); exists != nil && exists.Error() == "no rows in result set" {
			err := db.con.QueryRow(context.Background(),
				`INSERT INTO url (full_link, short_link) VALUES ($1, $2) RETURNING short_link;`,
				url, hash[i-10:i]).Scan(&hash)
			return hash, err
		} else if value == url {
			/* This URL is already on hash rable */
			return hash[i-10 : i], fmt.Errorf("url is already on hash table (%s)", hash[i-10:i])
		}
	}
	return "", errors.New("collision not resolved")
}
