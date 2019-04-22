package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNotAllowedHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/any", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleNotAllowed)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusOK)
	}

	expected := Response(405, "Method not allowed")
	if !reflect.DeepEqual(rec.Body.Bytes(), expected) {
		t.Errorf("Handler returned unexpected body: got %s expected %s", rec.Body.Bytes(), expected)
	}

}

func TestHandleHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleHealth)

	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusOK)
	}

	expected := Response(405, "Method not allowed")
	if reflect.DeepEqual(rec.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %s expected %s", rec.Body.String(), expected)
	}
}

func TestGetNotExistingResource(t *testing.T) {
	req, err := http.NewRequest("GET", "/issues/25", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleResources)

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusNotFound)
	}

	expected := Response(404, "Not found")
	if reflect.DeepEqual(rec.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %s expected %s", rec.Body.String(), expected)
	}
}

func TestPostResourceWithoutBody(t *testing.T) {
	// create a resource with an empty body
	const url = "/issues/25"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleResources)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusNoContent)
	}

	// try to get the resource /issue/25
	// expecting a 404
	if st, _ := getResource(url); st != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v expected %v", st, http.StatusOK)
	}
}

func getResource(url string) (int, *bytes.Buffer) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleResources)
	handler.ServeHTTP(rec, req)

	return rec.Code, rec.Body
}

func TestGetExistingResource(t *testing.T) {
	// at first - create a resource
	req, err := http.NewRequest("POST", "/issues/25", strings.NewReader("{'a':'b'}"))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleResources)

	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusCreated)
	}

	expected := Response(201, "created")
	if reflect.DeepEqual(rec.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %s expected %s", rec.Body.String(), expected)
	}

	// now - get the resource
	req, err = http.NewRequest("GET", "/issues/25", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = httptest.NewRecorder()
	handler = http.HandlerFunc(HandleResources)

	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusOK)
	}

	body := "{'a':'b'}"
	if body != rec.Body.String() {
		t.Errorf("Handler returned unexpected body: got %s expected %s", rec.Body.String(), body)
	}
}
