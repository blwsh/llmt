package main

import (
	"fmt"
	"os"

	"github.com/blwsh/llmt/lib/logger"
)

type impl struct{}

func newLogger() logger.Logger {
	return &impl{}
}

func (i impl) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (i impl) Infof(template string, args ...interface{}) {
	fmt.Printf(template, args...)
}

func (i impl) Warnf(template string, args ...interface{}) {
	fmt.Printf(template, args...)
}

func (i impl) Error(args ...interface{}) {
	fmt.Printf("%v", args...)
}

func (i impl) Fatal(args ...interface{}) {
	fmt.Printf("%v", args...)
	os.Exit(1)
}