package main

// curl -X POST -d "1234" http://diarizer.blabbertabber.com:8080/api/v1/upload

import (
	"fmt"
	"net/http"

	"github.com/satori/go.uuid"
	"os"
	"path/filepath"
	"io/ioutil"
	"crypto/tls"
	"time"
	"log"
)

const PORT = 9443

var dataRootDir = filepath.FromSlash("/var/blabbertabber/UploadServer/")
var keyPath = filepath.FromSlash("/etc/pki/nginx/private/diarizer.blabbertabber.com.key")
var certPath = filepath.FromSlash("/etc/pki/nginx/diarizer.blabbertabber.com.crt")

func handler(w http.ResponseWriter, r *http.Request) {
	conversationUUID := uuid.NewV4()
	dataDir := filepath.Join(dataRootDir, conversationUUID.String())
	_ = os.MkdirAll(dataDir, 0777)                                           // TODO(brian) handle error
	bytes, _ := ioutil.ReadAll(r.Body)                                       // TODO(brian) handle error
	_ = ioutil.WriteFile(filepath.Join(dataDir, "meeting.wav"), bytes, 0644) // TODO(brian) handle error
	// return weblink to client "http://diarizer.blabbertabber.com:8080/results/UUID"
	w.Write([]byte(fmt.Sprint("https://diarizer.blabbertabber.com/", conversationUUID.String())))
	// kick off diarization in the background
}

func main() {
	//&tls.Config{
	//	// Causes servers to use Go's default ciphersuite preferences,
	//	// which are tuned to avoid attacks. Does nothing on clients.
	//	PreferServerCipherSuites: true,
	//	// Only use curves which have assembly implementations
	//	CurvePreferences: []tls.CurveID{
	//		tls.CurveP256,
	//		tls.X25519, // Go 1.8 only
	//	},
	//	// If you can take the compatibility loss of the Modern configuration, you should then also set MinVersion and CipherSuites.
	//	MinVersion: tls.VersionTLS12,
	//	CipherSuites: []uint16{
	//		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	//		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	//		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
	//		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
	//		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	//		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	//
	//		// Best disabled, as they don't provide Forward Secrecy,
	//		// but might be necessary for some clients
	//		// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	//		// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	//	},
	//}

	//srv := &http.Server{
	//	ReadTimeout:  5 * time.Second,
	//	WriteTimeout: 5 * time.Second,
	//	Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	//		w.Header().Set("Connection", "close")
	//		url := "https://" + req.Host + req.URL.String()
	//		http.Redirect(w, req, url, http.StatusMovedPermanently)
	//	}),
	//}

	http.HandleFunc("/api/v1/upload", handler)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", PORT), certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
