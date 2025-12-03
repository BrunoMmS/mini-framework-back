package server

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	ERROR
	DEBUG
)

type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
	level  LogLevel
}

func InitLogger(level LogLevel) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
		level:  level,
	}
}

func (l *Logger) log(level LogLevel, msg string) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var lvl string

	switch level {
	case INFO:
		lvl = "[INFO]"
	case ERROR:
		lvl = "[ERROR]"
	case DEBUG:
		lvl = "[DEBUG]"
	}

	l.logger.Printf("%s %s %s", timestamp, lvl, msg)
}

func (l *Logger) LogRequest(method, path, ip string) {
	l.log(INFO, fmt.Sprintf("Request %s %s from %s", method, path, ip))
}

func (l *Logger) LogError(method, path string, err error, ip string) {
	l.log(ERROR, fmt.Sprintf("Error %s %s from %s: %s", method, path, ip, err.Error()))
}

func (l *Logger) Debug(msg string) {
	l.log(DEBUG, msg)
}
