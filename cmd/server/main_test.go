package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
	"urlshort/internal/api"
	"urlshort/internal/urlservice"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpcGatewayStart вызывается асихронно через grpcStart и открывает gRPC Gateway сервер
func grpcGatewayStartTest() {
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

// gprcStart запускает сервер gRPC (в своей рутине) и сервер gRPC Gateway (в корутине)
func gprcStartTest(service *urlservice.Service) {
	grpcServ := grpc.NewServer()
	api.RegisterURLServer(grpcServ, service)
	lstn, err := net.Listen("tcp", ":"+os.Getenv("GPRC_PORT"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("gRPC: Started server on port " + os.Getenv("GPRC_PORT"))
		go grpcGatewayStartTest()
	}
	if err := grpcServ.Serve(lstn); err != nil {
		log.Fatal(err)
	}
}

func TestServer(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	service := &urlservice.Service{}
	service.HashMap = make(map[string]string)
	service.Method = memoryMethod
	go gprcStartTest(service)
	log.Println("Servers started, waiting 1 sec. timeout")
	time.Sleep(1000 * time.Millisecond)

	for i := 0; i < 10; i++ {
		want := strings.Repeat("a", i+1)

		/* POST */
		postURL := "http://localhost:" + os.Getenv("GPRC_GATEWAY_PORT") + "/post"
		data := []byte(`{"Value":"` + want + `"}`)
		r := bytes.NewReader(data)
		resPost, err := http.Post(postURL, "application/json", r)
		if err != nil {
			t.Fatal(err)
		}
		var received map[string]string
		json.NewDecoder(resPost.Body).Decode(&received)
		shortLink := received["Value"]
		received = nil
		resPost.Body.Close()
		log.Println("POST: \"" + want + "\", status code: " + strconv.Itoa(resPost.StatusCode) + ", received: " + shortLink)

		/* GET */
		getURL := "http://localhost:" + os.Getenv("GPRC_GATEWAY_PORT") + "/get/" + shortLink
		resGet, err := http.Get(getURL)
		if err != nil {
			t.Fatal(err)
		}
		json.NewDecoder(resGet.Body).Decode(&received)
		log.Println(received)
		fullLink := received["Value"]
		log.Println("GET: \"" + shortLink + "\", status code: " + strconv.Itoa(resGet.StatusCode) + ", received: " + fullLink)
		resGet.Body.Close()

		if fullLink != want {
			t.Fatal("Mismatch URL!")
		}
	}

	/* GET: Проверка несуществующей ссылки */
	getURL := "http://localhost:" + os.Getenv("GPRC_GATEWAY_PORT") + "/get/get-wrong-link"
	res, err := http.Get(getURL)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 500 {
		t.Fatal("Wrong implementation of GET for non-existent item")
	}

	/* POST: Добавление уже существующей ссылки не должно приводить к ошибке */
	postURL := "http://localhost:" + os.Getenv("GPRC_GATEWAY_PORT") + "/post"
	data := []byte(`{"Value":"` + "aaaaa" + `"}`)
	r := bytes.NewReader(data)
	resPost, err := http.Post(postURL, "application/json", r)
	if err != nil {
		t.Fatal(err)
	}
	if resPost.StatusCode != 200 {
		t.Fatal("Adding an existing link for POST should not be an error")
	}
}
