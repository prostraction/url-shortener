package main

import (
	"log"
	"net"
	"urlshort/internal/urlshort"
	"urlshort/pkg/api"

	"google.golang.org/grpc"
)

const (
	memoryMethod = iota
	dbMethod
)

func main() {
	grpc := grpc.NewServer()
	app := &urlshort.Application{}
	app.HashMap = make(map[string]string)

	app.Method = dbMethod
	if app.Method == dbMethod {
		err := app.InitDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	api.RegisterURLServer(grpc, app)

	lstn, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("gRPC server started on port 50001")
		log.Println("gRPC server using memory to store links")
	}
	if err := grpc.Serve(lstn); err != nil {
		log.Fatal(err)
	}
}
