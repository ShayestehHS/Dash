package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Code      string                 `json:"code"`
	File      string                 `json:"file"`
	Line      int                    `json:"line"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

func log(level LogLevel, code string, fields map[string]interface{}) {
	// Caller(2) skips 2 frames: 0=runtime.Caller, 1=log(), 2=Debug/Info/Warn/Error
	// This captures the actual caller of the logging function, not the internal log() function
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	} else {
		file = filepath.Base(file)
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Code:      code,
		File:      file,
		Line:      line,
		Fields:    fields,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, `{"timestamp":"%s","level":"ERROR","code":"%s","file":"%s","line":"%d","fields":"%v","error":"%v"}`+"\n", time.Now().UTC().Format(time.RFC3339), code, file, line, fields, err)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

func Debug(code string, fields map[string]interface{}) {
	log(LevelDebug, code, fields)
}

func Info(code string, fields map[string]interface{}) {
	log(LevelInfo, code, fields)
}

func Warn(code string, fields map[string]interface{}) {
	log(LevelWarn, code, fields)
}

func Error(code string, fields map[string]interface{}) {
	log(LevelError, code, fields)
}
