package utils

import (
	// "encoding/json"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	myErrs "wt/pkg/my_errors"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
)

func HashThePassword(password string) (string, error) {
	passInBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", myErrs.InternalServerErrMaker(fmt.Errorf("error encrypting the password : %w", err))
	}
	return string(passInBytes), nil
}

func MatchPasswords(password string, passFromDb string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passFromDb), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
			return myErrs.BadReqErrMaker(ErrIncorrectPassword)
		}
		return myErrs.InternalServerErrMaker(fmt.Errorf("some error occured while authenticating the password : %w", err))
	}
	return nil
}

