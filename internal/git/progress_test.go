package git

import (
	"sync"
	"testing"
)

func TestProgress_Increment(t *testing.T) {
	t.Run("concurrent increments are safe", func(t *testing.T) {
		p := NewProgress(100)
		// Don't start the spinner (no TTY in tests)

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.Increment()
			}()
		}
		wg.Wait()

		if got := p.completed.Load(); got != 100 {
			t.Errorf("expected 100, got %d", got)
		}
	})
}

func TestNewProgress_UsesDefaultLabel(t *testing.T) {
	p := NewProgress(10)
	if p.label != DefaultProgressLabel {
		t.Errorf("label = %q, want %q", p.label, DefaultProgressLabel)
	}
	if p.total != 10 {
		t.Errorf("total = %d, want 10", p.total)
	}
}

func TestNewProgressWithLabel(t *testing.T) {
	tests := []struct {
		name      string
		label     string
		wantLabel string
	}{
		{"custom label is used", "Reconciling workspace...", "Reconciling workspace..."},
		{"empty label falls back to the default", "", DefaultProgressLabel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProgressWithLabel(5, tt.label)
			if p.label != tt.wantLabel {
				t.Errorf("label = %q, want %q", p.label, tt.wantLabel)
			}
			// The lifecycle must stay intact regardless of label.
			p.Start()
			p.Increment()
			p.Stop()
			if got := p.completed.Load(); got != 1 {
				t.Errorf("completed = %d, want 1", got)
			}
		})
	}
}

func TestProgress_StopWithoutStart(t *testing.T) {
	// Stop should not panic if Start was never called
	p := &Progress{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	close(p.done) // simulate non-TTY (done is closed immediately)
	p.Stop()      // should not block or panic
}

func TestProgress_NonTTY(t *testing.T) {
	// In test environment, stderr is not a TTY
	p := NewProgress(10)
	p.Start() // should close done immediately since not a TTY
	p.Stop()  // should not block
}

func TestProgress_StartStop(t *testing.T) {
	// Verify Start/Stop lifecycle doesn't hang
	p := NewProgress(5)
	p.Start()
	p.Increment()
	p.Increment()
	p.Stop()

	// Double stop should be safe
	p.Stop()
}
