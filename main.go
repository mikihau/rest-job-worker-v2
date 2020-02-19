package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/mikihau/rest-job-worker-v2/handlers"
)

// This http server is adapted from: https://gist.github.com/enricofoltran/10b4a980cd07cb02836f70a4ab3e72d7.
func main() {
	tlsMode := true

	listenAddr := ":8080"
	if tlsMode {
		listenAddr = ":443"
	}

	var logger = log.New(os.Stdout, "", log.LstdFlags|log.LUTC|log.Lshortfile)
	logger.Printf("Server is starting...")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello/{name}", handlers.VerifyAuth(handlers.Hello, logger)).Methods(http.MethodGet)

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      handlers.Logging(logger)(router),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Printf("Server is shutting down...")

		// shut down the server
		ctxServer, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctxServer); err != nil {
			logger.Printf("Failed to gracefully shutdown the server: %v\n", err)
		}

		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	var err error
	if tlsMode {
		err = server.ListenAndServeTLS("tls/ca.pem", "tls/ca.key")
	} else {
		// reference to setting up self-signed ssl certificate:
		// https://devcenter.heroku.com/articles/ssl-certificate-self
		// https://gist.github.com/mtigas/952344
		err = server.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Printf("Server stopped")
}
