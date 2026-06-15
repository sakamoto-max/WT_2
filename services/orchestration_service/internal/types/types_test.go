package types

import (
	"encoding/json"
	"testing"

	"github.com/go-openapi/testify/assert"
)

func Test_ConvertToBytes(t *testing.T) {

	data := Data{
		DbId:          "123",
		TargetService: "target",
	}

	bytes, err := data.ConvertToBytes()
	assert.NoError(t, err)

	var convertedData Data

	json.Unmarshal(*bytes, &convertedData)

	assert.Equal(t, convertedData, data)

}
