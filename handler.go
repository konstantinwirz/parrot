package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
)

// Response creates a json encoded response object
// with following structure
// {
//	  "code": int
//    "message": string
// }
func Response(code int, message string) []byte {
	b, err := json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: message,
	})

	if err != nil {
		panic(err)
	}

	return b
}

// HandleHealth handles health requests
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleNotAllowed(w, r)
		return
	}

	writeResponse(w, 200, "OK")
}

// HandleNotAllowed handles not allowed request
func HandleNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, 405, "Method not allowed")
}

func HandleResources(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostResource(w, r)
	case http.MethodGet:
		handleGetResource(w, r)
	default:
		HandleNotAllowed(w, r)
	}
}

var resourcePattern = regexp.MustCompile("/(\\S+)/(\\S+)")
var resources = make(map[string]map[string][]byte)

func handlePostResource(w http.ResponseWriter, r *http.Request) {
	groups := resourcePattern.FindStringSubmatch(r.URL.EscapedPath())
	if len(groups) != 3 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "%s", Response(400, "Bad path, expected /resources/id"))
		return
	}

	id := groups[2]
	resource := groups[1]

	if r.Body == nil {
		writeResponse(w, 204, "No content")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "%s", Response(400, "expected a body"))
		return
	}

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", Response(500, err.Error()))
		return
	}

	if resources[resource] == nil {
		resources[resource] = make(map[string][]byte)
	}

	resources[resource][id] = body
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", Response(201, "created"))
}

func handleGetResource(w http.ResponseWriter, r *http.Request) {
	groups := resourcePattern.FindStringSubmatch(r.URL.EscapedPath())
	if len(groups) != 3 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "%s", Response(400, "Bad path, expected /resources/id"))
		return
	}

	resource, ok := resources[groups[1]]

	notFound := func() {
		w.WriteHeader(404)
		fmt.Fprintf(w, "%s", Response(404, "Not found"))
	}

	if !ok {
		notFound()
		return
	}

	content, ok := resource[groups[2]]
	if !ok {
		notFound()
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", content)
}

// WithHeaders is a HandlerFunc that set all the necessary headers
func WithHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Server", "Parrot ("+runtime.Version()+")")
		handler(w, r)
	}
}

func writeResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", Response(code, message))
}
