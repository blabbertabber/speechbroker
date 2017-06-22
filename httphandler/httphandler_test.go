package httphandler_test

import (
	. "github.com/blabbertabber/speechbroker/httphandler"

	"errors"
	"github.com/blabbertabber/speechbroker/httphandler/httphandlerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
	"net/http/httptest"
	"bytes"
)

type fakeWriter struct {
	Written string
}

func (fw fakeWriter) Write(p []byte) (n int, err error) {
	fw.Written += string(p)
	return len(p), nil
}
func (fw fakeWriter) Header() http.Header {
	return http.Header{}
}
func (fw fakeWriter) WriteHeader(i int) {
}

var _ = Describe("Httphandler", func() {

	var handler Handler
	var r *http.Request
	var fw fakeWriter
	var ffs *httphandlerfakes.FakeFileSystem

	BeforeEach(func() {
		fakeUuid := new(httphandlerfakes.FakeUuid)
		fakeUuid.GetUuidReturns("xyz")
		ffs = new(httphandlerfakes.FakeFileSystem)
		ffs.MkdirAllReturns( nil)
		ffs.CreateReturns(os.Create("/dev/null"))
		fakeDockerRunner := new(httphandlerfakes.FakeDockerRunner)

		fakeUuid.GetUuidReturns("fake-uuid")

		handler = Handler{
			Uuid:           fakeUuid,
			FileSystem:     ffs,
			SoundRootDir:   "/a/b",
			ResultsRootDir: "/c/d",
			DockerRunner:   fakeDockerRunner,
		}

		r = httptest.NewRequest(http.MethodPost,"https:/api/v1/upload", bytes.NewBufferString("hallo!"))
		//r = &http.Request{
		//	Header: http.Header{
		//		"Diarizer": {
		//			"IBM",
		//		},
		//		"Transcriber": {
		//			"IBM",
		//		},
		//	},
		//}

		fw = fakeWriter{}
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
		Context("when everything works", func() {
			It("should NOT panic", func() {
				handler.ServeHTTP(fw, r)
				Expect(func() { handler.ServeHTTP(fw, r) }).ToNot(Panic())
			})
		})
	})

})
