package main

import (
	grpcclient "api_gateway/grpc_client"
	"api_gateway/handlers"
	"api_gateway/routes"
	"fmt"
	"log"
	"net/http"
)

type httpServer struct {
	addr string
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func (h *httpServer) Run() {

	client := grpcclient.NewgrpcClient()
	client = client.ConnectToClients()
	client.PingAll()

	fmt.Println("all services are up and running")

	handler := handlers.NewHandler(client.AuthClient, client.PlanClient, client.ExerClient, client.TrackClient)

	router := routes.NewRouter(handler)

	fmt.Printf("gateway has started listening at %v", h.addr)

	if err := http.ListenAndServe(h.addr, router); err != nil {
		log.Fatalf("error opening http server at %v : %v", h.addr, err)
	}
}
