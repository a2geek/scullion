package log

import (
	"fmt"
	"io/ioutil"
	golog "log"
	"os"
)

func NewLogger(prefix, level string, noDate bool) (Logger, error) {
	flags := golog.LstdFlags | golog.Lmsgprefix
	if noDate {
		flags = 0
	}
	prefix = fmt.Sprintf("[%s] ", prefix)
	switch level {
	case "DEBUG":
		return logger{
			debugLog: golog.New(os.Stdout, prefix, flags),
			infoLog:  golog.New(os.Stdout, prefix, flags),
			warnLog:  golog.New(os.Stdout, prefix, flags),
			errorLog: golog.New(os.Stderr, prefix, flags),
		}, nil
	case "INFO":
		return logger{
			debugLog: golog.New(ioutil.Discard, prefix, flags),
			infoLog:  golog.New(os.Stdout, prefix, flags),
			warnLog:  golog.New(os.Stdout, prefix, flags),
			errorLog: golog.New(os.Stderr, prefix, flags),
		}, nil
	case "WARN":
		return logger{
			debugLog: golog.New(ioutil.Discard, prefix, flags),
			infoLog:  golog.New(ioutil.Discard, prefix, flags),
			warnLog:  golog.New(os.Stdout, prefix, flags),
			errorLog: golog.New(os.Stderr, prefix, flags),
		}, nil
	case "ERROR":
		return logger{
			debugLog: golog.New(ioutil.Discard, prefix, flags),
			infoLog:  golog.New(ioutil.Discard, prefix, flags),
			warnLog:  golog.New(ioutil.Discard, prefix, flags),
			errorLog: golog.New(os.Stderr, prefix, flags),
		}, nil
	}
	return nil, fmt.Errorf("unknown log level '%s'", level)
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

type logger struct {
	debugLog *golog.Logger
	infoLog  *golog.Logger
	warnLog  *golog.Logger
	errorLog *golog.Logger
}

func (l logger) Debug(v ...interface{}) {
	l.debugLog.Output(2, fmt.Sprint(v...))
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.debugLog.Output(2, fmt.Sprintf(format, v...))
}

func (l logger) Info(v ...interface{}) {
	l.infoLog.Output(2, fmt.Sprint(v...))
}

func (l logger) Infof(format string, v ...interface{}) {
	l.infoLog.Output(2, fmt.Sprintf(format, v...))
}

func (l logger) Warn(v ...interface{}) {
	l.warnLog.Output(2, fmt.Sprint(v...))
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.warnLog.Output(2, fmt.Sprintf(format, v...))
}

func (l logger) Error(v ...interface{}) {
	l.warnLog.Output(2, fmt.Sprint(v...))
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.errorLog.Output(2, fmt.Sprintf(format, v...))
}
