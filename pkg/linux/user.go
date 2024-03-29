package linux

import (
	"os/exec"
)

// 适用于 alpine

func UserHomePath(username string) string {
	return "/home/" + username
}

func CreateUser(username string) error {
	return exec.Command("adduser", "-D", "-s", "/bin/sh", "-G", "common", "-h", UserHomePath(username), username).Run()
}

func DelUserPasswd(username string) error {
	return exec.Command("passwd", "-d", username).Run()
}

func DeleteUser(username string) error {
	return exec.Command("deluser", "--remove-home", username).Run()
}

func UserKillAll(username string) error {
	return exec.Command("pkill", "-u", username).Run()
}
