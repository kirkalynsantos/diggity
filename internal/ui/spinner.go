package ui

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

var disabled bool

// InitSpinner generates simple spinner
func InitSpinner(text string) *progressbar.ProgressBar {
	if disabled {
		return nil
	}
	spinner := progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetDescription(text),
		progressbar.OptionClearOnFinish(),
	)
	return spinner
}

// RunSpinner starts a spinner
func RunSpinner(spinner *progressbar.ProgressBar) {
	if disabled {
		return
	}
	for {
		spinner.Add(1)
		time.Sleep(5 * time.Millisecond)
	}
}

// DoneSpinner stops and closes a spiner
func DoneSpinner(spinner *progressbar.ProgressBar) {
	if disabled {
		return
	}
	spinner.Finish()
	spinner.Close()
}

// Disable spinner
func Disable() {
	disabled = true
}
