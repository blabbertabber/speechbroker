package main

// curl -F "a=1234" https://test.diarizer.com:9443/api/v1/upload
// curl -F "meeting.wav=@/dev/null" http://test.diarizer.com:8080/api/v1/upload
// curl -F "meeting.wav=@/Users/cunnie/Google Drive/BlabberTabber/ICSI-diarizer-sample-meeting.wav" https://test.diarizer.com:9443/api/v1/upload
// curl --trace - -F "meeting.wav=@/dev/null" http://test.diarizer.com:8080/api/v1/upload
// cleanup: sudo -u diarizer find /var/blabbertabber -name "*-*-*" -exec rm -rf {} \;

import (
	"flag"
	"fmt"
	"github.com/blabbertabber/speechbroker/httphandler"
	"log"
	"net/http"
	"path/filepath"
)

const CLEAR_PORT = ":8080" // for troubleshooting in cleartext
const SSL_PORT = ":9443"

func main() {
	var IBMConfigFile = flag.String("IBMConfigFile", "",
		"pathname to JSON-formatted IBM Bluemix Watson Speech to Text service credentials")
	var keyPath = flag.String("keyPath", "/etc/pki/nginx/private/server.key", "path to HTTPS private key")
	var certPath = flag.String("certPath", "/etc/pki/nginx/server.crt", "path to HTTPS certificate")

	flag.Parse()

	fmt.Printf("I got these creds: %s", *IBMConfigFile)

	h := httphandler.Handler{
		Uuid:           httphandler.UuidReal{},
		FileSystem:     httphandler.FileSystemReal{},
		DockerRunner:   httphandler.DockerRunnerReal{},
		SoundRootDir:   filepath.FromSlash("/var/blabbertabber/soundFiles/"),
		ResultsRootDir: filepath.FromSlash("/var/blabbertabber/diarizationResults/"),
	}
	http.HandleFunc("/api/v1/upload", h.ServeHTTP)

	go func() {
		log.Fatal(http.ListenAndServe(CLEAR_PORT, nil))

	}()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(SSL_PORT), *certPath, *keyPath, nil))
}
