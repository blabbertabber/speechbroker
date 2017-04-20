package main

// curl -F "a=1234" https://diarizer.blabbertabber.com:9443/api/v1/upload
// curl -F "meeting.wav=@/Users/cunnie/Google Drive/BlabberTabber/ICSI-diarizer-sample-meeting.wav" https://diarizer.blabbertabber.com:9443/api/v1/upload
// curl --trace - -F "meeting.wav=@/dev/null" http://diarizer.blabbertabber.com:8080/api/v1/upload
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

const CLEAR_PORT = ":8080" // for troubleshooting in cleartext
const SSL_PORT = ":9443"

var soundRootDir = filepath.FromSlash("/var/blabbertabber/soundFiles/")
var resultsRootDir = filepath.FromSlash("/var/blabbertabber/diarizationResults/")
var keyPath = filepath.FromSlash("/etc/pki/nginx/private/diarizer.blabbertabber.com.key")
var certPath = filepath.FromSlash("/etc/pki/nginx/diarizer.blabbertabber.com.crt")

func handler(w http.ResponseWriter, r *http.Request) {
	// NewV1() works via timestamp, which I like. Has mutex to avoid collisions
	conversationUUID := uuid.NewV1()
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
	dst, err := os.Create(filepath.Join(resultsDir, "00_upload_begun"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
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
	dst, err = os.Create(filepath.Join(resultsDir, "01_upload_finished"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
	// return weblink to client "https://diarizer.blabbertabber.com/UUID"
	w.Write([]byte(fmt.Sprint("https://diarizer.blabbertabber.com/", uuid)))
	dst, err = os.Create(filepath.Join(resultsDir, "03_transcription_begun"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
	go diarize(resultsDir, uuid)
	go transcribe(resultsDir, uuid)
}

func diarize(resultsDir string, uuid string) {
	dst, err := os.Create(filepath.Join(resultsDir, "diarization_begun"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
	diarizationCommand := exec.Command("docker",
		"run",
		"--volume=/var/blabbertabber:/blabbertabber",
		"--workdir=/speaker-diarization",
		"blabbertabber/aalto-speech-diarizer",
		"/speaker-diarization/spk-diarization2.py",
		fmt.Sprintf("/blabbertabber/soundFiles/%s/meeting.wav", uuid),
		"-o", fmt.Sprintf("/blabbertabber/diarizationResults/%s/results.txt", uuid))
	err = diarizationCommand.Run()
	dst, err = os.Create(filepath.Join(resultsDir, "diarization_finished"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
}

func transcribe(resultsDir string, uuid string) {
	dst, err := os.Create(filepath.Join(resultsDir, "transcription_begun"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
	transcriptionCommand := exec.Command("docker",
		"run",
		"--volume=/var/blabbertabber:/blabbertabber",
		"--workdir=/speaker-diarization",
		"blabbertabber/aalto-speech-diarizer",
		"/speaker-diarization/spk-diarization2.py",
		fmt.Sprintf("/blabbertabber/soundFiles/%s/meeting.wav", uuid),
		"-o", fmt.Sprintf("/blabbertabber/diarizationResults/%s/results.txt", uuid))
	err = transcriptionCommand.Run()
	dst, err = os.Create(filepath.Join(resultsDir, "transcription_finished"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
}

func main() {
	http.HandleFunc("/api/v1/upload", handler)

	go func() {
		log.Fatal(http.ListenAndServe(CLEAR_PORT, nil))

	}()
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf( SSL_PORT), certPath, keyPath, nil))
}
