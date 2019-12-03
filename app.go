package main

import (
	"github.com/athanatius/godir"
	"log"
	"net/http"

	// "encoding/json"
	directory "github.com/athanatius/godir/controllers/directory"
	users "github.com/athanatius/godir/controllers/users"
	// "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	// "github.com/rs/cors"
)

var largePool chan func()

func main() {
	utils.ConnectMongoDB() //Test Connection
	router := mux.NewRouter()
	// headers := handlers.AllowedHeaders([]string{"X-Requested-With"})
	// origins := handlers.AllowedOrigins([]string{"*"})
	// methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.Headers("Content-Type", "Application/JSON")

	prefix := router.PathPrefix("/v1/").Subrouter()
	prefix.HandleFunc("/users", users.GetAllUsers).Methods("GET", "OPTIONS")
	prefix.HandleFunc("/users", users.CreateUsers).Methods("POST", "OPTIONS")
	prefix.HandleFunc("/auth", users.Auth).Methods("POST", "OPTIONS")
	prefix.HandleFunc("/users/{id}", users.DeleteUsers).Methods("DELETE", "OPTIONS")

	prefix.HandleFunc("/directory", directory.GetDirectory).Methods("POST", "OPTIONS")
	prefix.HandleFunc("/directory/delete", directory.DeleteDirectory).Methods("POST", "OPTIONS")
	prefix.HandleFunc("/directory/create", directory.CreateFolder).Methods("POST", "OPTIONS")
	prefix.HandleFunc("/directory/upload", directory.UploadFile).Methods("POST", "OPTIONS")
	prefix.HandleFunc("/directory/rename", directory.RenameFolder).Methods("PUT", "OPTIONS")
	prefix.HandleFunc("/directory/download", directory.DownloadFile).Methods("GET", "OPTIONS")
	// prefix.HandleFunc("/directory/zip", directory.ZipFile).Methods("PUT", "OPTIONS")

	router.Use(mux.CORSMethodMiddleware(router))
	log.Println("Connection Successfull! Api running at http://localhost:9000")
	defer log.Println("Connection Closed")

	log.Fatal(http.ListenAndServe(":9000", router))
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			//handle preflight in here
		} else {
			h.ServeHTTP(w, r)
		}
	}
}
