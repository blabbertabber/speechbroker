package main

// curl -F "a=1234" https://diarizer.blabbertabber.com:9443/api/v1/upload
// curl -F "meeting.wav=@/Users/cunnie/Google Drive/BlabberTabber/ICSI-diarizer-sample-meeting.wav" https://test.diarizer.com:9443/api/v1/upload
// curl --trace - -F "meeting.wav=@/dev/null" http://diarizer.blabbertabber.com:8080/api/v1/upload
// cleanup: sudo -u diarizer find /var/blabbertabber -name "*-*-*" -exec rm -rf {} \;

import (
	"fmt"
	"github.com/blabbertabber/speechbroker/httphandler"
	"log"
	"net/http"
	"path/filepath"
)

const CLEAR_PORT = ":8080" // for troubleshooting in cleartext
const SSL_PORT = ":9443"

var keyPath = filepath.FromSlash("/etc/pki/nginx/private/server.key")
var certPath = filepath.FromSlash("/etc/pki/nginx/server.crt")

func main() {
	h := httphandler.HttpHandler{
		Uuid:           httphandler.UuidReal{},
		DockerRunner:   httphandler.DockerRunnerReal{},
		SoundRootDir:   filepath.FromSlash("/var/blabbertabber/soundFiles/"),
		ResultsRootDir: filepath.FromSlash("/var/blabbertabber/diarizationResults/"),
	}
	http.HandleFunc("/api/v1/upload", h.Handler)

	go func() {
		log.Fatal(http.ListenAndServe(CLEAR_PORT, nil))

	}()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(SSL_PORT), certPath, keyPath, nil))
}
