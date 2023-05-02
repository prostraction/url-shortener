package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"urlshort/internal/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// main парсит аргументы программы и вызывает методы gRPC саервера
func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Not enough arguments. Syntax: (toShort URL | toFull URL)")
	}

	operation := flag.Arg(0)
	url := flag.Arg(1)

	conn, err := grpc.Dial(":"+os.Getenv("GPRC_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
