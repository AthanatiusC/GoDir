package main

import (
	"github.com/athanatius/godir"
	"log"
	"net/http"

	"encoding/json"
	directory "github.com/athanatius/godir/controllers/directory"
	users "github.com/athanatius/godir/controllers/users"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var largePool chan func()

func main() {
	utils.ConnectMongoDB() //Test Connection
	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.Headers("Content-Type", "application/json")

	router.HandleFunc("/users", users.GetAllUsers).Methods("GET")
	router.HandleFunc("/users", users.CreateUsers).Methods("POST")
	router.HandleFunc("/auth", users.Auth).Methods("POST")
	router.HandleFunc("/users/{id}", users.DeleteUsers).Methods("DELETE")

	router.HandleFunc("/directory", directory.GetDirectory).Methods("POST")
	router.HandleFunc("/directory/delete", directory.DeleteDirectory).Methods("POST")

	log.Println("Connection Successfull! Api running at http://localhost:9000")
	defer log.Println("Connection Closed")
	log.Panic(http.ListenAndServe(":9000", handlers.CORS(origins, headers, methods)(router)))
}

func handler1(w http.ResponseWriter, r *http.Request) {
	var job struct{ URL string }

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		// Block until there are fewer than cap(largePool) light-work
		// goroutines running.
		// largePool <- struct{}{}
		defer func() { <-largePool }() // Let everyone that we are done

		http.Get(job.URL)
	}()

	w.WriteHeader(http.StatusAccepted)
}
