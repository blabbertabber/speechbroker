package emitblabbertabber

import (
	"encoding/json"
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
	"log"
	"math"
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
	transaction.SpeakerLabels, err = SquashSpeakerLabels(transaction.SpeakerLabels)
	log.Fatal("I was unable to squash the speakers!")

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

func SquashSpeakerLabels(speakerLabels []parseibm.SpeakerLabel) ([]parseibm.SpeakerLabel, error) {
	squashed := []parseibm.SpeakerLabel{}
	newSpeaker := math.MaxInt64
	for _, speakerLabel := range speakerLabels {
		if newSpeaker != speakerLabel.Speaker {
			newSpeaker = speakerLabel.Speaker
			squashed = append(squashed, parseibm.SpeakerLabel{
				Confidence: speakerLabel.Confidence,
				Final:      speakerLabel.Final,
				From:       speakerLabel.From,
				Speaker:    speakerLabel.Speaker,
				To:         speakerLabel.To,
			})
		} else {
			squashed[len(squashed)-1].To = speakerLabel.To
		}
	}
	return squashed, nil
}
