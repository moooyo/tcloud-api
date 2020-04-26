package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	time2 "time"
)

const defaultConfigPath = "logs"

var logFile *os.File = nil
var logLevel = 0
var logger chan string

func getLogFile() *os.File {
	if logFile == nil {
		path := defaultConfigPath
		config := GetConfig().Log
		if config.Path != "" {
			path = config.Path
		}
		logLevel = config.Level

		logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal("Error open log file error.\n", err)
		}
		return logFile
	}

	return logFile
}

func int2level(level int) string {
	switch level {
	case 1:
		return "INFO"
	case 2:
		return "WARN"
	case 3:
		return "ERROR"
	default:
		return "????"
	}
}

func formatWriteLog(level int, format string, args ...interface{}) {
	if logLevel > level {
		return
	}
	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(2)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")

		filename = filepath.Base(filename)
	}
	time := time2.Now().Format("2016-01-01 12:12:12")
	str := fmt.Sprintf("[%s] %s:%d:%s %s %s\n", int2level(level), filename, line, funcname, time, fmt.Sprintf(format, args))
	logger <- str
}

func INFO(format string, args ...interface{}) {
	formatWriteLog(1, format, args)
}
func WARN(format string, args ...interface{}) {
	formatWriteLog(2, format, args)
}
func ERROR(format string, args ...interface{}) {
	formatWriteLog(3, format, args)
}

func InitLogger() {
	logger = make(chan string, 1024)
	file := getLogFile()
	go func() {
		str := <-logger
		_, err := file.WriteString(str)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}()
}
