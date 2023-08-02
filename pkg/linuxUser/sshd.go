package linuxUser

import (
	"os"
	"os/exec"
)

func StartSshd() error {
	command := exec.Command("/usr/sbin/sshd", "-D", "-e")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Start()
}
