package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"urlshort/internal/api"

	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.Parse()
	// TO DO
	if flag.NArg() < 2 {
		log.Fatal("not enough arguments")
	}

	operation := flag.Arg(0)
	url := flag.Arg(1)

	conn, err := grpc.Dial(":50001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewURLClient(conn)
	switch strings.ToLower(operation) {
	case "toshort":
		if res, err := client.ToShortLink(context.Background(), &api.FullURL{Value: url}); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(res.Value)
		}

	case "tofull":
		if res, err := client.ToFullLink(context.Background(), &api.ShortURL{Value: url}); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(res.Value)
		}
	}
}
