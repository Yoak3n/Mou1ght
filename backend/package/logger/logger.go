package logger

import (
	"Mou1ght-Server/package/core"
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
	// because can't init in main.go before in logger.go ,so the function moves here
	core.MkDir()
	file, err := os.OpenFile("data/logs/errors.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Can't open file error.log:", err)
	}
	Trace = log.New(io.Discard, "[TRACE] ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	Info = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	Warning = log.New(os.Stdout, "[WARNING] ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr), "[ERROR] ", log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
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
