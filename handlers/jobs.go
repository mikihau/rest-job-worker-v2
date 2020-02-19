package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Hello returns a greeting message. This is a placeholder for the upcoming real handlers.
func Hello(w http.ResponseWriter, r *http.Request, l *log.Logger) {
	name := mux.Vars(r)["name"]
	l.Printf("Returning hello message: Hello %v!\n", name)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello %v!\n", name)))
}
