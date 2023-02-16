package main

import (
	"errors"
	"net/http"
)

// currently not using this middleware.

func (app Config) TokenAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			app.errorJSON(w, errors.New("token not found in the body"))
			return
		}
		var requestPayload struct {
			Email string `json:"email"`
			Data  string `json:"data"`
		}
		err := app.readJSON(w, r, &requestPayload)
		if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}
		claims, msg := app.Jwt.ValidateToken("access token", token)
		if msg != "" {
			if msg == "error while validation" {
				app.errorJSON(w, errors.New(msg), http.StatusInternalServerError)
				return
			}
			app.errorJSON(w, errors.New(msg), http.StatusUnauthorized)
			return
		}

		if claims.Email != requestPayload.Email {
			app.errorJSON(w, errors.New("token is invalid"))
			return
		}
		next.ServeHTTP(w, r)
	})

}
