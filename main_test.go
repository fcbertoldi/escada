package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func proxiedBackendServer(userAgentSpy *string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*userAgentSpy = r.Header.Get("User-Agent")
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

func TestProxyPageHandler_OK(t *testing.T) {
	t.Parallel()
	userAgentSpy := ""
	backend := proxiedBackendServer(&userAgentSpy)
	defer backend.Close()

	req := httptest.NewRequest(http.MethodGet, "/pages/", nil)
	req.SetPathValue("page", url.QueryEscape(backend.URL+"/test"))
	recorder := httptest.NewRecorder()

	NewHandler().ProxyPage(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code to be 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "TEST" {
		t.Errorf("Expected response body to be 'TEST', got '%s'", string(body))
	}

	if userAgentSpy != googlebotUserAgent {
		t.Errorf("Expected User-Agent to be '%s', got '%s'", googlebotUserAgent, userAgentSpy)
	}
}

func TestProxyPageHandler_NoHttpScheme(t *testing.T) {
	t.Parallel()
	userAgentSpy := ""
	backend := proxiedBackendServer(&userAgentSpy)
	defer backend.Close()

	req := httptest.NewRequest(http.MethodGet, "/pages/", nil)
	req.SetPathValue("page", url.QueryEscape(backend.URL[len("http://"):]+"/test"))
	recorder := httptest.NewRecorder()

	NewHandler().ProxyPage(recorder, req)

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

func TestProxyPageHandler_InvalidEncoding(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/pages/", nil)
	req.SetPathValue("page", url.QueryEscape("%%%%%%"))
	recorder := httptest.NewRecorder()

	NewHandler().ProxyPage(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code to be 400 Bad Request, got %d", resp.StatusCode)
	}
}

func TestProxyPageHandler_ProxiedError(t *testing.T) {
	t.Parallel()
	userAgentSpy := ""
	backend := proxiedBackendServer(&userAgentSpy)
	defer backend.Close()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("page", url.QueryEscape(backend.URL+"/410"))
	recorder := httptest.NewRecorder()

	NewHandler().ProxyPage(recorder, req)

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
