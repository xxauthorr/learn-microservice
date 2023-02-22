package main

import (
	"context"
	"log"
	"logger/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogINFO(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	})
	if err != nil {
		log.Println("error writing to mongo ")
		return err
	}

	*resp = "Processed payload via RPC : " + payload.Name
	return nil
}
