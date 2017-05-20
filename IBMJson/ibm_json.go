// IBMJson < xxx.json > yyy.json
package main

import (
	"encoding/json"
	"fmt"
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
	"io/ioutil"
	"log"
	"os"
)

/*
    B L A B B E R T A B B E R   D A T A   S T R U C T U R E S

    typical JSON-marshaled format:

{
	"speaker_totals": {
		"0": 34.7,
		"1": 35
	},
	"transcription": [
		{
			"speaker": "0",
			"words": "I love my dog",
			"from": 2.37,
			"to": 5.4
		}
	]
}
*/

func main() {
	// source := []byte(`"result_index": 0, "results": [], "speaker_labels": []`)
	source, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var input parseibm.IBMTranscription // the whole, complete transcription
	err = json.Unmarshal(source, &input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(input)
}
