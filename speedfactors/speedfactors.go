// ibmservicecreds converts JSON-formtted IBM creds into a Golang struct
package speedfactors

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type Speedfactors struct {
	Diarizer    map[string]float64 `json:"diarizer"`
	Transcriber map[string]float64 `json:"transcriber"`
}

func ReadCredsFromPath(path string) (Speedfactors, error) {
	file, err := os.Open(path)
	if err != nil {
		return Speedfactors{}, err
	}
	return ReadCredsFromReader(file)
}
func ReadCredsFromReader(r io.Reader) (creds Speedfactors, err error) {
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
