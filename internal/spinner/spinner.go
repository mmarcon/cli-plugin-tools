package spinner

import (
	"time"

	"github.com/schollz/progressbar/v3"
)

type Progress struct {
	Bar *progressbar.ProgressBar
}

// New creates a new Progress instance and starts the spinner
func New(description string) *Progress {
	p := &Progress{
		Bar: progressbar.NewOptions(-1, progressbar.OptionSpinnerType(14), progressbar.OptionSetDescription(description)),
	}
	go func() {
		for {
			p.Bar.Add(1)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return p
}

// Stop stops the spinner
func (p *Progress) Stop() {
	p.Bar.Clear()
	p.Bar.Close()
	// Clear the line completely
	// Probably a hack, but it works
	print("\r\033[K")
}
