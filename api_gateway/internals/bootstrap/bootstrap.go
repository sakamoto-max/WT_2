package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/config"
	grpcclient "github.com/sakamoto-max/wt_2/api_gateway/internals/grpc_client"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/handlers"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/routes"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	"go.uber.org/zap"
)

type app struct {
	addr    string
	logger  *logger.MyLogger
	config  config.Config
	router  *chi.Mux
	clients *grpcclient.GrpcClient
}

func NewApp(config config.Config) *app {

	client := grpcclient.ConnectToClients(config)

	handler := handlers.NewHandler(
		client.AuthClient,
		client.PlanClient,
		client.ExerClient,
		client.TrackClient,
	)

	router := routes.NewRouter(handler)

	return &app{
		addr: config.HttpServer.Addr,
		config: config,
		router: router,
		clients: client,
	}
}

func (h *app) Run() {

	// logger.Log.Info("starting the server")

	// logger.Log.Info("connected to the grpc clients")

	// logger.Log.Info("created the handlers")

	server := http.Server{
		Addr:    h.addr,
		Handler: h.router,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		h.logger.Log.Infow("server has started", zap.String("addr", h.addr))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.logger.Log.Warnw("failed to run the http server", zap.Error(err))
			sigChan <- os.Interrupt
		}

	}()

	sig := <-sigChan
	h.logger.Log.Infow("shutdown signal received", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		h.logger.Log.Errorf("unable to shutdown the server :", err)
	}

	if err := h.clients.Close(); err != nil {
		h.logger.Log.Errorf("unable to close the clients : %v", err)
	}

	h.logger.Log.Info("graceful shutdown complete")
}
