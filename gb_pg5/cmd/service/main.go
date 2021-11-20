package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"


func main()  {
	srv, err := initServer()
	if err !=nil{
		log.Fatalf("[ERR]: failed to initialize server: %v", err)
	}
	log.Println("Let's Go!")
	if err := srv.ListenAndServe(); err != nil{
		log.Fatalf("[ERR]: %v", err)
	}
}

func initServer() (*http.Server, error) {
	handler, err := registerRoutes()
	if err != nil {
		return nil, fmt.Errorf("failed to register routes: %w", err)
	}
	srv := &http.Server{
		Addr: ":8080",
		Handler: handler,
	}
	return srv, nil
}

func registerRoutes() (http.Handler, error) {
	r := mux.NewRouter()
	r.HandleFunc("/quantity/{product_id}", func(w http.ResponseWriter, r *http.Request){
		product_idHint.GetQuantityByProduct_id(w, r, mux.Vars(r)["product_id"])
	}).Methods("GET")
	addDBMiddleware, err := createAddDBMiddleware()
	if err != nil {
		return nil, fmt.Errorf("failed to create AddDBMiddleware: %w", err)
	}
	r.Use(addDBMiddleware)
	return r, err
}

func createAddDBMiddleware() (func(next http.Handler) http.Handler, error) {
	connStr, err := getConnString()
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string info for DB connection: %w", err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			db, err := storage.NewDB(connStr)
			if err != nil {
				log.Println("[ERR]: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer db.Close()
			r = r.WithContext(context.WithValue(r.Context(), storage.ContextKeyDB, db))

			next.ServeHTTP(w, r)
		})
	}, nil
}