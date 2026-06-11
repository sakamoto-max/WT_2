package database

import (
	"auth_service/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPgConn(t *testing.T) {

	tests := []struct {
		name     string
		userName string
		pass     string
		host     string
		port     string
		dbname   string
		sslmode  string
		wantErr  bool
	}{
		{name: "parse failes", wantErr: true},
	}

	for _, test := range tests {
		_, err := NewPgConn(config.Config{
			Db: config.DbConfig{
				PgUser:         test.userName,
				PgPass:         test.pass,
				PgHost:         test.host,
				PgPort:         test.port,
				PgDatabaseName: test.dbname,
				PgSSLMode:      test.sslmode,
			},
		})

		if test.wantErr {
			assert.Error(t, err)
			return
		}

		assert.NoError(t, err)
	}
}
