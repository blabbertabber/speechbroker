// IBMJson < xxx.json > yyy.json
package main

import (
	"encoding/json"
	"fmt"
	"github.com/blabbertabber/DiarizerServer/IBMJson/emitblabbertabber"
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
	"io/ioutil"
	"log"
	"os"
)

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
	btJSON, err := emitblabbertabber.Coerce(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(btJSON))
}
