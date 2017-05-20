package emitblabbertabber_test

import (
	. "github.com/blabbertabber/DiarizerServer/IBMJson/emitblabbertabber"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
)

var _ = Describe("Emitblabbertabber", func() {
	It("should emit an empty struct to an empty JSON properly", func() {
		emptyTrans := parseibm.IBMTranscription{}
		out, err := Coerce(emptyTrans)
		Expect(err).To(BeNil())
		Expect(out).To(Equal([]byte(`{}`)))
	})

})
