// ibmservicecreds converts JSON-formtted IBM creds into a Golang struct
package timesandsize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type TimesAndSizeToPath interface {
	WriteTimesAndSize(io.Writer)
}

type TimesAndSize struct {
	WaveFileSizeInBytes              int64     `json:"wav_file_size_in_bytes"`
	Diarizer                         string    `json:"diarizer"`
	Transcriber                      string    `json:"transcriber"`
	DiarizationProcessingRatio       float64   `json:"diarization_processing_ratio"`
	TranscriptionProcessingRatio     float64   `json:"transcription_processing_ratio"`
	EstimatedDiarizationFinishTime   time.Time `json:"estimated_diarization_finish_time"`
	EstimatedTranscriptionFinishTime time.Time `json:"estimated_transcription_finish_time"`
}

// for JSON marshalling & unmarshalling with custom time.Time format
type Doppelgänger TimesAndSize

var timeFormat = "2006-01-02T15:04:05-0700"

func (tas *TimesAndSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		*Doppelgänger
		EstimatedDiarizationFinishTime   string `json:"estimated_diarization_finish_time"`
		EstimatedTranscriptionFinishTime string `json:"estimated_transcription_finish_time"`
	}{
		Doppelgänger:                     (*Doppelgänger)(tas),
		EstimatedDiarizationFinishTime:   fmt.Sprintf("%s", time.Time(tas.EstimatedDiarizationFinishTime).Format(timeFormat)),
		EstimatedTranscriptionFinishTime: fmt.Sprintf("%s", time.Time(tas.EstimatedTranscriptionFinishTime).Format(timeFormat)),
	})
}

// the following is a convenience method for tests; is not
// used by production code (we only marshal, never unmarshal)
func (tas *TimesAndSize) UnmarshalJSON(data []byte) error {
	// strip out the quotes otherwise time.Parse will fail
	doppelganger := &struct {
		EstimatedDiarizationFinishTime   string `json:"estimated_diarization_finish_time"`
		EstimatedTranscriptionFinishTime string `json:"estimated_transcription_finish_time"`
		*Doppelgänger
	}{
		Doppelgänger: (*Doppelgänger)(tas),
	}
	if err := json.Unmarshal(data, &doppelganger); err != nil {
		panic(err)
	}
	dt, err := time.Parse(timeFormat, doppelganger.EstimatedDiarizationFinishTime)
	if err != nil {
		panic(err)
	}
	tt, err := time.Parse(timeFormat, doppelganger.EstimatedTranscriptionFinishTime)
	if err != nil {
		panic(err)
	}
	tas.EstimatedDiarizationFinishTime = dt
	tas.EstimatedTranscriptionFinishTime = tt
	return nil
}

func (tas *TimesAndSize) WriteTimesAndSize(w io.Writer) {
	b, err := json.Marshal(tas)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(b)
	_, err = buf.WriteTo(w)
	if err != nil {
		panic(err)
	}
}
