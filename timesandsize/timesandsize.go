// ibmservicecreds converts JSON-formtted IBM creds into a Golang struct
package timesandsize

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

type TimesAndSize struct {
	WaveFileSizeInByte               int       `json:"wav_file_size_in_bytes"`
	Diarizer                         string    `json:"diarizer"`
	Transcriber                      string    `json:"transcriber"`
	EstimatedTranscriptionFinishTime time.Time `json:"estimated_transcription_finish_time"`
	EstimatedDiarizationFinishTime   time.Time `json:"estimated_diarization_finish_time"`
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
