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
	lstn, err := net.Listen("tcp", ":"+os.Getenv("GPRC_PORT"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("gRPC: Started server on port " + os.Getenv("GPRC_PORT"))
		go grpcGatewayStart()
	}
	return grpcServ.Serve(lstn)
}

func GRPCCallTest(t *testing.T) {
	service := &urlservice.Service{}
	service.Method = dbMethod
	table := "url"
	err := service.InitDbWithTable(table)
	if err != nil {
		log.Fatal(err)
		return
	}
	if gprcStartTest(service, t) != nil {
		log.Fatal(err)
		return
	}
}

func TestServer(t *testing.T) {
	GRPCCallTest(t)
}
