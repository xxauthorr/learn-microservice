package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Ping(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "This is the authentication-server",
	}
	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err.Error())
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("error :"+err.Error()), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("error in matching password"), http.StatusUnauthorized)
		return
	}
	err = app.logRequest("authentication", fmt.Sprintf("%s has logged in", requestPayload.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	token, err := app.Jwt.GenerateToken("admin@example.com")
	if err != nil {

		app.errorJSON(w, errors.New("error in generating token  "+err.Error()), http.StatusInternalServerError)
		return
	}

	var data struct {
		Auth any `json:"auth"`
		Data any `json:"data"`
	}
	data.Auth = token
	data.Data = user

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    data,
	}
	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) logRequest(name, data string) error {
	logServiceURL := "http://logger-service/log"

	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonRequest, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return err
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		var jsonFromService jsonResponse
		err := json.NewDecoder(response.Body).Decode(&jsonFromService)
		if err != nil {
			return err
		}
		return errors.New("jsonFromService.Error")
	}

	return nil
}

// func (app Config) verifyToken(email, token string) error {
// 	if token == "" {
// 		return errors.New("token not found")
// 	}
// 	claims, msg := app.Jwt.ValidateToken("access token", token)
// 	if msg != "" {
// 		if msg == "error while validation" {
// 			return errors.New("server validation error")
// 		}
// 		return errors.New(msg)
// 	}
// 	if claims.Email != email {
// 		return errors.New("token invalid")
// 	}
// 	return nil
// }
