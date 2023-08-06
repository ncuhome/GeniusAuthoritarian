package linuxUser

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

func StartSshd() error {
	command := exec.Command("/usr/sbin/sshd", "-D", "-e")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func DaemonSshd() {
	for {
		err := StartSshd()
		log.Errorln("sshd error: ", err)
		time.Sleep(time.Second * 2)
	}
}
