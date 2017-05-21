package emitblabbertabber

import (
	"encoding/json"
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
	"strings"
)

type Transcriptions []Utterance

type Utterance struct {
	Speaker    int
	From       float64
	To         float64
	Transcript string
}

func Coerce(transaction parseibm.IBMTranscription) (bytes []byte, err error) {
	if transaction.Results == nil {
		return []byte(`{}`), nil
	}
	speaker := transaction.SpeakerLabels[0].Speaker
	from := transaction.SpeakerLabels[0].From
	to := transaction.SpeakerLabels[0].To
	var transcription []string
	results := transaction.Results
	for _, result := range results {
		timestamps := result.Alternatives[0].Timestamps
		for _, timestamp := range timestamps {
			if timestamp.Word != "%HESITATION" {
				transcription = append(transcription, timestamp.Word)
			}
		}
	}
	transcript := strings.Join(transcription, " ")

	utterance := Utterance{speaker, from, to, transcript}

	transcriptions := Transcriptions{utterance}

	bytes, err = json.Marshal(transcriptions)

	return bytes, err
}
