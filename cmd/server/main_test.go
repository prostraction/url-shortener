package main

import (
	"log"
	"testing"
	"urlshort/internal/urlservice"
)

func GRPCTest(t *testing.T) {
	service := &urlservice.Service{}
	service.Method = dbMethod
	table := "url"
	err := service.InitDbWithTable(table)
	if err != nil {
		log.Fatal(err)
		return
	}
	gprcStart(service)
}

func GRPCGatewayTest(t *testing.T) {
	grpcGatewayStart()
}

func TestServer(t *testing.T) {
	go GRPCTest(t)
	go GRPCGatewayTest(t)

}
