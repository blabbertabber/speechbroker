package emitblabbertabber

import (
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
)

func Coerce(transaction parseibm.IBMTranscription) (bytes []byte, err error) {
	bytes=[]byte(`{}`)
	return bytes, nil
}