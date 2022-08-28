package io

import (
	"github.com/AlecAivazis/survey/v2"
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
	//PromptUser prompts the user to select any of the given options, with the specific option
	PromptUser(prompt string, options []string, defaultSelection string) (string, error)
	//Tightly coupled interface just so we can moq it
	AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
}

func NewTerminal() *terminal {
	return &terminal{
		Logger: NewLogger(),
	}
}

type terminal struct {
	Logger
}

func (t terminal) PromptUser(prompt string, options []string, defaultSelection string) (string, error) {
	panic("implement me")
}

func (t terminal) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(p, response)
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
	z, _ := zap.NewDevelopment(zap.Development())
	return &logger{
		z.Sugar(),
	}
}
