package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)


var(
	CurrentUserEmail string = "test8@gmail.com"
	Password string = "Password123"
	CurrentUserName string = "test8"
)

func Test_SignUp(t *testing.T) {
	client := &http.Client{}

	payload := map[string]string{
		"email":    CurrentUserEmail,
		"password": Password,
		"name":     CurrentUserName,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "http://localhost:5000/wt/user/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil{
		t.Fatalf("request failed : %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated{
		t.Fatalf("expected 201, got %v", resp.StatusCode)
	}

	var res map[string]any

	json.NewDecoder(resp.Body).Decode(&res)
	if res["name"] == nil{
		t.Fatalf("expected name but got nil")
	}
	
	if res["email"] == nil {
		t.Fatalf("expected email but got nil")
		
	}

	if res["role"] == nil {
		t.Fatalf("expected role but got nil")
	}
}
func Test_DuplicateSignUp(t *testing.T) {
	client := &http.Client{}

	payload := map[string]string{
		"email":    CurrentUserEmail,
		"password": Password,
		"name":     CurrentUserName,
	}
	
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "http://localhost:5000/wt/user/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil{
		t.Fatalf("request failed : %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest{
		t.Fatalf("expected 400, got %v", resp.StatusCode)
	}
}

