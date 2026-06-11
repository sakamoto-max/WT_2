package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HashThePassowrd(t *testing.T) {

	password := "jon snow king in the north"

	hashedPass, err := HashThePassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, hashedPass, password)
}

func Test_MatchPasswords(t *testing.T) {

	tests := []struct {
		name      string
		password  string
		wantErr   bool
		wrongPass string
	}{
		{name: "success", password: "jonsnowKingInTheNorth"},
		{name: "wrong password", password: "jonsnowKingInTheNorth", wrongPass: "jofferyistheking", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hashedPass, _ := HashThePassword(test.password)
			if test.wrongPass != "" {
				err := MatchPasswords(test.wrongPass, hashedPass)
				assert.Error(t, err)
				return
			}

			err := MatchPasswords(test.password, hashedPass)
			if test.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}

}
