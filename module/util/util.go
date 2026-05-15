package util

import "strings"

func TopLevel(vpath string) string {
	args := strings.Split(vpath, "/")
	var arg string
	for _, a := range args {
		if a != "" {
			arg = a
			break
		}
	}

	return "/" + arg
}
