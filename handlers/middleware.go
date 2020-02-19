// Package handlers contains handlers and middleware used by the rest job worker.
package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func convertToken(t string) string {
	hashed := sha256.Sum256([]byte(t))
	slice := hashed[:]
	return base64.StdEncoding.EncodeToString(slice)
}

// Logging adds logging capabilities to incoming http requests.
func Logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Incoming request:", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

			defer func() {
				// TODO: figure out a way to log the status code
				logger.Println("Response returned:", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// VerifyAuth verifies that incoming requests are authenticated and authorized.
// It rejects requests with no or improper auth information by responding to requests with HTTP error codes.
// This middleware also injects a logger for handlers to log into.
func VerifyAuth(handle func(http.ResponseWriter, *http.Request, *log.Logger), logger *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// assume we have an authorized token (or a set of tokens) ready to be looked up
		authorizedToken := convertToken("123456")

		logger.Printf("Checking headers for auth ...\n")
		authHeader := ""
		if authHeader = strings.TrimSpace(r.Header.Get("Authorization")); authHeader == "" {
			logger.Printf("No Authorization header found.\n")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if hashedToken := convertToken(authHeader); hashedToken != authorizedToken {
			logger.Printf("Token is unauthorized.\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		logger.Printf("Token is authorized, proceeding to handling the request ...\n")
		handle(w, r, logger)
	}
}
