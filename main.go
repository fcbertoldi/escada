package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

const (
	proxyPort          = "9982"
	googlebotUserAgent = "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.119 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
)

func extractUrl(u *url.URL) (reqUrl *url.URL, err error) {
	reqUrlString, err := url.QueryUnescape(u.RequestURI())
	if err != nil {
		return nil, err
	}
	if len(reqUrlString) > 0 && reqUrlString[0] == '/' {
		reqUrlString = reqUrlString[1:]
	}
	reqUrl, err = url.Parse(reqUrlString)
	if err != nil {
		return nil, err
	}

	return reqUrl, nil
}

func copyHeader(dst, src http.Header) {
	clear(dst)
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func ProxySite(w http.ResponseWriter, r *http.Request) {
	reqUrl, err := extractUrl(r.URL)
	if err != nil {
		slog.Error("ProxySite", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<h1>Invalid Origin URL</h1>")
		return
	}

	slog.Info("Proxied Request", slog.String("url", reqUrl.String()))

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		slog.Error("ProxySite", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
		return
	}
	req.Header.Set("User-Agent", googlebotUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("ProxySite", slog.String("error", err.Error()))
		if resp != nil {
			w.WriteHeader(resp.StatusCode)
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				slog.Error("ProxySite", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
			}
		}
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		slog.Error("ProxySite", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
	}
}

func main() {
	http.HandleFunc("GET /...", ProxySite)

	addr := fmt.Sprintf("127.0.0.1:%s", proxyPort)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("ListenAndServe:", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
