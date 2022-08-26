package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
	"unicode"

	"github.com/gorilla/mux"
	"github.com/lacunaverse/cirrus"
)
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/analyze", Analyze).Methods("POST")

	r.NotFoundHandler = NotFound{}
	log.Fatal(http.ListenAndServe(":8000", r))
}
