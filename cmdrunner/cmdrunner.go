// cmdrunner runs commands, typically Docker commands but can be anything.
// it does not return any data and will panic() if the command doesn't start
// or exits with a non-zero value.
package cmdrunner

import (
	"log"
	"os/exec"
)

// `counterfeiter cmdrunner/cmdrunner.go CmdRunner`
type CmdRunner interface {
	Run(cmdArgs ...string)
}

type CmdRunnerReal struct{}

func (d CmdRunnerReal) Run(cmdArgs ...string) {
	command := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	stdOutStdErr, err := command.CombinedOutput()
	log.Println(cmdArgs)
	log.Println(string(stdOutStdErr))
	if err != nil {
		log.Panic(err)
	}
}
