package rpc

import "fmt"

func LinuxAccountName(uid uint) string {
	return "home" + fmt.Sprint(uid)
}
