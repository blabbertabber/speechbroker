// IBMJson < xxx.json > yyy.json
package main

import (
	"encoding/json"
	"fmt"
)

type IBMTranscription struct {
	ResultIndex   int       `json:"result_index"`
	Results       []Result  `json:"results"`
	SpeakerLabels []Speaker `json:"speaker_labels"`
}

type Result struct {
	Alternatives []Alternative `json:"alternatives"`
	Final        bool          `json:"final"`
}

type Alternative struct {
	Confidence float64     `json:"confidence"`
	Timestamps []Timestamp `json:"timestamps"`
	Transcript string      `json:"transcript"`
}

type Timestamp struct {
	entry []interface{}
}

type Speaker struct {
	Confidence float64 `json:"confidence"`
	Final      bool    `json:"final"`
	From       float64 `json:"from"`
	To         float64 `json:"to"`
}

func main() {
	source := []byte(`"result_index": 0, "results": [], "speaker_labels": []`)
	var input IBMTranscription
	json.Unmarshal(source, &input)
	fmt.Println(input.ResultIndex, input.Results)
}
