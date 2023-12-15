package core

import "os"

func MkDir() {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	err := makeDirNotExists("data/")
	err = makeDirNotExists("data/db/")
	err = makeDirNotExists("data/logs/")
	if err != nil {
		return
	}

}

func makeDirNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
