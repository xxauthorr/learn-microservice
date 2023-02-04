package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// response from the service
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// request.body
type RequestPayload struct {
	Action string        `json:"action"`
	Auth   AuthPayload   `json:"auth,omitempty"`
	Log    LoggerPayload `json:"logger,omitempty"`
}

// json body for action authenticate
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// json body for action logger
type LoggerPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "This is the broker-server",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

	// Everything writen in helper function writeJSON instead of this here
	// out, _ := json.MarshalIndent(payload, "", "\t")
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	switch requestPayload.Action {
	case "ping-auth":
		app.pingAuthService(w)
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "ping-logger":
		app.pingLogger(w)
	case "log":
		app.logItem(w, requestPayload.Log)
	default:
		log.Println("request payload: ", requestPayload)
		app.errorJSON(w, errors.New("unknown action"))
	}
}

// ping the authentication service

func (app *Config) pingAuthService(w http.ResponseWriter) {
	request, err := http.NewRequest("GET", "http://authentication-service/", nil)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("server returns error"), http.StatusInternalServerError)
		return
	}
	// read the response.body into a variable
	var jsonFromService jsonResponse
	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message))
		return
	}

	payload := jsonFromService
	app.writeJSON(w, http.StatusAccepted, payload)

}

// ping the logger service
func (app *Config) pingLogger(w http.ResponseWriter) {
	request, err := http.NewRequest("GET", "http://logger-service/", nil)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("server returns error"))
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message))
		return
	}
	payload := jsonFromService
	app.writeJSON(w, http.StatusAccepted, payload)
}

// get json and send it to the auth microservice
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	//call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()
	//make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("user is unauthorized"), response.StatusCode)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid credentials"), response.StatusCode)
		return
	}

	// read the response.body into a variable
	var jsonFromService jsonResponse
	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message))
		return
	}

	var payload jsonResponse
	payload.Error = jsonFromService.Error
	payload.Message = "authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItem(w http.ResponseWriter, l LoggerPayload) {
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err, response.StatusCode)
		return
	}

	var payload jsonResponse

	err = json.NewDecoder(request.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if payload.Error {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
