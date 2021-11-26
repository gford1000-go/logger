package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

// CorrelatedLog provides logging with a level and correlation id for each message
type CorrelatedLog func(lvl LogLevel, correlationID, format string, a ...interface{})

type logger struct {
	lvl LogLevel
	*log.Logger
}

func (l *logger) log(lvl LogLevel, correlationID, format string, a ...interface{}) {
	if lvl > singletonLogger.lvl {
		return
	}

	_, file, line, _ := runtime.Caller(2)
	fileParts := strings.Split(file, "/")

	singletonLogger.Println(fmt.Sprintf("%v:%v %v %v %v", fileParts[len(fileParts)-1], line, correlationID, lvl, fmt.Sprintf(format, a...)))
}

var singletonLogger *logger
var once sync.Once

// NewFileLogger creates a new log file at te specified location, which will contain
// entries up to the specified LogLevel
func NewFileLogger(fname string, lvl LogLevel, prefix string) (CorrelatedLog, error) {
	file, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
	if err != nil {
		return nil, err
	}
	return NewLogger(file, lvl, prefix), nil
}

// NewLogger writes entries to the specified io.Writer, if the level of the statement
// to be logged is less than the specified LogLevel
func NewLogger(w io.Writer, lvl LogLevel, prefix string) CorrelatedLog {
	once.Do(func() {
		singletonLogger = &logger{
			lvl:    lvl,
			Logger: log.New(w, prefix, log.Ldate|log.Ltime|log.Lmicroseconds),
		}
	})
	return func(lvl LogLevel, correlationID, format string, a ...interface{}) {
		singletonLogger.log(lvl, correlationID, format, a...)
	}
}

// GetLogger returns access to the CorrelatedLog function for the initialised logging
func GetLogger() CorrelatedLog {
	return singletonLogger.log
}
