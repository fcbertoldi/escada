package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestProxySiteHandler_OK(t *testing.T) {
	t.Parallel()
    backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/test" {
            w.WriteHeader(http.StatusOK)
            if _, err := w.Write([]byte("TEST")); err != nil {
                t.Errorf("Failed to write response: %v", err)
            }
        } else {
            w.WriteHeader(http.StatusNotFound)
        }
    }))
    defer backend.Close()

    proxy := httptest.NewServer(http.HandlerFunc(ProxySite))
    defer proxy.Close()

    testURL := proxy.URL + "/" + url.QueryEscape(backend.URL + "/test")
	resp, err := http.Get(testURL)
    if err != nil {
        t.Fatalf("Failed to make request: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Failed to read response body: %v", err)
    }

    if string(body) != "TEST" {
        t.Errorf("Expected response body to be 'TEST', got '%s'", string(body))
    }
}

// TODO test proxied without http:// prefix

func TestProxySiteHandler_InvalidEncoding(t *testing.T) {
	t.Parallel()
	proxy := httptest.NewServer(http.HandlerFunc(ProxySite))
	defer proxy.Close()

	testURL := proxy.URL + "/" + url.QueryEscape("%%%%%%")
	resp, err := http.Get(testURL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code to be 400 Bad Request, got %d", resp.StatusCode)
	}
}

// TODO test for error from proxied request
