package main

import (
	"log"
	"net"
	"os"
	"testing"
	"urlshort/internal/api"
	"urlshort/internal/urlservice"

	"google.golang.org/grpc"
)

func gprcStartTest(service *urlservice.Service, t *testing.T) error {
	grpcServ := grpc.NewServer()
	api.RegisterURLServer(grpcServ, service)
	_, err := net.Listen("tcp", ":"+os.Getenv("GPRC_PORT_TEST"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("gRPC: Started server on port " + os.Getenv("GPRC_PORT_TEST"))
		go grpcGatewayStart()
	}
	return nil
}

func GRPCCallTest(t *testing.T) {
	service := &urlservice.Service{}
	service.Method = memoryMethod
	if err := gprcStartTest(service, t); err != nil {
		log.Fatal(err)
		return
	}
}

func TestServer(t *testing.T) {
	GRPCCallTest(t)
}
