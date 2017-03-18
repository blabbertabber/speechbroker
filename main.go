package main

// curl -X POST -d "1234" http://diarizer.blabbertabber.com:8080/api/v1/upload

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"os"
	"path/filepath"
)

const PORT=8080
var dataDir = filepath.FromSlash("/opt/blabbertabber")

func handler(w http.ResponseWriter, r *http.Request) {
	conversationUUID := uuid.NewV4() // TODO(brian) handle error
	dataDir = filepath.Join(dataDir, conversationUUID.String())
	os.MkdirAll(dataDir, 0777)
	w.Write([]byte(conversationUUID.String()))
}

func main() {
	http.HandleFunc("/api/v1/upload", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
