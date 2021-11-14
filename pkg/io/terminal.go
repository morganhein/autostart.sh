package io

import (
	"github.com/morganhein/autostart.sh/pkg/T"
	"go.uber.org/zap"
)

type LogLevel string

const (
	Informational LogLevel = "informational"
	Warning       LogLevel = "warning"
	Error         LogLevel = "error"
)

// Terminal A terminal/shell with user, to prompt questions from user and print/log
type Terminal interface {
	Logger
	//Prompts the user to select any of the given options, with the specific option
	PromptUser(prompt string, options []string, defaultSelection string) (string, error)
}

type Logger interface {
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type logger struct {
	zap *zap.SugaredLogger
}

func (l logger) Debugf(format string, args ...interface{}) {
	l.zap.Debugf(format, args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.zap.Infof(format, args...)
}

func (l logger) Warningf(format string, args ...interface{}) {
	l.zap.Warnf(format, args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	l.zap.Errorf(format, args)
}

func NewLogger() Logger {
	z, err := zap.NewDevelopment(zap.Development())
	_ = T.Log(err)
	return &logger{
		z.Sugar(),
	}
}
