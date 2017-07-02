// IBMJson < xxx.json > yyy.json
package main

import (
	"encoding/json"
	"fmt"
	"github.com/blabbertabber/speechbroker/IBMJson/emitblabbertabber"
	"github.com/blabbertabber/speechbroker/IBMJson/parseibm"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// source := []byte(`"result_index": 0, "results": [], "speaker_labels": []`)
	source, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err.Error())
	}
	var input parseibm.IBMTranscription // the whole, complete transcription
	err = json.Unmarshal(source, &input)
	if err != nil {
		panic(err.Error())
	}
	transcriptions, err := emitblabbertabber.Coerce(input)
	if err != nil {
		panic(err.Error())
	}
	bytes, err := json.Marshal(transcriptions)

	fmt.Println(string(bytes))
}
