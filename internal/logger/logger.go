package logger

import (
	"log"
	"os"
)

const (
	LevelCritical = iota
	LevelError
	LevelInfo
	LevelDebug
)

type Logger struct {
	stdout *log.Logger
	stderr *log.Logger
	level  int
}

var l *Logger

func Init(level int) {
	l = &Logger{
		stdout: log.New(os.Stdout, "", log.LstdFlags),
		stderr: log.New(os.Stderr, "", log.LstdFlags),
		level:  level,
	}
}

func Fatalf(format string, v ...interface{}) {
	if l.level >= LevelCritical {
		l.stderr.Fatalf("[CRIT] "+format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if l.level >= LevelError {
		l.stderr.Printf("[ERR] "+format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	if l.level >= LevelInfo {
		l.stdout.Printf("[INFO] "+format, v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if l.level >= LevelDebug {
		l.stdout.Printf("[DEBUG] "+format, v...)
	}
}
