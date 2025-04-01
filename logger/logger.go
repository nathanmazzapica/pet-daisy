package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var ErrLog = log.New(os.Stderr, "[ ERROR ]", log.LstdFlags|log.Ldate)

var errFile *os.File

func InitLog() {
	var err error

	currentTime := time.Now()
	formatted := currentTime.Format("2006-01-02-15")

	fileName := fmt.Sprintf("./logs/%s.log", formatted)

	errFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	ErrLog.SetOutput(errFile)
}

func CloseLog() {
	errFile.Close()
}

func LogError(err error) {
	ErrLog.Println(err)
}

func LogFatalError(err error) {
	ErrLog.Fatalln(err)
}
