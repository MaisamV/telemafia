package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	// DEBUG level for detailed information
	DEBUG Level = iota
	// INFO level for general information
	INFO
	// WARN level for warning messages
	WARN
	// ERROR level for error messages
	ERROR
)

// Logger represents the logger instance
type Logger struct {
	level  Level
	logger *log.Logger
}

// New creates a new Logger instance
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level >= l.level {
		prefix := fmt.Sprintf("[%s] %s ", time.Now().Format("2006-01-02 15:04:05"), getLevelString(level))
		message := fmt.Sprintf(format, args...)
		l.logger.Println(prefix + message)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func getLevelString(level Level) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
} 