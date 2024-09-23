package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/pterm/pterm"
)

var isQuiet bool
var isQuieter bool
var showExecutedCommands bool

// Default prefix texts
var (
	infoPrefixText    = "INFO"
	executePrefixText = "EXECUTE"
	warningPrefixText = "WARNING"
	errorPrefixText   = "ERROR"
	debugPrefixText   = "DEBUG"
)

// Stdin returns os.Stdin
func Stdin() *os.File {
	return os.Stdin
}

// Stdout returns os.Stdout or os.DevNull depending on the verbosity level
func Stdout() io.Writer {
	if isQuieter {
		return io.Discard
	}
	return os.Stdout
}

// Stderr returns os.Stderr or os.DevNull depending on the verbosity level
func Stderr() io.Writer {
	if isQuieter {
		return io.Discard
	}
	return os.Stderr
}

// Disables info messages
func SetQuiet(quiet bool) {
	isQuiet = quiet
}

// Disables commands output
func SetQuieter(quiet bool) {
	if quiet {
		SetQuiet(quiet)
	}
	isQuieter = quiet

}

func SetPrefix(prefix string) {
	infoPrefixText = prefix
	executePrefixText = prefix
}

func SetDebug(debug bool) {
	if debug {
		pterm.EnableDebugMessages()
	} else {
		pterm.DisableDebugMessages()
	}
}

func ShowExecutedCommands(showCommands bool) {
	showExecutedCommands = showCommands
}

func Info(format string, args ...any) {
	if isQuiet {
		return
	}
	pterm.Info.Prefix = pterm.Prefix{
		Text:  infoPrefixText,
		Style: pterm.NewStyle(pterm.BgGreen, pterm.FgBlack),
	}
	pterm.Info.MessageStyle = pterm.NewStyle(pterm.FgDefault)
	msg := fmt.Sprintf(format, args...)
	pterm.Info.Println(msg)
}

func Warning(format string, args ...any) {
	if isQuieter {
		return
	}
	pterm.Warning.Prefix = pterm.Prefix{
		Text:  warningPrefixText,
		Style: pterm.NewStyle(pterm.BgYellow, pterm.FgBlack),
	}
	pterm.Info.MessageStyle = pterm.NewStyle(pterm.FgYellow)
	msg := fmt.Sprintf(format, args...)
	pterm.Warning.Println(msg)
}

func Debug(format string, args ...any) {
	pterm.Debug.Prefix = pterm.Prefix{
		Text:  debugPrefixText,
		Style: pterm.NewStyle(pterm.BgGray, pterm.FgBlack),
	}
	pterm.Info.MessageStyle = pterm.NewStyle(pterm.FgDefault)
	msg := fmt.Sprintf(format, args...)
	pterm.Debug.Println(msg)
}

func Error(format string, args ...any) {
	pterm.Error.Prefix = pterm.Prefix{
		Text:  errorPrefixText,
		Style: pterm.NewStyle(pterm.BgRed, pterm.FgLightWhite),
	}
	pterm.Error.MessageStyle = pterm.NewStyle(pterm.FgRed)
	msg := fmt.Sprintf(format, args...)
	pterm.Error.Println(msg)
}

func Execute(msg string) {
	if isQuiet || !showExecutedCommands {
		return
	}
	pterm.Info.Prefix = pterm.Prefix{
		Text:  executePrefixText,
		Style: pterm.NewStyle(pterm.BgGreen, pterm.FgBlack),
	}
	pterm.Info.MessageStyle = pterm.NewStyle(pterm.FgGray)
	pterm.Info.Println(msg)
}
