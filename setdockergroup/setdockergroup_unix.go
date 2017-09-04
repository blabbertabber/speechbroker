// +build !windows

package setdockergroup

import (
	"log"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
)

func SetDockerGroup() {
	dockerGroup, err := user.LookupGroup("docker")
	if err != nil {
		log.Println("I couldn't lookup the group 'docker', error: " + err.Error())
		return
	}
	gid, err := strconv.Atoi(dockerGroup.Gid)
	if err != nil {
		panic("I couldn't convert '" + dockerGroup.Gid + "' to integer, error: " + err.Error())
	}
	if runtime.GOOS != "windows" {
		err = syscall.Setgroups([]int{gid})
		if err != nil {
			panic("I couldn't setGroups() to 'docker', make sure I'm in the docker group in /etc/group, error: " + err.Error())
		}
	}
}
