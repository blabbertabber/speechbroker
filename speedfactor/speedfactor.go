// ibmservicecreds converts JSON-formtted IBM creds into a Golang struct
package speedfactor

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type Speedfactor struct {
	Diarizer    map[string]float64 `json:"diarizer"`
	Transcriber map[string]float64 `json:"transcriber"`
}

func ReadCredsFromPath(path string) (Speedfactor, error) {
	file, err := os.Open(path)
	if err != nil {
		return Speedfactor{}, err
	}
	return ReadCredsFromReader(file)
}
func ReadCredsFromReader(r io.Reader) (creds Speedfactor, err error) {
	buf := []byte{}
	buf, err = ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(buf, &creds)
	if err != nil {
		panic(err)
	}
	return creds, err
}
