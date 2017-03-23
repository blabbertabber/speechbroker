package main

// curl -X POST -d "1234" http://diarizer.blabbertabber.com:8080/api/v1/upload

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"os"
	"path/filepath"
	"io/ioutil"
)

const PORT = 8080

var dataRootDir = filepath.FromSlash("/opt/blabbertabber")

func handler(w http.ResponseWriter, r *http.Request) {
	conversationUUID := uuid.NewV4()
	dataDir := filepath.Join(dataRootDir, conversationUUID.String())
	_ = os.MkdirAll(dataDir, 0777)                                           // TODO(brian) handle error
	bytes, _ := ioutil.ReadAll(r.Body)                                       // TODO(brian) handle error
	_ = ioutil.WriteFile(filepath.Join(dataDir, "meeting.wav"), bytes, 0644) // TODO(brian) handle error
	// return weblink to client "http://diarizer.blabbertabber.com:8080/results/UUID"
	w.Write([]byte(fmt.Sprint("http://diarizer.blabbertabber.com:8080/results/", conversationUUID.String())))
	// kick off diarization in the background
}

func main() {
	http.HandleFunc("/api/v1/upload", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
