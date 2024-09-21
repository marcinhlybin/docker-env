package logger

import (
	"fmt"

	"github.com/pterm/pterm"
)

var (
	infoPrefixText    = "INFO"
	executePrefixText = "EXECUTE"
	warningPrefixText = "WARNING"
	errorPrefixText   = "ERROR"
	debugPrefixText   = "DEBUG"
	executeEnabled    = true
)

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

func ShowCommands(show bool) {
	executeEnabled = show
}

func Info(format string, args ...any) {
	pterm.Info.Prefix = pterm.Prefix{
		Text:  infoPrefixText,
		Style: pterm.NewStyle(pterm.BgGreen, pterm.FgWhite),
	}
	pterm.Info.MessageStyle = pterm.NewStyle(pterm.FgDefault)
	msg := fmt.Sprintf(format, args...)
	pterm.Info.Println(msg)
}

func Warning(format string, args ...any) {
	pterm.Warning.Prefix = pterm.Prefix{
		Text:  warningPrefixText,
		Style: pterm.NewStyle(pterm.BgYellow, pterm.FgBlack),
	}
	msg := fmt.Sprintf(format, args...)
	pterm.Warning.Println(msg)
}

func Debug(format string, args ...any) {
	pterm.Debug.Prefix = pterm.Prefix{
		Text:  debugPrefixText,
		Style: pterm.NewStyle(pterm.BgLightRed, pterm.FgBlack),
	}
	msg := fmt.Sprintf(format, args...)
	pterm.Debug.Println(msg)
}

func Error(format string, args ...any) {
	pterm.Error.Prefix = pterm.Prefix{
		Text:  errorPrefixText,
		Style: pterm.NewStyle(pterm.BgRed, pterm.FgLightWhite),
	}
	msg := fmt.Sprintf(format, args...)
	pterm.Error.Println(msg)
}

func Execute(msg string) {
	if !executeEnabled {
		return
	}
	pterm.Info.Prefix = pterm.Prefix{
		Text:  executePrefixText,
		Style: pterm.NewStyle(pterm.BgGreen, pterm.FgWhite),
	}
	pterm.Info.MessageStyle = pterm.NewStyle(pterm.BgDefault, pterm.FgGray)
	pterm.Info.Println(msg)
}
