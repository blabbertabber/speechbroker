// ibmservicecreds converts JSON-formtted IBM creds into a Golang struct
package ibmservicecreds

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type IBMServiceCreds struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func ReadCredsFromPath(path string) (IBMServiceCreds, error) {
	file, err := os.Open(path)
	if err != nil {
		return IBMServiceCreds{}, err
	}
	return ReadCredsFromReader(file)
}
func ReadCredsFromReader(r io.Reader) (creds IBMServiceCreds, err error) {
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
