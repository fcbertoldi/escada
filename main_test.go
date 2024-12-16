package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func proxiedBackendServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test" {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("TEST")); err != nil {
				log.Printf("Failed to write response: %v", err)
			}
		} else if r.URL.Path == "/410" {
			w.WriteHeader(http.StatusGone)
			if _, err := w.Write([]byte("GONE")); err != nil {
				log.Printf("Failed to write response: %v", err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))

}

func TestProxySiteHandler_OK(t *testing.T) {
	t.Parallel()
	backend := proxiedBackendServer()
	defer backend.Close()

	testPath := "/" + url.QueryEscape(backend.URL+"/test")
	log.Printf("Test path: %s", testPath)
	req := httptest.NewRequest(http.MethodGet, testPath, nil)
	recorder := httptest.NewRecorder()

	NewHandler().ProxySite(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "TEST" {
		t.Errorf("Expected response body to be 'TEST', got '%s'", string(body))
	}
}

func TestProxySiteHandler_NoHttpScheme(t *testing.T) {
	t.Parallel()
	backend := proxiedBackendServer()
	defer backend.Close()

	testPath := "/" + url.QueryEscape(backend.URL[len("http://"):]+"/test")

	log.Printf("Test path: %s", testPath)
	req := httptest.NewRequest(http.MethodGet, testPath, nil)
	recorder := httptest.NewRecorder()

	NewHandler().ProxySite(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "TEST" {
		t.Errorf("Expected response body to be 'TEST', got '%s'", string(body))
	}
}

func TestProxySiteHandler_InvalidEncoding(t *testing.T) {
	t.Parallel()
	testPath := "/" + url.QueryEscape("%%%%%%")
	req := httptest.NewRequest(http.MethodGet, testPath, nil)
	recorder := httptest.NewRecorder()

	NewHandler().ProxySite(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code to be 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestProxySiteHandler_ProxiedError(t *testing.T) {
	t.Parallel()
	backend := proxiedBackendServer()
	defer backend.Close()

	testPath := "/" + url.QueryEscape(backend.URL+"/410")
	log.Printf("Test path: %s", testPath)
	req := httptest.NewRequest(http.MethodGet, testPath, nil)
	recorder := httptest.NewRecorder()

	NewHandler().ProxySite(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusGone {
		t.Errorf("Expected status code to be 410 Gone, got %d", resp.StatusCode)
	}

	if string(body) != "GONE" {
		t.Errorf("Expected response body to be 'TEST', got '%s'", string(body))
	}
}
