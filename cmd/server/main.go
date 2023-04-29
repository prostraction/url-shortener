package main

import (
	"log"
	"net"
	"urlshort/internal/urlshort"
	"urlshort/pkg/api"

	"google.golang.org/grpc"
)

func main() {
	grpc := grpc.NewServer()
	app := &urlshort.Application{}
	app.HashMap = make(map[string]string)
	api.RegisterURLServer(grpc, app)

	lstn, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatal(err)
	}
	if err := grpc.Serve(lstn); err != nil {
		log.Fatal(err)
	}
}
