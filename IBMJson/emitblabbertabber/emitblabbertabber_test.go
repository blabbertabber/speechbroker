package emitblabbertabber_test

import (
	. "github.com/blabbertabber/DiarizerServer/IBMJson/emitblabbertabber"

	"encoding/json"
	"github.com/blabbertabber/DiarizerServer/IBMJson/parseibm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Emitblabbertabber", func() {
	Context("#Coerce", func() {
		Context("when the IBMTranscription is an empty object", func() {
			It("should emit an empty struct to an empty JSON properly", func() {
				emptyTrans := parseibm.IBMTranscription{}
				out, err := Coerce(emptyTrans)
				Expect(err).To(BeNil())
				Expect(out).To(Equal([]Utterance{}))
			})
		})
		Context("when the IBMTranscription has one result (with one timestamp) and one speaker_label", func() {
			It("should output a correctly transformed JSON", func() {
				source, err := ioutil.ReadFile("../../assets/test/ibm_1.json")
				Expect(err).To(BeNil())
				trans := parseibm.IBMTranscription{}
				err = json.Unmarshal(source, &trans)
				out, err := Coerce(trans)
				Expect(err).To(BeNil())
				expectation := []Utterance{
					Utterance{
						Speaker:    0,
						From:       2.37,
						To:         17.54,
						Transcript: "design swift transaction so you go through when you put all all the things you need to do",
					},
				}
				Expect(out).To(Equal(expectation))
			})
		})

		Context("when the IBMTranscription has one result (with multiple timestamps) and one speaker_label", func() {
			It("should aggregate utterances and strip \"%HESITATION\"", func() {
				source, err := ioutil.ReadFile("../../assets/test/ibm_2.json")
				Expect(err).To(BeNil())
				trans := parseibm.IBMTranscription{}
				err = json.Unmarshal(source, &trans)
				out, err := Coerce(trans)
				Expect(err).To(BeNil())
				expectation := []Utterance{
					Utterance{
						Speaker:    0,
						From:       2.37,
						To:         9.55,
						Transcript: "design swift transaction sure",
					},
				}
				Expect(out).To(Equal(expectation))
			})
		})

		Context("when the IBMTranscription has one result (with multiple timestamps) and one speaker_label", func() {
			It("should aggregate utterances and strip \"%HESITATION\"", func() {
				source, err := ioutil.ReadFile("../../assets/test/ibm_2.json")
				Expect(err).To(BeNil())
				trans := parseibm.IBMTranscription{}
				err = json.Unmarshal(source, &trans)
				out, err := Coerce(trans)
				Expect(err).To(BeNil())
				expectation := []Utterance{
					Utterance{
						Speaker:    0,
						From:       2.37,
						To:         9.55,
						Transcript: "design swift transaction sure",
					},
				}
				Expect(out).To(Equal(expectation))
			})
		})

		Context("when the IBMTranscription has one result and multiple speaker_labels", func() {
			It("should coalesce the utterances", func() {
				source, err := ioutil.ReadFile("../../assets/test/ibm_3.json")
				Expect(err).To(BeNil())
				trans := parseibm.IBMTranscription{}
				err = json.Unmarshal(source, &trans)
				Expect(err).To(BeNil())
				out, err := Coerce(trans)
				Expect(err).To(BeNil())
				expectation := []Utterance{
					Utterance{
						Speaker:    0,
						From:       2.37,
						To:         7.2,
						Transcript: "design swift transaction",
					},
				}
				Expect(out).To(Equal(expectation))
			})
		})

		Context("when the IBMTranscription has multiple results and multiple speaker_labels", func() {
			It("should coalesce the utterances", func() {
				source, err := ioutil.ReadFile("../../assets/test/ibm_4.json")
				Expect(err).To(BeNil())
				trans := parseibm.IBMTranscription{}
				err = json.Unmarshal(source, &trans)
				Expect(err).To(BeNil())
				out, err := Coerce(trans)
				Expect(err).To(BeNil())
				expectation := []Utterance{
					Utterance{
						Speaker:    0,
						From:       2.37,
						To:         4.03,
						Transcript: "design",
					},
					Utterance{
						Speaker:    1,
						From:       4.09,
						To:         5.48,
						Transcript: "swift",
					},
					Utterance{
						Speaker:    0,
						From:       5.99,
						To:         7.2,
						Transcript: "transaction",
					},
				}
				Expect(out).To(Equal(expectation))
			})
		})

	})
	Context("#SquashSpeakerLabels", func() {
		Context("when there are adjacent speaker_labels belonging to the same speaker", func() {
			It("merges them", func() {
				source, err := ioutil.ReadFile("../../assets/test/speaker_labels.json")
				Expect(err).To(BeNil())
				sls := []parseibm.SpeakerLabel{}
				err = json.Unmarshal(source, &sls)
				Expect(err).To(BeNil())
				expectation := []parseibm.SpeakerLabel{
					{
						Confidence: 0.488,
						Final:      false,
						From:       2.37,
						Speaker:    0,
						To:         4.03,
					},
					{
						Confidence: 0.639,
						Final:      false,
						From:       4.09,
						Speaker:    1,
						To:         5.48,
					},
					{
						Confidence: 0.344,
						Final:      false,
						From:       5.99,
						Speaker:    0,
						To:         9.55,
					},
				}
				out, err := SquashSpeakerLabels(sls)
				Expect(err).To(BeNil())
				Expect(out).To(Equal(expectation))
			})
		})
	})
	Context("#CalcTotals", func() {
		Context("when there are no utterances", func() {
			It("returns an empty Summary object", func() {
				expectation := Summary{}
				out, err := CalcTotals([]Utterance{})
				Expect(err).To(BeNil())
				Expect(out).To(Equal(expectation))
			})
		})
	})
})
