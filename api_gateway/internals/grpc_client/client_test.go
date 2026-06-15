package grpcclient

import (
	"testing"

	"github.com/go-openapi/testify/assert"
	"github.com/sakamoto-max/wt_2_pkg/logger"
)

// NewgrpcClient()
// ConnectToClients()
// close()

func Test_NewGrpcClient(t *testing.T) {
	client := NewgrpcClient()

	assert.NotNil(t, client)
}

func Test_ConnectToClients(t *testing.T) {

	client := NewgrpcClient()

	logger := logger.NewLogger()

	client = client.ConnectToClients(logger)

	assert.NotNil(t, client.authConn)
	assert.NotNil(t, client.planConn)
	assert.NotNil(t, client.exerConn)
	assert.NotNil(t, client.trackConn)
	assert.NotNil(t, client.AuthClient)
	assert.NotNil(t, client.PlanClient)
	assert.NotNil(t, client.ExerClient)
	assert.NotNil(t, client.TrackClient)

	client.Close()
}

func Test_Close(t *testing.T) {

	client := NewgrpcClient()

	logger := logger.NewLogger()

	client = client.ConnectToClients(logger)

	err := client.Close()

	assert.NoError(t, err)
}
