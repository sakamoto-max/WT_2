package utils

import (
	"encoding/json"
	"testing"
	"github.com/go-openapi/testify/assert"
)

func Test_ConvertToJson(t *testing.T) {

	tests := []struct {
		name    string
		data    any
		wantErr bool
	}{
		{
			name: "success",
			data: map[string]any{"name": "this is the name"},
		},
		{
			name:    "failes",
			data:    "some random data",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			bytes, _ := json.Marshal(test.data)

			dataAfterConversion, err := ConvertToJson(bytes)

			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, dataAfterConversion, test.data)

		})
	}
}
