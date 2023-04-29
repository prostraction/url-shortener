package urlshort

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

func test() {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, "postgres://postgress:simpleDBpassword@localhost:5432/")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(db)
	}
}
