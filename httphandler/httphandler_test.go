package httphandler_test

import (
	. "github.com/blabbertabber/speechbroker/httphandler"

	"errors"
	"github.com/blabbertabber/speechbroker/diarizerrunner/diarizerrunnerfakes"
	"github.com/blabbertabber/speechbroker/httphandler/httphandlerfakes"
	"github.com/blabbertabber/speechbroker/ibmservicecreds"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/url"
	"os"
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

	var handler Handler
	var r *http.Request
	var fw *fakeWriter
	var ffs *httphandlerfakes.FakeFileSystem
	var fdr *diarizerrunnerfakes.FakeDiarizerRunner
	//var frc *fakeReadCloser
	var boundary = "ILoveMyDogCherieSheIsSoWarmAndCuddly"

	BeforeEach(func() {
		fakeUuid := new(httphandlerfakes.FakeUuid)
		fakeUuid.GetUuidReturns("fake-uuid")
		ffs = new(httphandlerfakes.FakeFileSystem)
		ffs.MkdirAllReturns(nil)
		ffs.CreateReturns(os.Create("/dev/null"))
		fdr = new(diarizerrunnerfakes.FakeDiarizerRunner)

		fakeUuid.GetUuidReturns("fake-uuid")

		handler = Handler{
			Uuid:            fakeUuid,
			IBMServiceCreds: ibmservicecreds.IBMServiceCreds{},
			FileSystem:      ffs,
			Runner:          fdr,
			SoundRootDir:    "/a/b",
			ResultsRootDir:  "/c/d",
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
				ffs.MkdirAllReturns(errors.New("Can't create dir"))
				Expect(func() { handler.ServeHTTP(fw, r) }).To(Panic())
			})
		})
		Context("when it's unable to create the second directory", func() {
			It("should panic", func() {
				ffs.MkdirAllReturnsOnCall(0, nil)
				ffs.MkdirAllReturnsOnCall(1, errors.New("Can't create dir"))
				Expect(func() { handler.ServeHTTP(fw, r) }).To(Panic())
			})
		})
		Context("when it's unable to create a file", func() {
			It("should panic", func() {
				ffs.CreateReturns(nil, errors.New("create file"))
				Expect(func() { handler.ServeHTTP(fw, r) }).To(Panic())
			})
		})
		Context("when it's unable to create the second file", func() {
			It("should panic", func() {
				fakefile, err := os.Create("/dev/null")
				if err != nil {
					panic(err.Error())
				}
				ffs.CreateReturnsOnCall(0, fakefile, nil)
				ffs.CreateReturnsOnCall(1, nil, errors.New("create file"))
				Expect(func() { handler.ServeHTTP(fw, r) }).To(Panic())
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
			})
		})
		Context("when using IBM for both transcription and Diarization", func() {
			//r.Header.Set("Transcriber", "IBM")
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
