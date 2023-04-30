package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestT(t *testing.T) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(db)
	}
}
