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
				const longForm = "Jan 2, 2006 at 3:04pm (MST)"
				t, _ := time.Parse(longForm, "Sep 4, 2017 at 9:49am (PST)")
				tas := TimesAndSize{
					Diarizer:                         "blah",
					Transcriber:                      "blech",
					WaveFileSizeInByte:               256,
					EstimatedDiarizationFinishTime:   t,
					EstimatedTranscriptionFinishTime: t,
				}
				bytewriter := bytes.NewBuffer([]byte{})
				tas.WriteTimesAndSizeToWriter(bytewriter)
				jsonwritten := bytewriter.String()

				expectation := `{"wav_file_size_in_bytes":256,"diarizer":"blah","transcriber":"blech",` +
					`"estimated_transcription_finish_time":"2017-09-04T09:49:00Z",` +
					`"estimated_diarization_finish_time":"2017-09-04T09:49:00Z"}`
				Expect(jsonwritten).To(Equal(expectation))
			})
		})
	})
})
