// httphandler handles web server requests (POSTs) and invokes
// back-end speech processors and returns a URL of results to the client.
package httphandler

import (
	"fmt"
	"github.com/blabbertabber/speechbroker/diarizerrunner"
	"github.com/blabbertabber/speechbroker/ibmservicecreds"
	"github.com/blabbertabber/speechbroker/speedfactors"
	"github.com/blabbertabber/speechbroker/timesandsize"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// `counterfeiter httphandler/httphandler.go Uuid`
type Uuid interface {
	GetUuid() string
}

type UuidReal struct{}

func (u UuidReal) GetUuid() string {
	return uuid.New().String()
}

// `counterfeiter  -o httphandler/httphandlerfakes/fake_file_info.go os.FileInfo`
// `counterfeiter httphandler/httphandler.go FileSystem`
type FileSystem interface {
	MkdirAll(string, os.FileMode) error
	Create(string) (*os.File, error)
	Copy(io.Writer, io.Reader) (int64, error)
	Stat(string) (os.FileInfo, error)
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

func (FileSystemReal) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

type Handler struct {
	IBMServiceCreds    ibmservicecreds.IBMServiceCreds
	Speedfactors       speedfactors.Speedfactors
	Uuid               Uuid
	FileSystem         FileSystem
	TimesAndSizeToPath timesandsize.TimesAndSizeToPath
	Runner             diarizerrunner.DiarizerRunner
	SoundRootDir       string
	ResultsRootDir     string
	WaitForDiarizer    bool
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// `diarizer` and `transcriber` _must_ match what the Android client sends
	// e.g.  https://github.com/blabbertabber/blabbertabber/blob/fea98684ad500380ef347cff584821ee52098c1e/app/src/main/java/com/blabbertabber/blabbertabber/RecordingActivity.java#L386-L392
	diarizer := r.Header.Get("Diarizer")       // "IBM" or "Aalto"
	transcriber := r.Header.Get("Transcriber") // "IBM" or "CMUSphinx4"

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
		// the 2nd portion ("|| err.Error() == ...") is to make the tests work; it shouldn't be necessary
		if err != nil && (err == io.EOF || err.Error() == "multipart: NextPart: EOF") {
			break
		}
		if err != nil {
			panic(err.Error())
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

	wavFileInfo, err := h.FileSystem.Stat(filepath.Join(soundDir, "meeting.wav"))
	if err != nil {
		panic(err.Error())
	}
	tas := timesandsize.TimesAndSize{
		WaveFileSizeInBytes:              wavFileInfo.Size(),
		Diarizer:                         diarizer,
		Transcriber:                      transcriber,
		EstimatedDiarizationFinishTime:   timesandsize.JSONTime(time.Now().Add(h.Speedfactors.EstimatedDiarizationTime(diarizer, wavFileInfo.Size()))),
		EstimatedTranscriptionFinishTime: timesandsize.JSONTime(time.Now().Add(h.Speedfactors.EstimatedTranscriptionTime(transcriber, wavFileInfo.Size()))),
		DiarizationProcessingRatio:       0,
		TranscriptionProcessingRatio:     0,
	}
	h.TimesAndSizeToPath.WriteTimesAndSizeToPath(&tas, filepath.Join(resultsDir, "times_and_size.json"))

	// sync.WaitGroup accommodates our testing requirements
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		h.Runner.Run(diarizer, meetingUuid, h.IBMServiceCreds)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		if diarizer != transcriber {
			h.Runner.Run(transcriber, meetingUuid, h.IBMServiceCreds)
		}
		wg.Done()
	}()
	if h.WaitForDiarizer {
		wg.Wait()
	} // tests need to wait but production should return immediately
}
