package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"urlshort/internal/api"
	"urlshort/internal/urlservice"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	memoryMethod = iota
	dbMethod
)

func grpcGatewayStart() {
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:"+os.Getenv("GPRC_PORT"),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("gPRC-Gateway: Dial failed:", err)
	}

	mux := runtime.NewServeMux()
	err = api.RegisterURLHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalln("gPRC-Gateway: Register gateway failed:", err)
	}
	gwServer := &http.Server{
		Addr:    ":" + os.Getenv("GPRC_GATEWAY_PORT"),
		Handler: mux,
	}

	log.Println("gRPC-Gateway: Started server on port " + os.Getenv("GPRC_GATEWAY_PORT"))
	log.Fatal(gwServer.ListenAndServe())
}

func gprcStart(service *urlservice.Service) {
	grpcServ := grpc.NewServer()
	api.RegisterURLServer(grpcServ, service)
	lstn, err := net.Listen("tcp", ":"+os.Getenv("GPRC_PORT"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("gRPC: Started server on port " + os.Getenv("GPRC_PORT"))
		go grpcGatewayStart()
	}
	if err := grpcServ.Serve(lstn); err != nil {
		log.Fatal(err)
	} else {
		log.Println("here")
	}
}

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var method string
	flag.StringVar(&method, "method", "db", "method used to store hash")
	flag.Parse()
	log.Println("gRPC server using " + method + " to store links")

	service := &urlservice.Service{}
	switch method {
	case "memory":
		service.HashMap = make(map[string]string)
		service.Method = memoryMethod
	case "db":
		service.Method = dbMethod
		table := "url"
		err := service.InitDbWithTable(table)
		if err != nil {
			log.Fatal(err)
		}
	}
	gprcStart(service)
}
