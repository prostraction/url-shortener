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
	con   *pgxpool.Pool
	Table string
}

func (db *DB) ConnectDB() (err error) {
	db.con, err = pgxpool.Connect(context.Background(), os.Getenv("DB_URL"))
	return err
}

func (db *DB) CreateTableDB(table string) (err error) {
	db.Table = table
	query := `CREATE TABLE IF NOT EXISTS ` + table + ` (id SERIAL, full_link text, short_link text);`
	err = db.con.QueryRow(context.Background(), query).Scan()
	if err != nil && err.Error() != "no rows in result set" {
		return
	}
	query = `CREATE INDEX IF NOT EXISTS ind ON ` + table + ` USING HASH (full_link);`
	err = db.con.QueryRow(context.Background(), query).Scan()
	if err != nil && err.Error() != "no rows in result set" {
		return
	} else {
		return nil
	}
}

func (db *DB) FromHash(hash string) (string, error) {
	var url string
	query := `SELECT full_link FROM ` + db.Table + ` WHERE short_link = $1`
	err := db.con.QueryRow(context.Background(), query, hash).Scan(&url)
	return url, err
}

func (db *DB) ToHash(url string) (string, error) {
	query := `INSERT INTO ` + db.Table + ` (full_link, short_link) VALUES ($1, $2) RETURNING short_link;`
	hash := hashfunc.GetBaseEnc(url)[:10]
	for i := 10; i < 32; i++ {
		if value, exists := db.FromHash(hash[i-10 : i]); exists != nil && exists.Error() == "no rows in result set" {
			err := db.con.QueryRow(context.Background(), query, url, hash[i-10:i]).Scan(&hash)
			return hash, err
		} else if value == url {
			/* This URL is already on hash rable */
			return hash[i-10 : i], fmt.Errorf("url is already on hash table (%s)", hash[i-10:i])
		}
	}
	return "", errors.New("collision not resolved")
}
