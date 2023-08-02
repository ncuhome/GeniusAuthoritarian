package linuxUser

import (
	"github.com/Mmx233/tool"
	"os"
	"os/exec"
	"path"
)

// 适用于 alpine

func UserHomePath(username string) string {
	return "/home/" + username
}

func Exist(username string) (bool, error) {
	return tool.File.Exists(UserHomePath(username))
}

func Create(username string) error {
	return exec.Command("adduser", "-D", "-s", "/bin/sh", "-G", "common", "-h", UserHomePath(username), username).Run()
}

func Delete(username string) error {
	return exec.Command("deluser", "--remove-home", username).Run()
}

func Chown(path, username string) error {
	return exec.Command("chown", username, path).Run()
}

func PrepareSshDir(username string) error {
	dirPath := path.Join(UserHomePath(username), ".ssh")
	if exist, err := tool.File.Exists(dirPath); err != nil {
		return err
	} else if !exist {
		err = os.Mkdir(dirPath, 700)
		if err != nil {
			return err
		}
		return Chown(dirPath, username)
	}
	return nil
}

func WriteAuthorizedKeys(username, publicKey string) error {
	filePath := path.Join(UserHomePath(username), ".ssh", "authorized_keys")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 400)
	if err != nil {
		return err
	}

	if _, err = f.WriteString(publicKey); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return Chown(filePath, username)
}

func StartSshd() error {
	command := exec.Command("/usr/sbin/sshd", "-D", "-e")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Start()
}
