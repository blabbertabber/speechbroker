// cmdrunner runs commands, typically Docker commands but can be anything.
// it does not return any data and will panic() if the command doesn't start
// or exits with a non-zero value.
package cmdrunner

import (
	"io/ioutil"
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
	stderr, err := command.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := command.Start(); err != nil {
		panic(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	log.Printf("%s\n", slurp)

	if err := command.Wait(); err != nil {
		panic(err)
	}
}
