package main

import (
	"log"
	"logger/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (app Config) Ping(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "This is the logger-server",
	}
	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload jsonPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		return
	}
	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		log.Print(err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = app.writeJSON(w, http.StatusAccepted, resp)

}
