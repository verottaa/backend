package logger

import (
	"log"
	"os"
	"strings"
)

type Logger struct {
	logger *log.Logger
	tag    string
}

type Logging interface {
	Error(err error)
	Info(info ...string)
	Warn(warn ...string)
}

func (l *Logger) Error(err error) {
	l.logger.Println("[ERROR][", l.tag, "]: ", err.Error())
}

func (l *Logger) Info(infos ...string) {
	messageArray := []string{"[INFO][", l.tag, "]: "}
	for _, v := range infos {
		messageArray = append(messageArray, v)
	}
	l.logger.Println(strings.Join(messageArray, ""))
}

func (l *Logger) Warn(warns ...string) {
	messageArray := []string{"[WARNING][", l.tag, "]: "}
	for _, v := range warns {
		messageArray = append(messageArray, v)
	}
	l.logger.Println(strings.Join(messageArray, ""))
}

func CreateLogger(tag string) *Logger {
	logger := new(Logger)
	logger.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.tag = tag
	return logger
}
