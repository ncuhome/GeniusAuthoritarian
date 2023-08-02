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
	return exec.Command("adduser", "-s", "/bin/sh", "-G", "common", username).Run()
}

func PrepareSshDir(username string) error {
	dirPath := path.Join(UserHomePath(username), ".ssh")
	if exist, err := tool.File.Exists(dirPath); err != nil {
		return err
	} else if !exist {
		return os.Mkdir(dirPath, 700)
	}
	return nil
}

func WriteAuthorizedKeys(username, publicKey string) error {
	f, err := os.OpenFile(path.Join(UserHomePath(username), ".ssh", "authorized_keys"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 400)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(publicKey)
	return err
}
