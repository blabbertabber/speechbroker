// ibmservicecreds converts JSON-formtted IBM creds into a Golang struct
package timesandsize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// `counterfeiter timesandsize/timesandsize.go TimesAndSizeToPath`
type TimesAndSizeToPath interface {
	WriteTimesAndSizeToPath(*TimesAndSize, string)
}

type TimesAndSize struct {
	WaveFileSizeInBytes              int64    `json:"wav_file_size_in_bytes"`
	Diarizer                         string   `json:"diarizer"`
	Transcriber                      string   `json:"transcriber"`
	EstimatedDiarizationFinishTime   JSONTime `json:"estimated_diarization_finish_time"`
	EstimatedTranscriptionFinishTime JSONTime `json:"estimated_transcription_finish_time"`
	DiarizationProcessingRatio       float64  `json:"diarization_processing_ratio"`
	TranscriptionProcessingRatio     float64  `json:"transcription_processing_ratio"`
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05-0700"))
	return []byte(stamp), nil
}

func (tas *TimesAndSize) WriteTimesAndSizeToWriter(w io.Writer) {
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

func WriteTimesAndSizeToPath(tas *TimesAndSize, tasPath string) {
	file, err := os.Create(tasPath)
	if err != nil {
		panic(err)
	}
	tas.WriteTimesAndSizeToWriter(file)
	if err = file.Close(); err != nil {
		panic(err)
	}
}
