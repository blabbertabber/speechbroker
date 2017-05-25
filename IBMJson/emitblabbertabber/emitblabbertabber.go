package emitblabbertabber

import (
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

func Coerce(transaction parseibm.IBMTranscription) (utterances Transcriptions, err error) {
	transaction.SpeakerLabels, err = SquashSpeakerLabels(transaction.SpeakerLabels)
	if err != nil {
		log.Fatal("I was unable to squash the speakers!")
	}

	if transaction.Results == nil {
		return Transcriptions{}, nil
	}
	for _, speakerLabel := range transaction.SpeakerLabels {
		speaker := speakerLabel.Speaker
		from := speakerLabel.From
		to := speakerLabel.To
		var transcription []string
		results := transaction.Results
		for _, result := range results {
			timestamps := result.Alternatives[0].Timestamps
			for _, timestamp := range timestamps {
				if (timestamp.From >= speakerLabel.From) && (timestamp.To <= speakerLabel.To) {
					if timestamp.Word != "%HESITATION" {
						transcription = append(transcription, timestamp.Word)
					}
				}
			}
		}
		transcript := strings.Join(transcription, " ")

		utterances = append(utterances, Utterance{speaker, from, to, transcript})
	}

	return utterances, err
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
