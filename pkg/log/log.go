package log

import (
	"log"
)

type LogLevel int

const (
	LevelD LogLevel = -1
	LevelI LogLevel = 0
	LevelW LogLevel = 1
	LevelE LogLevel = 2
)

var Logger *log.Logger
var Level LogLevel = LevelI
var Prefix map[LogLevel]string = map[LogLevel]string{
	LevelD: "DEBUG: ",
	LevelI: "INFO: ",
	LevelW: "WARN: ",
	LevelE: "ERR!: ",
}

func E(v interface{}) {
	print(LevelE, v)
}

func Ef(format string, v ...interface{}) {
	printf(LevelE, format, v...)
}

func W(v interface{}) {
	print(LevelW, v)
}

func Wf(format string, v ...interface{}) {
	printf(LevelW, format, v...)
}

func I(v interface{}) {
	print(LevelI, v)
}

func If(format string, v ...interface{}) {
	printf(LevelI, format, v...)
}

func D(v interface{}) {
	print(LevelD, v)
}

func Df(format string, v ...interface{}) {
	printf(LevelD, format, v...)
}

func Fatal(v ...interface{}) {
	if Logger == nil {
		log.Fatal(v...)
	}
	Logger.Fatal(v...)
}

func print(level LogLevel, v interface{}) {
	if level < Level {
		return
	}

	prefix := Prefix[level]
	if Logger == nil {
		log.Print(prefix, v)
		return
	}
	Logger.Print(prefix, v)
}

func printf(level LogLevel, format string, v ...interface{}) {
	if level < Level {
		return
	}

	f := Prefix[level] + format
	if Logger == nil {
		log.Printf(f, v...)
		return
	}
	Logger.Printf(f, v...)
}
