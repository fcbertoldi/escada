package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	proxyPort          = "9982"
	googlebotUserAgent = "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.119 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
)

func extractUrl(u *url.URL) (reqUrl *url.URL, err error) {
	reqUrlString, err := url.QueryUnescape(u.RequestURI())
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	reqUrl, err = url.Parse(reqUrlString)
	if err != nil {
		return nil, err
	}

	return reqUrl, nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func ProxySite(w http.ResponseWriter, r *http.Request) {
	reqUrl, err := extractUrl(r.URL)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Proxying request to: ", reqUrl)

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	req.Header.Set("User-Agent", googlebotUserAgent)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	http.HandleFunc("GET /...", ProxySite)

	addr := fmt.Sprintf(":%s", proxyPort)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
