package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/amku91/lam/api/app/order"
	"github.com/amku91/lam/mongo"
	"github.com/rs/cors"
	"github.com/amku91/lam/api/config"
)

func main() {

	var (
		err error
	)

	//Initialize the routes
	r := chi.NewRouter()

	// CORS config
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Initialise the router middlewares

	r.Use(corsConfig.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Lam Ping Successful"))
	})

	r.Mount("/order", order.Controller{}.Routes())

	//Init system functions
	initialise()

	log.Println("Starting Application...")

	log.Println("Application will listen at port 8080")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println("Error while initializing the Application: " + err.Error())
		return
	}
}

//All initialisation related function calls goes here
func initialise() {

	initDB()
}

//Make Mongo DB Connection
func initDB() {

	mongo.MaxPool = config.MONGO_MAX_POOL
	mongo.PATH = config.MONGO_DSN
	mongo.DBNAME = config.MONGO_DATABASE
	mongo.CheckAndInitServiceConnection()
}
