package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"urlshort/pkg/hashfunc"
)

func TestDatabase(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	var test_db DB
	var str string
	var err error
	log.Println("Connecting to DB")
	err = test_db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Creating table")
	err = test_db.CreateTableDB("test_url")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Filling table with values and validating")
	for i := 0; i < 10; i++ {
		want := strings.Repeat("a", i+1)
		str, err = test_db.ToHash(want)
		if err != nil {
			log.Fatal(err)
		}
		str, err = test_db.FromHash(str)
		if err != nil {
			log.Fatal(err)
		}
		if str != want {
			log.Fatal(errors.New("str != want"))
		}
	}
	log.Println("Getting wrong number")
	str, err = test_db.FromHash("bbbbbbbbbb")
	if str != "" {
		log.Fatal(fmt.Errorf("waiting for empty response, got '%s' with error '%s", str, err.Error()))
	}
	log.Println("Dropping table")
	err = test_db.con.QueryRow(context.Background(), "DROP TABLE test_url").Scan()
	if err != nil && err.Error() != "no rows in result set" {
		log.Fatal(err)
	}
	log.Println("Validating drop")

	want := strings.Repeat("a", 5)
	str = hashfunc.GetBaseEnc(want)[:10]
	_, err = test_db.FromHash(str)
	if err == nil {
		log.Fatal("table is still loaded")
	} else if err != nil && strings.Contains(err.Error(), "does not exist") {
		log.Println("Done")
	} else {
		log.Fatal(err)
	}
}
