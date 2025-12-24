package util

import (
	"errors"
	"os"
)

func CreateDirNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		e := os.MkdirAll(dir, os.ModePerm)
		if e != nil {
			return errors.New("Error creating directory: " + e.Error())
		}
	}
	return nil
}
