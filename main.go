package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"testing"
)

const (
	defaultAddr        = "127.0.0.1"
	defaultPort        = 9982
	googlebotUserAgent = "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.119 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
)

func extractUrl(u string, urlScheme string) (reqUrl *url.URL, err error) {
	reqUrlString, err := url.QueryUnescape(u)
	if err != nil {
		return nil, err
	}

	if len(reqUrlString) == 0 {
		return nil, fmt.Errorf("empty URL")
	}
	if reqUrlString[0] == '/' {
		reqUrlString = reqUrlString[1:]
	}
	if matched, _ := regexp.MatchString(`^https?://`, reqUrlString); !matched {
		reqUrlString = urlScheme + reqUrlString
	}
	reqUrl, err = url.Parse(reqUrlString)
	if err != nil {
		return nil, err
	}

	return reqUrl, nil
}

func copyResponseHeaders(dst, src http.Header) {
	clear(dst)
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

//go:embed public/index.html
var indexHTML []byte

func ProxyForm(w http.ResponseWriter, r *http.Request) {
	logger.Debug("ProxyForm", slog.String("url", r.URL.String()))
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(indexHTML); err != nil {
		logger.Error("ProxyForm", slog.String("error", err.Error()))
	}
}

type Handler struct {
	urlScheme string
}

func NewHandler() *Handler {
	if testing.Testing() {
		return &Handler{urlScheme: "http://"}
	}
	return &Handler{urlScheme: "https://"}
}

func (h *Handler) ProxyPage(w http.ResponseWriter, r *http.Request) {
	reqUrl, err := extractUrl(r.PathValue("page"), h.urlScheme)
	if err != nil {
		logger.Error("ProxyPage", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<h1>Invalid Origin URL</h1>")
		return
	}

	logger.Debug("Proxied Request", slog.String("url", reqUrl.String()))

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		logger.Error("ProxyPage", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
		return
	}

	for header, values := range r.Header {
		req.Header[header] = values
	}
	req.Header.Set("User-Agent", googlebotUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("ProxyPage", slog.String("error", err.Error()))
		if resp != nil {
			w.WriteHeader(resp.StatusCode)
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				logger.Error("ProxyPage", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
			}
		}
		return
	}
	defer resp.Body.Close()

	copyResponseHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		logger.Error("ProxyPage", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<h1>Internal Server Error</h1>")
	}
}

func main() {
	addrFlag := flag.String("addr", defaultAddr, "address to listen on")
	portFlag := flag.Int("port", defaultPort, "port to listen on")
	helpFlag := flag.Bool("help", false, "display help")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *portFlag <= 0 {
		logger.Error("Invalid port number", slog.Int("port", *portFlag))
		os.Exit(1)
	}

	http.HandleFunc("/pages/{page...}", NewHandler().ProxyPage)
	http.HandleFunc("/", ProxyForm)

	addr := fmt.Sprintf("%s:%d", *addrFlag, *portFlag)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error("ListenAndServe:", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
