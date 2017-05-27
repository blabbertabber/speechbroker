package emitblabbertabber

import (
	"fmt"
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
	"log"
	"math"
	"strings"
)

type Summary struct {
	TotalSpeakingTime float64       `json:"total_speaking_time"`
	LeaderBoard       []SpeakerStat `json:"leader_board"`
	Utterances        []Utterance   `json:"utterances"`
}

type SpeakerStat struct {
	Speaker   string  `json:"speaker"`
	TotalTime float64 `json:"total_time"`
}

type Utterance struct {
	Speaker    int     `json:"speaker"`
	From       float64 `json:"from"`
	To         float64 `json:"to"`
	Transcript string  `json:"transcript"`
}

func CalcSummary(utterances []Utterance) (summary Summary, err error) {
	//summary = Summary{}
	var totalSpeakingTime float64 = 0.0
	speakerStatMap := make(map[int]float64)
	leaderBoard := []SpeakerStat{}
	for _, utterance := range utterances {
		elapsedTime := utterance.To - utterance.From
		elapsedTime = math.Floor(elapsedTime*100+.5) / 100
		totalSpeakingTime += elapsedTime
		speakerStatMap[utterance.Speaker] += elapsedTime
	}
	for key, value := range speakerStatMap {
		leaderBoard = append(leaderBoard, SpeakerStat{
			Speaker:   fmt.Sprintf("%d", key),
			TotalTime: value,
		})
	}
	summary.LeaderBoard = leaderBoard
	summary.TotalSpeakingTime = totalSpeakingTime
	summary.Utterances = utterances
	return summary, nil
}

func Coerce(transaction parseibm.IBMTranscription) (utterances []Utterance, err error) {
	transaction.SpeakerLabels, err = SquashSpeakerLabels(transaction.SpeakerLabels)
	if err != nil {
		log.Fatal("I was unable to squash the speakers!")
	}

	if transaction.Results == nil {
		return []Utterance{}, nil
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
