package emitblabbertabber_test

import (
	. "github.com/blabbertabber/speechbroker/IBMJson/emitblabbertabber"

	"encoding/json"
	"github.com/blabbertabber/speechbroker/IBMJson/parseibm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"sort"
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
	Context("#CalcSummary", func() {
		Context("when there are no utterances", func() {
			It("returns an empty Summary object", func() {
				expectation := Summary{
					LeaderBoard: []SpeakerStat{},
					Utterances:  []Utterance{},
				}
				out, err := CalcSummary([]Utterance{})
				Expect(err).To(BeNil())
				Expect(out).To(Equal(expectation))
			})
		})
		Context("when there is one utterance", func() {
			It("should return a correct Summary with one Utterance", func() {
				source, err := ioutil.ReadFile("../../assets/test/ibm_1.json")
				Expect(err).To(BeNil())
				trans := parseibm.IBMTranscription{}
				err = json.Unmarshal(source, &trans)
				utterances, err := Coerce(trans)
				Expect(err).To(BeNil())
				out, err := CalcSummary(utterances)
				expectation := Summary{
					TotalSpeakingTime: 15.17,
					LeaderBoard: []SpeakerStat{
						{
							Speaker:   "0",
							TotalTime: 15.17,
						},
					},
					Utterances: []Utterance{
						{
							Speaker:    0,
							From:       2.37,
							To:         17.54,
							Transcript: "design swift transaction so you go through when you put all all the things you need to do",
						},
					},
				}
				Expect(out).To(Equal(expectation))
			})
		})
		Context("when there are many utterances and multiple speakers", func() {
			It("should return a correct Summary", func() {
				source, err := ioutil.ReadFile("../../assets/test/ibm_4.json")
				Expect(err).To(BeNil())
				trans := parseibm.IBMTranscription{}
				err = json.Unmarshal(source, &trans)
				utterances, err := Coerce(trans)
				Expect(err).To(BeNil())
				out, err := CalcSummary(utterances)
				expectation := Summary{
					TotalSpeakingTime: 4.26,
					LeaderBoard: []SpeakerStat{
						{Speaker: "0", TotalTime: 2.87},
						{Speaker: "1", TotalTime: 1.39},
					},
					Utterances: []Utterance{
						{Speaker: 0, From: 2.37, To: 4.03, Transcript: "design"},
						{Speaker: 1, From: 4.09, To: 5.48, Transcript: "swift"},
						{
							Speaker:    0,
							From:       5.99,
							To:         7.2,
							Transcript: "transaction",
						},
					},
				}
				Expect(out).To(Equal(expectation))
			})
		})

	})
	Context("sort.Sort(SpeakerStats)", func() {
		Context("When given a slew of speakers", func() {
			It("should sort them in ascending order by TotalTime", func() {
				expectation := SpeakerStats{
					{Speaker: "6", TotalTime: 0.0},
					{Speaker: "5", TotalTime: 1.0},
					{Speaker: "4", TotalTime: 2.0},
					{Speaker: "3", TotalTime: 3.0},
					{Speaker: "2", TotalTime: 4.0},
					{Speaker: "1", TotalTime: 5.0},
					{Speaker: "0", TotalTime: 6.0},
				}
				out := SpeakerStats{
					{Speaker: "1", TotalTime: 5},
					{Speaker: "5", TotalTime: 1},
					{Speaker: "6", TotalTime: 0},
					{Speaker: "2", TotalTime: 4},
					{Speaker: "3", TotalTime: 3},
					{Speaker: "4", TotalTime: 2},
					{Speaker: "0", TotalTime: 6},
				}
				sort.Sort(SpeakerStats(out))
				Expect(out).To(Equal(expectation))
			})
		})
	})
})
