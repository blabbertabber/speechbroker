package emitblabbertabber_test

import (
	. "github.com/blabbertabber/DiarizerServer/IBMJson/emitblabbertabber"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
	"io/ioutil"
	"encoding/json"
)

var _ = Describe("Emitblabbertabber", func() {
	It("should emit an empty struct to an empty JSON properly", func() {
		emptyTrans := parseibm.IBMTranscription{}
		out, err := Coerce(emptyTrans)
		Expect(err).To(BeNil())
		Expect(out).To(Equal([]byte(`{}`)))
	})
	It("should output a correctly transformed JSON", func() {
		source, err := ioutil.ReadFile("../../assets/test/ibm_1.json")
		Expect(err).To(BeNil())
		trans := parseibm.IBMTranscription{}
		err = json.Unmarshal(source, &trans)
		out, err := Coerce(trans)
		Expect(err).To(BeNil())
		expectation := Transcriptions{
			Utterance{
				Speaker:    0,
				From:       2.37,
				To:         17.54,
				Transcript: "design swift transaction so you go through when you put all all the things you need to do",
			},
		}
		expectedJson, err := json.Marshal(expectation)
		Expect(out).To(Equal(expectedJson))
	})

	It("should output a correctly transformed JSON", func() {
		source, err := ioutil.ReadFile("../../assets/test/ibm_2.json")
		Expect(err).To(BeNil())
		trans := parseibm.IBMTranscription{}
		err = json.Unmarshal(source, &trans)
		out, err := Coerce(trans)
		Expect(err).To(BeNil())
		expectation := Transcriptions{
			Utterance{
				Speaker:    0,
				From:       2.37,
				To:         9.55,
				Transcript: "design swift transaction sure",
			},
		}
		expectedJson, err := json.Marshal(expectation)
		Expect(out).To(Equal(expectedJson))
	})

})
