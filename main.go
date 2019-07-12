package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/zeitgeistlabs/aether/internal/auth"
	"github.com/zeitgeistlabs/aether/internal/kubernetes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type InitialAuthResponse struct {
	Token string
}

type ExampleResponse struct {
	User string
}

// wrapper function for http logging
func logger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer log.Printf("%s - %s", r.Method, r.URL)
		fn(w, r)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	clientset := kubernetes.LoadClientSet()

	// validate the token from the request
	aetherToken, err := auth.InitialAuthentication(clientset, r.Header.Get("Token"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	resp := InitialAuthResponse{
		aetherToken,
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json marshaling error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

// Placeholder
func keyValueHandler(w http.ResponseWriter, r *http.Request) {
	user := auth.Authenticate(r.Header.Get("token"))
	if len(user) == 0 {
		w.WriteHeader(403)
		return
	}

	resp := ExampleResponse{
		user,
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json marshaling error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func main() {
	log.Printf("starting compute aether backend")

	// WARNING: path matching is strict, EG /auth/ doesn't match.
	r := mux.NewRouter()
	r.HandleFunc("/auth", logger(authHandler)).Methods("POST")
	r.HandleFunc("/keyvalue", logger(keyValueHandler))

	s := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Printf("shutdown signal received, exiting")

	s.Shutdown(context.Background())
}
