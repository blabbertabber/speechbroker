package speedfactors_test

import (
	. "github.com/blabbertabber/speechbroker/speedfactors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"strings"
	"time"
)

var _ = Describe("Speedfactors", func() {
	Context(".ReadCredsFromReader", func() {
		Context("When a reader is passed an empty JSON", func() {
			It("returns an empty struct", func() {
				sourceReader := strings.NewReader("{}")
				expectation := Speedfactors{}
				readCreds, err := ReadCredsFromReader(sourceReader)
				Expect(err).To(BeNil())
				Expect(readCreds).To(Equal(expectation))
			})
		})
		Context("When a reader is passed a populated JSON", func() {
			It("returns a populated struct", func() {
				source, err := ioutil.ReadFile("../assets/test/speedfactors.json")
				Expect(err).To(BeNil())
				sourceReader := strings.NewReader(string(source))
				expectation := Speedfactors{
					Diarizer: map[string]float64{
						"IBM":   2.4,
						"Aalto": 0.5,
					},
					Transcriber: map[string]float64{
						"IBM":        2.4,
						"CMUSphinx4": 8.0,
					},
				}
				readCreds, err := ReadCredsFromReader(sourceReader)
				Expect(err).To(BeNil())
				Expect(readCreds).To(Equal(expectation))
			})
		})
	})
	Context(".ReadCredsFromPath", func() {
		Context("When path is exists and is valid JSON", func() {
			It("returns the expected struct", func() {
				readCreds, err := ReadCredsFromPath("../assets/test/speedfactors.json")
				Expect(err).To(BeNil())
				expectation := Speedfactors{
					Diarizer: map[string]float64{
						"IBM":   2.4,
						"Aalto": 0.5,
					},
					Transcriber: map[string]float64{
						"IBM":        2.4,
						"CMUSphinx4": 8.0,
					},
				}
				Expect(readCreds).To(Equal(expectation))
			})
		})
		Context("When path is non-existent", func() {
			It("returns an error", func() {
				_, err := ReadCredsFromPath("/non/existent/path")
				Expect(err.Error()).To(MatchRegexp("open /non/existent/path:"))
			})
		})
	})
	Context("EstimatedDiarizationTime", func() {
		It("calculates the expected diarization time based on file size", func() {
			sf := Speedfactors{
				Diarizer: map[string]float64{
					"Aalto": 0.5,
				},
			}
			// 32,000 bytes/second, 10-minute file is 19,200,000
			Expect(sf.EstimatedDiarizationTime("Aalto", 19200000)).To(Equal(time.Minute * 5))
		})
	})
	Context("EstimatedTranscriptionTime", func() {
		It("calculates the expected transcription time based on file size", func() {
			sf := Speedfactors{
				Transcriber: map[string]float64{
					"CMUSphinx4": 8.0,
				},
			}
			// 32,000 bytes/second, 10-minute file is 19,200,000
			Expect(sf.EstimatedTranscriptionTime("CMUSphinx4", 19200000)).To(Equal(time.Minute * 80))
		})
	})
})
