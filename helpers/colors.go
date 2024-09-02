package helpers

import "github.com/pterm/pterm"

func NormalText(s string) string {
	color := pterm.NewStyle(pterm.FgDefault, pterm.BgDefault)
	reset := pterm.NewStyle(pterm.Reset)

	// Wrapping in reset and bold is necessary
	// to maintain proper spacing in table formatting
	return reset.Sprint(color.Sprint(s))
}

func BoldText(s string) string {
	color := pterm.NewStyle(pterm.FgDefault, pterm.BgDefault)
	bold := pterm.NewStyle(pterm.Bold)
	return bold.Sprint(color.Sprint(s))
}

func GreenText(s string) string {
	color := pterm.NewStyle(pterm.FgGreen, pterm.BgDefault)
	reset := pterm.NewStyle(pterm.Reset)
	return reset.Sprint(color.Sprint(s))
}
