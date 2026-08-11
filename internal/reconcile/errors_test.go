package reconcile

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// makeErrors builds n distinct errors.
func makeErrors(n int) []error {
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		errs[i] = fmt.Errorf("scan problem %d", i+1)
	}
	return errs
}

func TestReportScanErrors(t *testing.T) {
	tests := []struct {
		name          string
		count         int
		wantOutput    bool
		wantListed    int
		wantRemainder int
	}{
		{name: "no errors writes nothing", count: 0, wantOutput: false},
		{name: "three errors are all listed", count: 3, wantOutput: true, wantListed: 3},
		{name: "exactly five errors are all listed", count: 5, wantOutput: true, wantListed: 5},
		{name: "six errors list five and note one more", count: 6, wantOutput: true, wantListed: 5, wantRemainder: 1},
		{name: "ten errors list five and note five more", count: 10, wantOutput: true, wantListed: 5, wantRemainder: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			ReportScanErrors(&buf, makeErrors(tt.count))
			out := buf.String()

			if !tt.wantOutput {
				if out != "" {
					t.Errorf("Expected no output, got %q", out)
				}
				return
			}

			wantHeader := fmt.Sprintf("Warning: %d errors occurred during scanning:", tt.count)
			if !strings.Contains(out, wantHeader) {
				t.Errorf("Expected header %q in output:\n%s", wantHeader, out)
			}

			listed := strings.Count(out, "  - scan problem")
			if listed != tt.wantListed {
				t.Errorf("Listed %d errors, want %d. Output:\n%s", listed, tt.wantListed, out)
			}

			remainderLine := fmt.Sprintf("  ... and %d more errors", tt.wantRemainder)
			if tt.wantRemainder > 0 {
				if !strings.Contains(out, remainderLine) {
					t.Errorf("Expected %q in output:\n%s", remainderLine, out)
				}
			} else if strings.Contains(out, "... and") {
				t.Errorf("Did not expect a remainder line, got:\n%s", out)
			}
		})
	}
}

func TestReportScanErrors_NilSlice(t *testing.T) {
	var buf bytes.Buffer
	ReportScanErrors(&buf, nil)
	if buf.String() != "" {
		t.Errorf("Expected no output for a nil slice, got %q", buf.String())
	}
}

// TestReportScanErrors_ExactFormat pins the output shape both commands share,
// so a change to one cannot silently diverge from the other.
func TestReportScanErrors_ExactFormat(t *testing.T) {
	var buf bytes.Buffer
	ReportScanErrors(&buf, makeErrors(7))

	want := strings.Join([]string{
		"Warning: 7 errors occurred during scanning:",
		"  - scan problem 1",
		"  - scan problem 2",
		"  - scan problem 3",
		"  - scan problem 4",
		"  - scan problem 5",
		"  ... and 2 more errors",
		"",
	}, "\n")

	if buf.String() != want {
		t.Errorf("Output mismatch.\ngot:\n%s\nwant:\n%s", buf.String(), want)
	}
}
