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
	logFile      *os.File    // Log file
	logFileMutex sync.Mutex  // Mutex for file handling
	logChannel   chan string // Channel for log entries
	currentHour  string      // Current hour for log file naming
}

func init() {
	defaultLogger = NewLogger()
	go defaultLogger.startFileWriter() // Start the goroutine for log writing
}

func NewLogger() *Logger {
	logger := &Logger{
		level:      LevelInfo,
		w:          os.Stderr,
		showDetail: false,
		logChannel: make(chan string, 100), // Buffered channel to avoid blocking
	}
	return logger
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

func GetLevel() Level {
	return defaultLogger.GetLevel()
}

func Info(format string, v ...any) {
	defaultLogger.Info(format, v...)
}

func Debug(format string, v ...any) {
	defaultLogger.Debug(format, v...)
}

func Error(format string, v ...any) {
	defaultLogger.Error(format, v...)
}

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
	msg := l.assembleMsg(format, v...)
	l.w.Write([]byte(InfoLevel + msg)) // Write to standard output
	l.logChannel <- "[INFO]" + msg     // Send log to channel for file writing
}

func (l *Logger) Debug(format string, v ...any) {
	if l.level > LevelDebug {
		return
	}
	msg := l.assembleMsg(format, v...)
	l.w.Write([]byte(DebugLevel + msg))
	l.logChannel <- "[DEBUG]" + msg
}

func (l *Logger) Error(format string, v ...any) {
	if l.level > LevelError {
		return
	}
	msg := l.assembleMsg(format, v...)
	l.w.Write([]byte(ErrorLevel + msg))
	l.logChannel <- "[ERROR]" + msg
}

func (l *Logger) AddProcessor(p Processor) {
	l.processors = append(l.processors, p)
}

func (l *Logger) assembleMsg(format string, v ...any) string {
	var msg strings.Builder
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

func (l *Logger) startFileWriter() {
	for msg := range l.logChannel {
		l.writeToFile(msg)
	}
}

func (l *Logger) writeToFile(msg string) {
	l.logFileMutex.Lock()
	defer l.logFileMutex.Unlock()

	currentHour := time.Now().Format("2006-01-02_15")
	if l.logFile == nil || l.currentHour != currentHour {
		if l.logFile != nil {
			l.logFile.Close()
		}
		filePath := fmt.Sprintf("log_%s.log", currentHour)
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		l.logFile = file
		l.currentHour = currentHour
	}

	if l.logFile != nil {
		l.logFile.WriteString(msg)
	}
}
