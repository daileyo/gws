package git

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"golang.org/x/term"
)

// Progress manages a spinner with progress count displayed on stderr.
// It is safe for concurrent use via atomic operations.
type Progress struct {
	total     int
	completed atomic.Int32
	isTTY     bool
	stop      chan struct{}
	done      chan struct{}
}

// NewProgress creates a new Progress tracker. The spinner only renders
// if stderr is a TTY.
func NewProgress(total int) *Progress {
	return &Progress{
		total: total,
		isTTY: term.IsTerminal(int(os.Stderr.Fd())),
		stop:  make(chan struct{}),
		done:  make(chan struct{}),
	}
}

// Start begins rendering the spinner in a background goroutine.
// The spinner only appears after a 200ms delay to avoid flicker for fast operations.
func (p *Progress) Start() {
	if !p.isTTY {
		close(p.done)
		return
	}

	go func() {
		defer close(p.done)

		frames := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
		frameIdx := 0

		// Wait 200ms before showing spinner
		select {
		case <-time.After(200 * time.Millisecond):
		case <-p.stop:
			return
		}

		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-p.stop:
				// Clear spinner line
				fmt.Fprintf(os.Stderr, "\r\033[K")
				return
			case <-ticker.C:
				c := p.completed.Load()
				fmt.Fprintf(os.Stderr, "\r%c Checking repos... %d/%d", frames[frameIdx], c, p.total)
				frameIdx = (frameIdx + 1) % len(frames)
			}
		}
	}()
}

// Increment marks one more repo as completed. Safe for concurrent use.
func (p *Progress) Increment() {
	p.completed.Add(1)
}

// Stop halts the spinner and clears the line. Blocks until cleanup is done.
func (p *Progress) Stop() {
	select {
	case <-p.stop:
		// Already stopped
	default:
		close(p.stop)
	}
	<-p.done
}
