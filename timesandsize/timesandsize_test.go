package timesandsize_test

import (
	"bytes"
	. "github.com/blabbertabber/speechbroker/timesandsize"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
					WaveFileSizeInByte:               256,
					EstimatedDiarizationFinishTime:   JSONTime(t),
					EstimatedTranscriptionFinishTime: JSONTime(t.Add(time.Second)),
				}
				bytewriter := bytes.NewBuffer([]byte{})
				tas.WriteTimesAndSizeToWriter(bytewriter)
				jsonwritten := bytewriter.String()

				expectation := `{` +
					`"wav_file_size_in_bytes":256,"diarizer":"blah","transcriber":"blech",` +
					`"estimated_diarization_finish_time":"2017-09-04T20:49:14-0700",` +
					`"estimated_transcription_finish_time":"2017-09-04T20:49:15-0700"` +
					`}`
				Expect(jsonwritten).To(Equal(expectation))
			})
		})
	})
})