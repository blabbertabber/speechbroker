package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const diarizerNonoIo = "10.2.0.99"

func handler(writer http.ResponseWriter, request *http.Request) {
	resp, err := http.Get(proxyEndpoint(request.URL))
	if err != nil {
		fmt.Println(err)
	}

	writer.WriteHeader(resp.StatusCode)
	io.Copy(writer, resp.Body)
}

func proxyEndpoint(url *url.URL) string {
	path := url.Path[1:]
	return fmt.Sprintf("http://%s/%s", diarizerNonoIo, path)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
