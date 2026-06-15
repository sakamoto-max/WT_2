package routes

import (
	"testing"

	"github.com/go-openapi/testify/assert"
	grpcclient "github.com/sakamoto-max/wt_2/api_gateway/internals/grpc_client"
	"github.com/sakamoto-max/wt_2/api_gateway/internals/handlers"
	"github.com/sakamoto-max/wt_2_pkg/logger"
)

func Test_NewRouter(t *testing.T) {

	logger := logger.NewLogger()
	client := grpcclient.NewgrpcClient().ConnectToClients(logger)
	handler := handlers.NewHandler(client.AuthClient, client.PlanClient, client.ExerClient, client.TrackClient)


	router := NewRouter(handler)

	assert.NotNil(t, router)

	client.Close()
}
