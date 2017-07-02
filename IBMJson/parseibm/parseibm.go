package parseibm

import (
	"encoding/json"
)

/*
   I B M   D A T A   S T R U C T U R E S
*/

type IBMTranscription struct {
	ResultIndex   int            `json:"result_index"`
	Results       []Result       `json:"results"`
	SpeakerLabels []SpeakerLabel `json:"speaker_labels"`
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
	Word string
	From float64
	To   float64
}

type SpeakerLabel struct {
	Confidence float64 `json:"confidence"`
	Final      bool    `json:"final"`
	From       float64 `json:"from"`
	Speaker    int     `json:"speaker"`
	To         float64 `json:"to"`
}

// http://attilaolah.eu/2013/11/29/json-decoding-in-go/
func (ts *Timestamp) UnmarshalJSON(b []byte) (err error) {
	j := []interface{}{"", 0, 0}
	if err = json.Unmarshal(b, &j); err == nil {
		ts.Word = j[0].(string)
		ts.From = j[1].(float64)
		ts.To = j[2].(float64)
		return
	}
	return
}
