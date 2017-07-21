package parseibm_test

import (
	. "github.com/blabbertabber/speechbroker/ibmjson/parseibm"

	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Parseibm", func() {

	It("should parse an empty JSON properly", func() {
		source := []byte(`{}`)
		expectation := IBMTranscription{}
		result := IBMTranscription{}
		err := json.Unmarshal(source, &result)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectation))
	})
	It("should parse a minimal JSON properly", func() {
		source := []byte(`{"result_index": 1, "results": [], "speaker_labels": []}`)
		expectation := IBMTranscription{ResultIndex: 1, Results: []Result{}, SpeakerLabels: []SpeakerLabel{}}
		result := IBMTranscription{}
		err := json.Unmarshal(source, &result)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectation))
	})
	It("should parse a full JSON properly", func() {
		source, err := ioutil.ReadFile("../../assets/test/ibm_0.json")
		Expect(err).To(BeNil())
		expectation := IBMTranscription{
			ResultIndex: 0,
			Results: []Result{{
				Alternatives: []Alternative{{
					Confidence: 0.694,
					Timestamps: []Timestamp{{
						Word: "design",
						From: 2.37,
						To:   3.13,
					}},
					Transcript: "design",
				}},
				Final: true,
			}},
			SpeakerLabels: []SpeakerLabel{{
				Confidence: 0.488,
				Final:      false,
				From:       2.37,
				To:         3.13,
			}},
		}
		result := IBMTranscription{}
		err = json.Unmarshal(source, &result)
		Expect(err).To(BeNil())
		Expect(result).To(Equal(expectation))
	})
})
