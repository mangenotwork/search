package utils

import "os"

func Mkdir(term string) {
	err := os.MkdirAll(term, os.ModePerm)
	if err != nil {
		return
	}
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
