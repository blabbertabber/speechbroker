package diarizerrunner_test

import (
	. "github.com/blabbertabber/speechbroker/diarizerrunner"

	"errors"
	"github.com/blabbertabber/speechbroker/cmdrunner/cmdrunnerfakes"
	"github.com/blabbertabber/speechbroker/ibmservicecreds"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("diarizerrunner", func() {
	var fdr *cmdrunnerfakes.FakeCmdRunner
	var r Runner
	var creds ibmservicecreds.IBMServiceCreds

	BeforeEach(func() {
		fdr = new(cmdrunnerfakes.FakeCmdRunner)
		r = Runner{CmdRunner: fdr}
		creds = ibmservicecreds.IBMServiceCreds{
			Username: "fake-ibm-username",
			Password: "fake-ibm-password",
		}
	})

	Context("When the runner is \"Aalto\"", func() {
		It("runs Docker with the correct arguments", func() {
			Expect(r.Run("Aalto", "fake-uuid", creds)).To(BeNil())
			Expect(fdr.RunArgsForCall(0)).To(Equal([]string{
				"docker",
				"run",
				"--volume=/var/blabbertabber:/blabbertabber",
				"--workdir=/speaker-diarization",
				"blabbertabber/aalto-speech-diarizer",
				"/speaker-diarization/spk-diarization2.py",
				"/blabbertabber/soundFiles/fake-uuid/meeting.wav",
				"-o",
				"/blabbertabber/diarizationResults/fake-uuid/diarization.txt",
			}))
		})
	})
	Context("When the runner is \"CMUSphinx4\"", func() {
		It("runs Docker with the correct arguments", func() {
			Expect(r.Run("CMUSphinx4", "fake-uuid", creds)).To(BeNil())
			Expect(fdr.RunArgsForCall(0)).To(Equal([]string{
				"docker",
				"run",
				"--volume=/var/blabbertabber:/blabbertabber",
				"blabbertabber/cmu-sphinx4-transcriber",
				"java",
				"-Xmx2g",
				"-cp",
				"/sphinx4-5prealpha-src/sphinx4-core/build/libs/sphinx4-core-5prealpha-SNAPSHOT.jar:/sphinx4-5prealpha-src/sphinx4-data/build/libs/sphinx4-data-5prealpha-SNAPSHOT.jar:.",
				"Transcriber",
				"/blabbertabber/soundFiles/fake-uuid/meeting.wav",
				"/blabbertabber/diarizationResults/fake-uuid/transcription.txt",
			}))
		})
	})
	Context("When the runner is \"IBM\"", func() {
		It("runs several commands with the correct arguments", func() {
			Expect(r.Run("IBM", "fake-uuid", creds)).To(BeNil())
			Expect(fdr.RunArgsForCall(0)).To(Equal([]string{
				"bash",
				"-c",
				"echo /blabbertabber/soundFiles/fake-uuid/meeting.wav > /var/blabbertabber/soundFiles/fake-uuid/wav_file_list.txt",
			}))
			Expect(fdr.RunArgsForCall(1)).To(Equal([]string{
				"docker",
				"run",
				"--volume=/var/blabbertabber:/blabbertabber",
				"blabbertabber/ibm-watson-stt",
				"python",
				"/speech-to-text-websockets-python/sttClient.py",
				"-credentials",
				"fake-ibm-username:fake-ibm-password",
				"-model",
				"en-US_NarrowbandModel",
				"-in",
				"/blabbertabber/soundFiles/fake-uuid/wav_file_list.txt",
				"-out",
				"/blabbertabber/diarizationResults/fake-uuid/ibm_out",
			}))
			Expect(fdr.RunArgsForCall(2)).To(Equal([]string{
				"/usr/local/bin/ibmjson",
				"-in",
				"/var/blabbertabber/diarizationResults/fake-uuid/ibm_out/0.json.txt",
				"-out",
				"/var/blabbertabber/diarizationResults/fake-uuid/ibm_out.json",
			}))
		})
		Context("when there are no credentials passed", func() {
			It("should return an error", func() {
				Expect(r.Run("IBM", "fake-uuid", ibmservicecreds.IBMServiceCreds{})).To(Equal(errors.New("invalid IBM creds")))
			})
		})
	})
	Context("When the runner is \"null\"", func() {
		It("doesn't do anything", func() {
			Expect(r.Run("null", "fake-uuid", creds)).To(BeNil())
			Expect(fdr.RunCallCount()).To(Equal(0))
		})
	})
	Context("When the runner is \"non-existent\"", func() {
		It("panics", func() {
			Expect(r.Run("non-existent", "fake-uuid", creds)).To(Equal(errors.New("No such back-end: \"non-existent\"")))
		})
	})
})
