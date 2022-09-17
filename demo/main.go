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

type Templates struct {
	index  *template.Template
	errors *template.Template
}

type NotFound struct {
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, cat string) error {
	switch cat {
	case "index":
		return t.index.ExecuteTemplate(w, name, data)
	case "errors":
		return t.errors.ExecuteTemplate(w, name, data)
	default:
		return t.errors.ExecuteTemplate(w, name, data)
	}
}

func (n NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.Render(w, "error.html", "", "errors")
}

var t = &Templates{
	index:  template.Must(template.ParseFiles("views/index.html")),
	errors: template.Must(template.ParseFiles("views/error.html")),
}

// Index route
func Index(w http.ResponseWriter, r *http.Request) {
	t.Render(w, "index.html", "", "index")
}

type Error struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&Error{Error: err})
}

type AnalysisRequest struct {
	Data string `json:"data"`
}

func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	})
}

type AnalysisResponse struct {
	Data []*cirrus.Result `json:"data"`
}

func Analyze(w http.ResponseWriter, r *http.Request) {
	req := &AnalysisRequest{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(req)
	if err != nil {
		sendError(w, "Invalid response.")
		return
	}

	results, err := cirrus.Recognize(req.Data)
	if err != nil {
		sendError(w, err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)

	enc := json.NewEncoder(w)
	res := &AnalysisResponse{
		Data: results,
	}
	enc.Encode(res)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/analyze", Analyze).Methods("POST")

	r.NotFoundHandler = NotFound{}
	log.Fatal(http.ListenAndServe(":8000", r))
}
