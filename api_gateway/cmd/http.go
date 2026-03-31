package main

import (
	grpcclient "api_gateway/grpc_client"
	"api_gateway/handlers"
	"api_gateway/routes"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	// "syscall"
	"time"

	"github.com/go-chi/chi/v5"
)


type httpServer struct {
	addr   string
	router *chi.Mux
}

func NewHttpServer(addr string) *httpServer {

	s := httpServer{addr: addr}

	return &s
}

func (h *httpServer) Run() {

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

	go func() {

		log.Printf("server has started at %v", h.addr)
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
