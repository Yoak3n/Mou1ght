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
	initDirectory()
	file, err := os.OpenFile("logs/errors.log",
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

func INFO(msg ...string) {
	Info.Println(msg)
}
func WARNING(msg ...string) {
	Warning.Println(msg)
}
func ERROR(msg ...string) {
	Error.Println(msg)
}
func TRACE(msg ...string) {
	Trace.Println(msg)
}

// because can't init in main.go before in logger.go ,so the function moves here
func initDirectory() {
	_, err := os.Stat("data")
	if err != nil {
		log.Println("data directory didn't exist")
		err = os.Mkdir("data", os.ModePerm)
		if err != nil {
			panic("Can't create directory 'data'")
		}
	}
	_, err = os.Stat("logs")
	if err != nil {
		if err != nil {
			log.Println("logs directory didn't exist")
			err = os.Mkdir("logs", os.ModePerm)
			if err != nil {
				panic("Can't create directory 'logs")
			}
		}
	}
}
