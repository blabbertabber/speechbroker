// ibmservicecreds converts JSON-formtted IBM creds into a Golang struct
package speedfactors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type Speedfactors struct {
	Diarizer    map[string]float64 `json:"diarizer"`
	Transcriber map[string]float64 `json:"transcriber"`
}

func ReadCredsFromPath(path string) (Speedfactors, error) {
	file, err := os.Open(path)
	if err != nil {
		return Speedfactors{}, err
	}
	return ReadCredsFromReader(file)
}
func ReadCredsFromReader(r io.Reader) (creds Speedfactors, err error) {
	buf := []byte{}
	buf, err = ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(buf, &creds)
	if err != nil {
		panic(err)
	}
	return creds, err
}

// the following functions return time.Duration whose underlying type is int64 (nanosecs),
func (sf Speedfactors) EstimatedDiarizationTime(diarizer string, soundFileSizeinBytes int64) (time.Duration, error) {
	if val, ok := sf.Diarizer[diarizer]; ok {
		return meetingLength(int64(val * float64(soundFileSizeinBytes))), nil
	} else {
		return time.Duration(0), errors.New(fmt.Sprintf("I couldn't find Diarizer[\"%s\"]!", diarizer))
	}
}

func (sf Speedfactors) EstimatedTranscriptionTime(transcriber string, soundFileSizeinBytes int64) (time.Duration, error) {
	if val, ok := sf.Transcriber[transcriber]; ok {
		return meetingLength(int64(val * float64(soundFileSizeinBytes))), nil
	} else {
		return time.Duration(0), errors.New(fmt.Sprintf("I couldn't find Transcriber[\"%s\"]!", transcriber))
	}
}

func ProcessingRatio(start time.Time, finish time.Time, meetingWaveFileSize int64) float64 {
	return float64(finish.Sub(start)) / float64((meetingLength(meetingWaveFileSize)))
}

func meetingLength(soundFileSizeinBytes int64) time.Duration {
	return time.Duration(soundFileSizeinBytes * 1e9 / 32000)
}
