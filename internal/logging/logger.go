package logging

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func GetLogger(logFileDir string) *Logger {
	// Setup loging
	file, err := os.OpenFile(logFileDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	multi := io.MultiWriter(file, os.Stdout)

	traceLog := log.New(multi,
		"[TRACE] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	infoLog := log.New(multi,
		"[INFO] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	warningLog := log.New(multi,
		"[WARNING] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	errorLog := log.New(multi,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile)

	out := Logger{Trace: traceLog, Info: infoLog, Warning: warningLog, Error: errorLog}

	// Add starting marker to log file
	log.New(file, "----------\n[STARTING] ", log.Ldate|log.Ltime).Println("")

	return &out
}
