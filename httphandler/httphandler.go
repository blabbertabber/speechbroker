package httphandler

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HttpHandler struct {
	Handler func(http.ResponseWriter, *http.Request)
}

func New() HttpHandler {
	return HttpHandler{
		Handler: handler,
	}
}

var soundRootDir = filepath.FromSlash("/var/blabbertabber/soundFiles/")
var resultsRootDir = filepath.FromSlash("/var/blabbertabber/diarizationResults/")

func handler(w http.ResponseWriter, r *http.Request) {
	diarizer := r.Header["Diarizer"]
	transcriber := r.Header["Transcriber"]
	fmt.Println("Diarizer: ", diarizer, "   Transcriber: ", transcriber)

	conversationUUID := uuid.New()
	meetingUuid := conversationUUID.String()
	soundDir := filepath.Join(soundRootDir, meetingUuid)
	err := os.MkdirAll(soundDir, 0777)
	if err != nil {
		log.Fatal("MkdirAll: ", err)
	}
	resultsDir := filepath.Join(resultsRootDir, meetingUuid)
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
	justTheHost := strings.Split(r.Host, ":")[0]
	w.Write([]byte(fmt.Sprintf("https://%s?meeting=%s", justTheHost, meetingUuid)))
	dst, err = os.Create(filepath.Join(resultsDir, "03_transcription_begun"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
	meetingWavFilepath := fmt.Sprintf("/blabbertabber/soundFiles/%s/meeting.wav", meetingUuid)
	diarizationFilepath := fmt.Sprintf("/blabbertabber/diarizationResults/%s/diarization.txt", meetingUuid)
	diarizationCommand := []string{
		"run",
		"--volume=/var/blabbertabber:/blabbertabber",
		"--workdir=/speaker-diarization",
		"blabbertabber/aalto-speech-diarizer",
		"/speaker-diarization/spk-diarization2.py",
		meetingWavFilepath,
		"-o",
		diarizationFilepath,
	}
	transcriptionFilepath := fmt.Sprintf("/blabbertabber/diarizationResults/%s/transcription.txt", meetingUuid)
	transcriptionCommand := []string{
		"run",
		"--volume=/var/blabbertabber:/blabbertabber",
		"blabbertabber/cmu-sphinx4-transcriber",
		"java",
		"-Xmx2g",
		"-cp",
		"/sphinx4-5prealpha-src/sphinx4-core/build/libs/sphinx4-core-5prealpha-SNAPSHOT.jar:/sphinx4-5prealpha-src/sphinx4-data/build/libs/sphinx4-data-5prealpha-SNAPSHOT.jar:.",
		"Transcriber",
		meetingWavFilepath,
		transcriptionFilepath,
	}
	go diarizeOrTranscribe("diarization", resultsDir, diarizationCommand...)
	go diarizeOrTranscribe("transcription", resultsDir, transcriptionCommand...)
}

func diarizeOrTranscribe(action string, resultsDir string, dockerCommandArgs ...string) {
	dst, err := os.Create(filepath.Join(resultsDir, action+"_begun"))
	log.Print(strings.Join(dockerCommandArgs, " "+"\n"))
	if err != nil {
		log.Fatal("Create: ", err)
	}
	dst.Close()
	command := exec.Command("docker", dockerCommandArgs...)
	stderr, err := command.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := command.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	log.Printf("%s\n", slurp)

	if err := command.Wait(); err != nil {
		log.Fatal(err)
	}
	dst, err = os.Create(filepath.Join(resultsDir, action+"_ended"))
}
