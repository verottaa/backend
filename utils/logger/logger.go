package logger

import "log"

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

func (l *Logger) Info(info ...string) {
	l.logger.Println("[INFO][", l.tag, "]: ", info)
}

func (l *Logger) Warn(warn ...string) {
	l.logger.Println("[WARNING][", l.tag, "]: ", warn)
}

func CreateLogger(tag string) *Logger {
	logger := new(Logger)
	logger.logger = new(log.Logger)
	logger.tag = tag
	return logger
}
