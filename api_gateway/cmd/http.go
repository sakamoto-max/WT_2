package main

import (
	grpcclient "api_gateway/grpc_client"
	"api_gateway/handlers"
	"api_gateway/routes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	// "syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type httpServer struct {
	addr string
}
type httpServer2 struct {
	addr   string
	router *chi.Mux
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}
func NewHttpServer2(addr string) *httpServer2 {

	s := httpServer2{addr: addr}

	return &s
}

func (h *httpServer2) Run() {

	client := grpcclient.NewgrpcClient().ConnectToClients()
	log.Println("created the clients")
	handler := handlers.NewHandler(client.AuthClient, client.PlanClient, client.ExerClient, client.TrackClient)
	log.Println("created the handlers")
	router := routes.NewRouter(handler)
	log.Println("created the router")

	server := http.Server{
		Addr:    h.addr,
		Handler: router,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error opening http server at %v : %v\n", h.addr, err)

		}

	}()

	sig := <- sigChan
	log.Printf("shutdown signal received : %v", sig.String())


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("unable to shutdown the server : %v", err)
	}

	if err := client.Close(); err != nil {
		log.Println(err)
	}
	

	log.Printf("graceful shutdown complete")
}

func (h *httpServer) Run() {

	// server := http.Server{Addr: h.addr}
	// server.ListenAndServe()
	// server.Shutdown(ctx)

	client := grpcclient.NewgrpcClient()
	client = client.ConnectToClients()

	fmt.Println("all services are up and running")

	handler := handlers.NewHandler(client.AuthClient, client.PlanClient, client.ExerClient, client.TrackClient)

	router := routes.NewRouter(handler)

	fmt.Printf("gateway has started listening at %v\n", h.addr)

	if err := http.ListenAndServe(h.addr, router); err != nil {
		log.Fatalf("error opening http server at %v : %v\n", h.addr, err)
	}

}
