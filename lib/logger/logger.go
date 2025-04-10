package logger

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

// Global Log variable
var (
	Log     *pterm.Logger
	Testing bool = false
)

// SetLogger sets the default logger verbosity level and returns a logger instance
func SetLogger(verbosity int) *pterm.Logger {
	pterm.EnableColor()

	if !viper.GetBool("color") {
		pterm.DisableColor()
	}

	Log = pterm.DefaultLogger.WithLevel(getLevel(verbosity))

	return Log
}

// Helper Functions
func Panic(err error) {
	Log.Error(err.Error())
	panic(err)
}

func Fatal(msg string, args ...any) {
	err := fmt.Errorf(msg, args...)

	// Panic if testing so we can catch it
	if Testing {
		Panic(err)
	} else {
		Log.Error(err.Error())
		os.Exit(1)
	}
}

func Error(msg string, args ...any) {
	Log.Error(fmt.Sprintf(msg, args...))
}

func Warn(msg string, args ...any) {
	Log.Warn(fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...any) {
	Log.Info(fmt.Sprintf(msg, args...))
}

func Debug(msg string, args ...any) {
	Log.Debug(fmt.Sprintf(msg, args...))
}

func Trace(msg string, args ...any) {
	Log.Trace(fmt.Sprintf(msg, args...))
}

// Private functions

func getLevel(verbosity int) pterm.LogLevel {
	if verbosity >= 3 {
		return pterm.LogLevelTrace
	}

	if verbosity >= 2 {
		return pterm.LogLevelDebug
	}

	if verbosity >= 1 {
		return pterm.LogLevelInfo
	}

	return pterm.LogLevelError
}

func init() {
	if Log == nil {
		SetLogger(0)
	}
}
