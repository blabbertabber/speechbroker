package timesandsize_test

import (
	"bytes"
	. "github.com/blabbertabber/speechbroker/timesandsize"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var _ = Describe("TimesAndSize", func() {
	Context(".WriteTimesAndSizeToWriter", func() {
		Context("When a writer is called on a valid TimesAndSize struct", func() {
			It("writes out the JSON struct to the writer", func() {
				const longForm = "Mon Jan 2 15:04:05 -0700 MST 2006"
				t, err := time.Parse(longForm, "Mon Sep 4 20:49:14 -0700 PDT 2017")
				if err != nil {
					panic(err)
				}
				tas := TimesAndSize{
					Diarizer:                         "blah",
					Transcriber:                      "blech",
					WaveFileSizeInBytes:              256,
					EstimatedDiarizationFinishTime:   JSONTime(t),
					EstimatedTranscriptionFinishTime: JSONTime(t.Add(time.Second)),
					DiarizationProcessingRatio:       0,
					TranscriptionProcessingRatio:     0,
				}
				bytewriter := bytes.NewBuffer([]byte{})
				tas.WriteTimesAndSizeToWriter(bytewriter)
				jsonwritten := bytewriter.String()

				expectation := `{` +
					`"wav_file_size_in_bytes":256,"diarizer":"blah","transcriber":"blech",` +
					`"estimated_diarization_finish_time":"2017-09-04T20:49:14-0700",` +
					`"estimated_transcription_finish_time":"2017-09-04T20:49:15-0700",` +
					`"diarization_processing_ratio":0,` +
					`"transcription_processing_ratio":0` +
					`}`
				Expect(jsonwritten).To(Equal(expectation))
			})
		})
	})
	Context(".WriteTimesAndSizeToPath", func() {
		Context("When the TimesAndPath is called on to write to a path", func() {
			It("writes out the JSON struct to the file", func() {
				const longForm = "Mon Jan 2 15:04:05 -0700 MST 2006"
				t, err := time.Parse(longForm, "Tue Sep 5 20:49:14 -0700 PDT 2017")
				if err != nil {
					panic(err)
				}
				tas := TimesAndSize{
					Diarizer:                         "boom",
					Transcriber:                      "bang",
					WaveFileSizeInBytes:              65536,
					EstimatedDiarizationFinishTime:   JSONTime(t),
					EstimatedTranscriptionFinishTime: JSONTime(t.Add(time.Minute)),
					DiarizationProcessingRatio:       0,
					TranscriptionProcessingRatio:     0,
				}

				tmpfile, err := ioutil.TempFile("", "times_and_sizes.json")
				if err != nil {
					log.Fatal(err)
				}

				defer os.Remove(tmpfile.Name()) // clean up

				WriteTimesAndSizeToPath(&tas, tmpfile.Name())

				jsonread, err := ioutil.ReadFile(tmpfile.Name())
				if err != nil {
					log.Fatal(err)
				}

				expectation := `{` +
					`"wav_file_size_in_bytes":65536,"diarizer":"boom","transcriber":"bang",` +
					`"estimated_diarization_finish_time":"2017-09-05T20:49:14-0700",` +
					`"estimated_transcription_finish_time":"2017-09-05T20:50:14-0700",` +
					`"diarization_processing_ratio":0,` +
					`"transcription_processing_ratio":0` +
					`}`
				Expect(string(jsonread)).To(Equal(expectation))
			})
		})
	})
})
