// IBMJson < xxx.json > yyy.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/blabbertabber/speechbroker/ibmjson/emitblabbertabber"
	"github.com/blabbertabber/speechbroker/ibmjson/parseibm"
	"io/ioutil"
	"os"
)

func main() {
	// source := []byte(`"result_index": 0, "results": [], "speaker_labels": []`)
	var inFile = flag.String("in", "",
		"path to JSON input file created by IBM's Speech-to-Text, e.g. '/blabbertabber/diarizationResults/0.json.txt'")
	var outFile = flag.String("out", "",
		"path to output file, e.g. 'summary.json'")
	flag.Parse()

	var source []byte
	var err error

	switch *inFile {
	case "":
		source, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err.Error())
		}
	default:
		inHandle, err := os.Open(*inFile)
		if err != nil {
			panic(err.Error())
		}
		source, err = ioutil.ReadAll(inHandle)
		if err != nil {
			panic(err.Error())
		}
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
	if err != nil {
		panic(err.Error())
	}

	switch *outFile {
	case "":
		fmt.Println(string(bytes))
	default:
		err := ioutil.WriteFile(*outFile, bytes, 0644)
		if err != nil {
			panic(err.Error())
		}
	}
}
