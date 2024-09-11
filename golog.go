package golog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError

	Reset      = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	Gray       = "\033[37m"
	White      = "\033[97m"
	Whitespace = " "
	Newline    = "\n"

	InfoLevel  = Green + "[INFO]" + Reset
	DebugLevel = Yellow + "[DEBUG]" + Reset
	ErrorLevel = Red + "[ERROR]" + Reset
)

var (
	defaultLogger *Logger
)

type Level int32

type Processor func(format string, v ...any) (string, []any)

type Logger struct {
	level        Level
	prefix       string
	fileLocation string
	showDetail   bool
	mutex        sync.Mutex
	buf          bytes.Buffer
	w            io.Writer
	processors   []Processor
}

func init() {
	defaultLogger = NewLogger()
}

/*
NewLogger
By convention, many programs output their log messages to os.Stderr
instead of os.Stdout for a couple of reasons:

1. Separation of Concerns: It allows for the separation of regular
program output from error messages and log information.
This can be very useful when you’re piping or redirecting output
in the command line.
2. Order of Messages: os.Stderr is unbuffered, while os.Stdout is buffered.
This means that if your program crashes or exits unexpectedly, messages
sent to os.Stdout might not get printed if the buffer isn’t flushed before the program exits.
But messages sent to os.Stderr will always get printed immediately.
*/
func NewLogger() *Logger {
	return &Logger{
		level:      LevelInfo,
		w:          os.Stderr,
		showDetail: false,
	}
}

/*
SetLevel
Enum of different levels:

	Debug: 0
	INFO:  1
	ERROR: 2
*/
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

/*
GetLevel
*/
func GetLevel() Level {
	return defaultLogger.GetLevel()
}

/*
Info
*/
func Info(format string, v ...any) {
	defaultLogger.Info(format, v...)
}

/*
Debug
*/
func Debug(format string, v ...any) {
	defaultLogger.Debug(format, v...)
}

/*
Error
*/
func Error(format string, v ...any) {
	defaultLogger.Error(format, v...)
}

/*
Add customized processor functions
*/
func AddProcessor(p Processor) {
	defaultLogger.AddProcessor(p)
}

func ShowDetail(b bool) {
	defaultLogger.showDetail = b
}

func (l *Logger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&l.level), int32(level))
}

func (l *Logger) GetLevel() Level {
	return Level(atomic.LoadInt32((*int32)(&l.level)))
}

func (l *Logger) Info(format string, v ...any) {
	if l.level > LevelInfo {
		return
	}
	l.w.Write([]byte(l.assembleMsg(InfoLevel, format, v...)))
}

func (l *Logger) Debug(format string, v ...any) {
	if l.level > LevelDebug {
		return
	}
	l.w.Write([]byte(l.assembleMsg(DebugLevel, format, v...)))
}

func (l *Logger) Error(format string, v ...any) {
	if l.level > LevelError {
		return
	}
	l.w.Write([]byte(l.assembleMsg(ErrorLevel, format, v...)))
}

func (l *Logger) AddProcessor(p Processor) {
	l.processors = append(l.processors, p)
}

func (l *Logger) assembleMsg(logLevel string, format string, v ...any) string {

	// https://golangnote.com/golang/golang-stringsbuilder-vs-bytesbuffer
	// Both `strings.Builder` and `bytes.Buffer` are used for efficient in Golang.
	// Here I choose strings.Builder

	var msg strings.Builder
	msg.WriteString(logLevel)
	msg.WriteString(Whitespace)

	if l.showDetail {
		msg.WriteString(time.Now().String())
		msg.WriteString(Whitespace)
		getFileLocation := func() string {
			_, file, line, ok := runtime.Caller(4)
			if !ok {
				file = "unknown file"
				line = -1
			}
			return fmt.Sprintf("%s:%d", filepath.Base(file), line) + " "
		}

		msg.WriteString(getFileLocation())
	}

	msg.WriteString(l.getContent(format, v...))
	msg.WriteString(Whitespace)
	msg.WriteString(Newline)

	return msg.String()
}

func (l *Logger) getContent(format string, v ...any) string {
	for _, process := range l.processors {
		format, v = process(format, v...)
	}
	return fmt.Sprintf(format, v...)
}
