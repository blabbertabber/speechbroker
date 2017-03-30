package main

// curl -X POST -F "a=1234" https://diarizer.blabbertabber.com:9443/api/v1/upload
// curl -X POST -F "meeting.wav=@/Users/cunnie/Google Drive/BlabberTabber/ICSI-diarizer-sample-meeting.wav" https://diarizer.blabbertabber.com:9443/api/v1/upload
// cleanup: sudo -u diarizer find /var/blabbertabber -name "*-*-*" -exec rm -rf {} \;

import (
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const PORT = 9443

var soundRootDir = filepath.FromSlash("/var/blabbertabber/soundFiles/")
var resultsRootDir = filepath.FromSlash("/var/blabbertabber/diarizationResults/")
var keyPath = filepath.FromSlash("/etc/pki/nginx/private/diarizer.blabbertabber.com.key")
var certPath = filepath.FromSlash("/etc/pki/nginx/diarizer.blabbertabber.com.crt")

func handler(w http.ResponseWriter, r *http.Request) {
	conversationUUID := uuid.NewV4()
	uuid := conversationUUID.String()
	soundDir := filepath.Join(soundRootDir, uuid)
	err := os.MkdirAll(soundDir, 0777)
	if err != nil {
		log.Fatal("MkdirAll: ", err)
	}
	resultsDir := filepath.Join(resultsRootDir, uuid)
	err = os.MkdirAll(resultsDir, 0777)
	if err != nil {
		log.Fatal("MkdirAll: ", err)
	}
	reader, err := r.MultipartReader()
	if err != nil {
		log.Fatal("MultipartReader: ", err)
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		dst, err := os.Create(filepath.Join(soundDir, "meeting.wav"))
		if err != nil {
			log.Fatal("Create: ", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, part); err != nil {
			log.Fatal("Copy: ", err)
		}
	}
	// return weblink to client "https://diarizer.blabbertabber.com/UUID"
	w.Write([]byte(fmt.Sprint("https://diarizer.blabbertabber.com/", uuid)))
	// kick off diarization in the background
	diarizationCommand := exec.Command("docker",
		"run",
		"--volume=/var/blabbertabber:/blabbertabber",
		"--workdir=/speaker-diarization",
		"blabbertabber/aalto-speech-diarizer",
		"/speaker-diarization/spk-diarization2.py",
		fmt.Sprintf("/blabbertabber/soundFiles/%s/meeting.wav", uuid),
		"-o", fmt.Sprintf("/blabbertabber/diarizationResults/%s/results.txt", uuid))
	err = diarizationCommand.Run()
	if err != nil {
		log.Fatal("Run: ", err)
	}
}

func main() {
	http.HandleFunc("/api/v1/upload", handler)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%d", PORT), certPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
