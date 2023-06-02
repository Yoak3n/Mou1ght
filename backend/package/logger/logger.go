package logger

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	file, err := os.OpenFile("data/errors.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Can't open file error.log:", err)
	}
	Trace = log.New(io.Discard, "[TRACE] ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	Warning = log.New(os.Stdout, "[WARNING] ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr), "[ERROR] ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
}

func LogOut(msg string, flag ...int) {
	// set default output choice is info
	choice := 3
	if len(flag) != 0 {
		choice = flag[0]
	}
	switch choice {
	case 0:
		Error.Println(msg)
	case 1:
		Warning.Println(msg)
	case 3:
		Info.Println(msg)
	case 4:
		Trace.Println(msg)
	}

}
