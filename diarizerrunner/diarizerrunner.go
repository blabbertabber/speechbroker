// diarizerrunner invokes the commands (via cmdrunner) necessary to
// run the diarization/transcription back-ends (e.g. Aalto, CMU Sphinx 4,
// IBM Speech to Text).
package diarizerrunner

import (
	"errors"
	"fmt"
	"github.com/blabbertabber/speechbroker/cmdrunner"
	"github.com/blabbertabber/speechbroker/ibmservicecreds"
)

type Runner struct {
	CmdRunner cmdrunner.CmdRunner
}

// `counterfeiter diarizerrunner/diarrizerrunner.go DiarizerRunner`
type DiarizerRunner interface {
	Run(flavor, uuid string, creds ibmservicecreds.IBMServiceCreds) error
}

func (r Runner) Run(flavor, meetingUuid string, creds ibmservicecreds.IBMServiceCreds) (err error) {
	err = nil
	switch flavor {
	case "Aalto":
		{
			AaltoCmd := []string{
				"docker",
				"run",
				"--volume=/var/blabbertabber:/blabbertabber",
				"--workdir=/speaker-diarization",
				"blabbertabber/aalto-speech-diarizer",
				"/speaker-diarization/spk-diarization2.py",
				fmt.Sprintf("/blabbertabber/soundFiles/%s/meeting.wav", meetingUuid),
				"-o",
				fmt.Sprintf("/blabbertabber/diarizationResults/%s/diarization.txt", meetingUuid),
			}
			r.CmdRunner.Run(AaltoCmd...)
		}
	case "CMUSphinx4":
		{
			CMUSphinx4Cmd := []string{
				"docker",
				"run",
				"--volume=/var/blabbertabber:/blabbertabber",
				"blabbertabber/cmu-sphinx4-transcriber",
				"java",
				"-Xmx2g",
				"-cp",
				"/sphinx4-5prealpha-src/sphinx4-core/build/libs/sphinx4-core-5prealpha-SNAPSHOT.jar:/sphinx4-5prealpha-src/sphinx4-data/build/libs/sphinx4-data-5prealpha-SNAPSHOT.jar:.",
				"Transcriber",
				fmt.Sprintf("/blabbertabber/soundFiles/%s/meeting.wav", meetingUuid),
				fmt.Sprintf("/blabbertabber/diarizationResults/%s/transcription.txt", meetingUuid),
			}
			r.CmdRunner.Run(CMUSphinx4Cmd...)
		}
	case "IBM":
		{
			if creds.Username == "" || creds.Password == "" {
				return errors.New("invalid IBM creds")
			}
			IBMCmd := []string{
				"docker",
				"run",
				"--volume=/var/blabbertabber:/blabbertabber",
				"blabbertabber/ibm-watson-stt",
				"python",
				"/sttClient.py",
				"-credentials",
				creds.Username + ":" + creds.Password,
				"-model",
				"en-US_NarrowbandModel",
				"-in",
				fmt.Sprintf("/blabbertabber/soundFiles/%s/meeting.wav", meetingUuid),
				"-out",
				fmt.Sprintf("/blabbertabber/diarizationResults/%s", meetingUuid),
			}
			r.CmdRunner.Run(IBMCmd...)
		}
	default:
		{
			err = errors.New(fmt.Sprintf("No such back-end: \"%s\"", flavor))
		}
	}
	return
}
