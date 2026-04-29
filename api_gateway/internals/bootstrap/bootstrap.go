package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/sakamoto-max/wt_2_pkg/logger"
	grpcclient "github.com/sakamoto-max/wt_2/api_gateway/internals/grpc_client"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/handlers"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/routes"
	"go.uber.org/zap"
)

type app struct {
	addr   string
	router *chi.Mux
}

func NewApp(addr string) *app {

	s := app{addr: addr}

	return &s
}

func (h *app) Run() {

	logger := logger.NewLogger()
	defer logger.Log.Sync()

	logger.Log.Info("starting the server")

	client := grpcclient.NewgrpcClient().ConnectToClients(logger)

	logger.Log.Info("connected to the grpc clients")

	handler := handlers.NewHandler(
		client.AuthClient,
		client.PlanClient,
		client.ExerClient,
		client.TrackClient,
	)

	router := routes.NewRouter(handler)

	logger.Log.Info("created the handlers")

	server := http.Server{
		Addr:    h.addr,
		Handler: router,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		logger.Log.Infow("server has started", zap.String("addr", h.addr))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Warnw("failed to run the http server", zap.Error(err))
			sigChan <- os.Interrupt
		}

	}()

	sig := <-sigChan
	logger.Log.Infow("shutdown signal received", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Errorf("unable to shutdown the server :", err)
	}

	if err := client.Close(); err != nil {
		logger.Log.Errorf("unable to close the clients : %v", err)
	}

	logger.Log.Info("graceful shutdown complete")
}
