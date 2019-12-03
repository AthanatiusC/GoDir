package utils

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func OnErr(err error) {
	if err != nil {

	}
}

func ConnectMongoDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	iserror := ErrorHandler(err)
	if iserror {
		return nil
	}
	// Check theconnection
	err = client.Ping(context.TODO(), nil)
	iserror = ErrorHandler(err)
	if iserror {
		return nil
	}

	return client.Database("GoDir")
}

func ErrorHandler(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	} else {
		return false
	}
}

type Payload struct {
	Message string      `json:"message"`
	Data    interface{} `json:"returned_data"`
}

func WriteResult(res http.ResponseWriter, data interface{}, message string) {
	res.Header().Add("Access-Control-Allow-Origin", "*")
	(res).Header().Set("Access-Control-Allow-Headers", "*")
	(res).Header().Set("Access-Control-Allow-Methods", "*")
	res.Header().Set("Content-Type", "Application/JSON")

	var payload Payload
	payload.Message = message
	payload.Data = data
	result, _ := json.Marshal(payload)

	res.WriteHeader(http.StatusAccepted)
	res.Write([]byte(result))
	fmt.Println(message)
}
