package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"urlshort/internal/api"
	"urlshort/internal/urlshort"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	memoryMethod = iota
	dbMethod
)

/*
var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)*/

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

func gprcStart(service *urlshort.Service) {
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
	}
}

func main() {
	log.Println("gRPC server using METHOD to store links")
	service := &urlshort.Service{}
	service.HashMap = make(map[string]string)
	service.Method = dbMethod
	if service.Method == dbMethod {
		err := service.InitDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	gprcStart(service)
}
