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

// `counterfeiter httphandler/httphandler.go Uuid`
type Uuid interface {
	GetUuid() string
}

type UuidReal struct{}

func (u UuidReal) GetUuid() string {
	return uuid.New().String()
}

// `counterfeiter httphandler/httphandler.go FileSystem`
type FileSystem interface {
	MkdirAll(string, os.FileMode) error
	Create(string) (*os.File, error)
	Copy(io.Writer, io.Reader) (int64, error)
}

type FileSystemReal struct{}

func (FileSystemReal) MkdirAll(path string, mode os.FileMode) error {
	return os.MkdirAll(path, mode)
}

func (FileSystemReal) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (FileSystemReal) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

// `counterfeiter httphandler/httphandler.go DockerRunner`
type DockerRunner interface {
	Run(action string, resultsDir string, dockerCommandArgs ...string)
}

type DockerRunnerReal struct{}

func (d DockerRunnerReal) Run(action string, resultsDir string, dockerCommandArgs ...string) {
	dst, err := os.Create(filepath.Join(resultsDir, action+"_begun"))
	log.Print(strings.Join(dockerCommandArgs, " "+"\n"))
	if err != nil {
		panic("Create: " + err.Error())
	}
	dst.Close()
	command := exec.Command("docker", dockerCommandArgs...)
	stderr, err := command.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := command.Start(); err != nil {
		panic(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	log.Printf("%s\n", slurp)

	if err := command.Wait(); err != nil {
		panic(err)
	}
	dst, err = os.Create(filepath.Join(resultsDir, action+"_ended"))
}

type Handler struct {
	Uuid           Uuid
	FileSystem     FileSystem
	DockerRunner   DockerRunner
	SoundRootDir   string
	ResultsRootDir string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// `diarizer` and `transcriber` _must_ match what the Android client sends
	// e.g.  https://github.com/blabbertabber/blabbertabber/blob/fea98684ad500380ef347cff584821ee52098c1e/app/src/main/java/com/blabbertabber/blabbertabber/RecordingActivity.java#L386-L392
	diarizer := r.Header["Diarizer"]       // "IBM" or "Aalto"
	transcriber := r.Header["Transcriber"] // "IBM" or "CMUSphinx4"
	fmt.Println("Diarizer: ", diarizer, "   Transcriber: ", transcriber)

	meetingUuid := h.Uuid.GetUuid()
	soundDir := filepath.Join(h.SoundRootDir, meetingUuid)
	err := h.FileSystem.MkdirAll(soundDir, 0777)
	if err != nil {
		panic(err.Error())
	}
	resultsDir := filepath.Join(h.ResultsRootDir, meetingUuid)
	err = h.FileSystem.MkdirAll(resultsDir, 0777)
	if err != nil {
		panic(err.Error())
	}
	dst, err := h.FileSystem.Create(filepath.Join(resultsDir, "00_upload_begun"))
	if err != nil {
		panic(err.Error())
	}
	dst.Close()
	reader, err := r.MultipartReader()
	if err != nil {
		panic(err.Error())
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		dst, err := h.FileSystem.Create(filepath.Join(soundDir, "meeting.wav"))
		if err != nil {
			panic(err.Error())
		}
		defer dst.Close()

		if _, err := h.FileSystem.Copy(dst, part); err != nil {
			panic(err.Error())
		}
	}
	dst, err = h.FileSystem.Create(filepath.Join(resultsDir, "01_upload_finished"))
	if err != nil {
		panic(err.Error())
	}
	dst.Close()
	// return weblink to client "https://diarizer.blabbertabber.com/UUID"
	justTheHost := strings.Split(r.Host, ":")[0]
	w.Write([]byte(fmt.Sprintf("https://%s?meeting=%s", justTheHost, meetingUuid)))
	dst, err = h.FileSystem.Create(filepath.Join(resultsDir, "03_transcription_begun"))
	if err != nil {
		panic(err.Error())
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
	go h.DockerRunner.Run("diarization", resultsDir, diarizationCommand...)
	go h.DockerRunner.Run("transcription", resultsDir, transcriptionCommand...)
}
