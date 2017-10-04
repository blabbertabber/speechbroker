package httphandler_test

import (
	. "github.com/blabbertabber/speechbroker/httphandler"

	"bytes"
	"encoding/json"
	"errors"
	"github.com/blabbertabber/speechbroker/diarizerrunner/diarizerrunnerfakes"
	"github.com/blabbertabber/speechbroker/httphandler/httphandlerfakes"
	"github.com/blabbertabber/speechbroker/ibmservicecreds"
	"github.com/blabbertabber/speechbroker/speedfactors"
	"github.com/blabbertabber/speechbroker/timesandsize"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type fakeWriter struct {
	Written string
}

func (fw *fakeWriter) Write(p []byte) (n int, err error) {
	fw.Written += string(p)
	return len(p), nil
}
func (fw *fakeWriter) Header() http.Header {
	return http.Header{}
}
func (fw *fakeWriter) WriteHeader(i int) {
}

type fakeReadCloser string

func (frc *fakeReadCloser) Read(p []byte) (n int, err error) {
	if *frc == "" {
		return 0, io.EOF
	}
	p = []byte(*frc)
	*frc = ""
	return len(p), nil
}

func (frc *fakeReadCloser) Close() error {
	return nil
}

var _ = Describe("Httphandler", func() {

	var miscWriter [2]*bytes.Buffer
	var handler Handler
	var r *http.Request
	var fw *fakeWriter
	var ffs *httphandlerfakes.FakeFileSystem
	var fdr *diarizerrunnerfakes.FakeDiarizerRunner
	var ffi *httphandlerfakes.FakeFileInfo
	var boundary = "ILoveMyDogCherieSheIsSoWarmAndCuddly"

	BeforeEach(func() {
		fakeUuid := new(httphandlerfakes.FakeUuid)
		fakeUuid.GetUuidReturns("fake-uuid")
		ffs = new(httphandlerfakes.FakeFileSystem)
		ffs.MkdirAllReturns(nil)
		// create distinct file handles for each os.Create()
		for i := 0; i < 3; i++ {
			fh, err := os.Create(os.DevNull)
			if err != nil {
				panic(err)
			}
			ffs.CreateReturnsOnCall(i, fh, nil)
		}
		for i := 0; i < 2; i++ {
			miscWriter[i] = bytes.NewBuffer([]byte{})
			fh, err := os.Create(os.DevNull)
			if err != nil {
				panic(err)
			}
			ffs.CreateWriterReturnsOnCall(i, fh, miscWriter[i], nil)
		}
		fdr = new(diarizerrunnerfakes.FakeDiarizerRunner)
		ffi = new(httphandlerfakes.FakeFileInfo)
		ffi.SizeReturns(38400000) // 20 minutes = 32000 bytes/sec * 60 * 20
		ffs.StatReturns(ffi, nil)

		fakeUuid.GetUuidReturns("fake-uuid")

		handler = Handler{
			IBMServiceCreds: ibmservicecreds.IBMServiceCreds{},
			Speedfactors: speedfactors.Speedfactors{
				Diarizer: map[string]float64{
					"Aalto": 0.6,
					"IBM":   2.4,
				},
				Transcriber: map[string]float64{
					"CMUSphinx4": 8.0,
					"IBM":        2.4,
				},
			},
			Uuid:            fakeUuid,
			FileSystem:      ffs,
			Runner:          fdr,
			SoundRootDir:    "/a/b",
			ResultsRootDir:  "/c/d",
			WaitForDiarizer: true,
		}

		frc := fakeReadCloser("--" + boundary + "\r\n" +
			"Content-Disposition: form-data; name=\"soundFile; filename=\"meeting.wav\"\r\n" +
			"Content-Type: application/octet-stream\r\n\r\n" +
			"contents_of_dot_wav_file" +
			"\r\n--" + boundary + "--\r\n")

		r = &http.Request{
			Method: "POST",
			Host:   "test.diarizer.com:8080",
			URL:    &url.URL{Path: "/api/v1/upload"},
			Header: http.Header{
				"Diarizer":          {"Aalto"},
				"Transcriber":       {"CMUSphinx4"},
				"User-Agent":        {"Mozilla/5.0 ( compatible )"},
				"Accept":            {"*/*"},
				"Connection":        {"Keep-Alive"},
				"Content-Type":      {"multipart/form-data; boundary=" + boundary},
				"Accept-Encoding":   {"identity"},
				"Transfer-Encoding": {"chunked"},
				"Host":              {"test.diarizer.com:8080"},
			},
			Proto:         "HTTP/1.1",
			ContentLength: -1,
			ProtoMajor:    1,
			ProtoMinor:    1,
			RemoteAddr:    "test.diarizer.com:8080",
			RequestURI:    "/api/v1/upload",
			Body:          &frc,
		}

		fw = &fakeWriter{}
	})
	Describe("ServeHTTP", func() {
		Context("when it's unable to create the first directory", func() {
			It("should panic", func() {
				errorTxt := "Can't create dir"
				ffs.MkdirAllReturns(errors.New(errorTxt))
				handler.ServeHTTP(fw, r)
				Expect(fw.Written).To(MatchRegexp(errorTxt))
			})
		})
		Context("when it's unable to create the second directory", func() {
			It("should panic", func() {
				errorTxt := "Can't create dir"
				ffs.MkdirAllReturnsOnCall(0, nil)
				ffs.MkdirAllReturnsOnCall(1, errors.New(errorTxt))
				handler.ServeHTTP(fw, r)
				Expect(fw.Written).To(MatchRegexp(errorTxt))
			})
		})
		Context("when it's unable to create a file", func() {
			It("should panic", func() {
				errorTxt := "Can't create file"
				ffs.CreateReturnsOnCall(0, nil, errors.New(errorTxt))
				handler.ServeHTTP(fw, r)
				Expect(fw.Written).To(MatchRegexp(errorTxt))
			})
		})
		Context("when it's unable to create the second file", func() {
			It("should panic", func() {
				fakefile, err := os.Create(os.DevNull)
				if err != nil {
					panic(err.Error())
				}
				errorTxt := "Can't create file"
				ffs.CreateReturnsOnCall(0, fakefile, nil)
				ffs.CreateReturnsOnCall(1, nil, errors.New(errorTxt))
				handler.ServeHTTP(fw, r)
				Expect(fw.Written).To(MatchRegexp(errorTxt))
			})
		})
		Context("when it's writing the 'times_and_sizes' JSON file", func() {
			It("should call WriteTimesAndSize()", func() {
				handler.ServeHTTP(fw, r)
				Expect(ffs.CreateWriterCallCount()).To(Equal(2))
				Expect(ffs.CreateWriterArgsForCall(0)).To(Equal(filepath.FromSlash("/c/d/fake-uuid/times_and_size.json")))
				regex := `{` +
					`"wav_file_size_in_bytes":38400000,"diarizer":"Aalto","transcriber":"CMUSphinx4",` +
					`"diarization_processing_ratio":0,` +
					`"transcription_processing_ratio":0,` +
					`"estimated_diarization_finish_time":".*",` +
					`"estimated_transcription_finish_time":".*"` +
					`}`
				Expect(miscWriter[0].String()).To(MatchRegexp(regex))
			})
		})
		Context("when using Aalto + CMU Sphinx", func() {
			It("send the correct value to the client", func() {
				handler.ServeHTTP(fw, r)
				Expect(fw.Written).To(Equal("https://test.diarizer.com?meeting=fake-uuid"))
			})
			It("invokes the diarizer runner with the correct arguments", func() {

				handler.ServeHTTP(fw, r)
				Expect(fdr.RunCallCount()).To(Equal(2))
				backEnd0, _, _ := fdr.RunArgsForCall(0)
				backEnd1, _, _ := fdr.RunArgsForCall(1)
				Expect(backEnd0).To(Not(Equal(backEnd1)))
				for i := 0; i < fdr.RunCallCount(); i++ {
					action, uuid, _ := fdr.RunArgsForCall(i)
					Expect(uuid).To(Equal("fake-uuid"))
					switch action {
					case "Aalto":
					case "CMUSphinx4":
					default:
						panic("I have no idea what action this should be: " + action)
					}
				}
				tas := timesandsize.TimesAndSize{}
				if err := json.Unmarshal(miscWriter[0].Bytes(), &tas); err != nil {
					panic(err)
				}
				Expect(tas.WaveFileSizeInBytes).To(Equal(int64(38400000)))
				// calculate delta by hand: 1200s meeting takes 720s to diarize
				// we subtract 1/2 second before rounding to make it reliably pass
				Expect(time.Time(tas.EstimatedDiarizationFinishTime)).
					To(Equal(time.Now().Add(time.Second * 720).Add(time.Millisecond * -500).
						Round(time.Second)))
				// Same expectation, but instead of calculating by hand we run
				// the calculations through the functions (maybe this is overkill).
				estElapDiarTime, _ := handler.Speedfactors.EstimatedDiarizationTime("Aalto", tas.WaveFileSizeInBytes)
				estElapTransTime, _ := handler.Speedfactors.EstimatedTranscriptionTime("CMUSphinx4", tas.WaveFileSizeInBytes)

				Expect(time.Time(tas.EstimatedDiarizationFinishTime)).
					To(Equal(time.Now().
						Add(estElapDiarTime).
						Add(time.Millisecond * -500).
						Round(time.Second)))
				// calculate by hand
				Expect(time.Time(tas.EstimatedTranscriptionFinishTime)).
					To(Equal(time.Now().Add(time.Second * 9600).Add(time.Millisecond * -500).
						Round(time.Second)))
				// calculate by function
				Expect(time.Time(tas.EstimatedTranscriptionFinishTime)).
					To(Equal(time.Now().
						Add(estElapTransTime).
						Add(time.Millisecond * -500).
						Round(time.Second)))
			})
		})
		Context("when using IBM for both transcription and Diarization", func() {
			It("invokes Docker but once", func() {
				r.Header.Del("Diarizer")
				r.Header.Add("Diarizer", "IBM")
				r.Header.Del("Transcriber")
				r.Header.Add("Transcriber", "IBM")
				handler.ServeHTTP(fw, r)
				Expect(fdr.RunCallCount()).To(Equal(1))
				backEnd, _, _ := fdr.RunArgsForCall(0)
				Expect(backEnd).To(Equal("IBM"))
			})
		})
	})
})
