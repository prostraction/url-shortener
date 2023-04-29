package urlshort

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/jackc/pgx/v4"
)

func TestT(t *testing.T) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, "postgres://postgress:simpleDBpassword@localhost:5432/url")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(db)
	}
}
