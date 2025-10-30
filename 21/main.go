package main

import (
	"log"
	"net/http"

	"restapi/cache"
	"restapi/db"
	"restapi/handler"
	"restapi/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	ps, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	rc, err := cache.NewRedisCache()
	if err != nil {
		log.Fatal(err)
	}

	h, err := handler.NewHandler(ps, rc)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/register", h.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", h.LoginHandler).Methods("POST")

	api := r.NewRoute().Subrouter()
	api.Use(middleware.AuthorizationMiddleware)

	api.HandleFunc("/tasks", h.CreateTaskHandler).Methods("POST")
	api.HandleFunc("/tasks/{id:[0-9]+}", h.GetTaskHandler).Methods("GET")
	api.HandleFunc("/tasks", h.GetSelectedTasksHandler).Methods("GET")
	api.HandleFunc("/tasks/{id:[0-9]+}", h.UpdateTaskHandler).Methods("PUT")
	api.HandleFunc("/tasks/{id:[0-9]+}", h.DeleteTaskHandler).Methods("DELETE")
	api.HandleFunc("/tasks/{id:[0-9]+}/comments", h.AddCommentToTaskHandler).Methods("POST")

	log.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
