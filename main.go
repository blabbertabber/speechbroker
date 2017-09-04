package main

// curl -F "a=1234" https://test.diarizer.com:9443/api/v1/upload
// curl -F "meeting.wav=@/dev/null" http://test.diarizer.com:8080/api/v1/upload
// curl -F "meeting.wav=@/Users/cunnie/Google Drive/BlabberTabber/ICSI-diarizer-sample-meeting.wav" https://test.diarizer.com:9443/api/v1/upload
// curl --trace - -F "meeting.wav=@/dev/null" http://test.diarizer.com:8080/api/v1/upload
// cleanup: sudo -u diarizer find /var/blabbertabber -name "*-*-*" -exec rm -rf {} \;

import (
	"flag"
	"fmt"
	"github.com/blabbertabber/speechbroker/cmdrunner"
	"github.com/blabbertabber/speechbroker/diarizerrunner"
	"github.com/blabbertabber/speechbroker/httphandler"
	"github.com/blabbertabber/speechbroker/ibmservicecreds"
	"github.com/blabbertabber/speechbroker/setdockergroup"
	"github.com/blabbertabber/speechbroker/speedfactors"
	"log"
	"net/http"
	"path/filepath"
)

const CLEAR_PORT = ":8080" // for troubleshooting in cleartext
const SSL_PORT = ":9443"

func main() {
	log.Println("speechbroker started.")
	var ibmServiceCredsPath = flag.String("ibmServiceCredsPath", "",
		"pathname to JSON-formatted IBM Bluemix Watson Speech to Text service credentials")
	var speedfactorsPath = flag.String("speedfactorsPath", "",
		"pathname to JSON-formatted hash of speech processing time factors")
	var keyPath = flag.String("keyPath", "/etc/pki/nginx/private/server.key", "path to HTTPS private key")
	var certPath = flag.String("certPath", "/etc/pki/nginx/server.crt", "path to HTTPS certificate")

	flag.Parse()

	ibmServiceCreds, err := ibmservicecreds.ReadCredsFromPath(*ibmServiceCredsPath)
	if err != nil {
		panic("I couldn't read the IBM service creds: " + *ibmServiceCredsPath + ", error: " + err.Error())
	}
	speedfactors, err := speedfactors.ReadCredsFromPath(*speedfactorsPath)
	if err != nil {
		panic("I couldn't read the Speech processing time factors: " + *speedfactorsPath + ", error: " + err.Error())
	}
	// must be in `docker` group to run containers for docker-ce
	setdockergroup.SetDockerGroup()

	h := httphandler.Handler{
		IBMServiceCreds: ibmServiceCreds,
		Speedfactors:    speedfactors,
		Uuid:            httphandler.UuidReal{},
		FileSystem:      httphandler.FileSystemReal{},
		Runner: diarizerrunner.Runner{
			CmdRunner: cmdrunner.CmdRunnerReal{},
		},
		SoundRootDir:   filepath.FromSlash("/var/blabbertabber/soundFiles/"),
		ResultsRootDir: filepath.FromSlash("/var/blabbertabber/diarizationResults/"),
	}
	http.HandleFunc("/api/v1/upload", h.ServeHTTP)

	go func() {
		log.Fatal(http.ListenAndServe(CLEAR_PORT, nil))

	}()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(SSL_PORT), *certPath, *keyPath, nil))
}
