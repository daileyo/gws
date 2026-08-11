package reconcile

import (
	"fmt"
	"io"
)

// maxReportedScanErrors is how many individual scan errors are listed before
// the reporter collapses the remainder into a count.
const maxReportedScanErrors = 5

// ReportScanErrors writes a summary of non-fatal scan errors. Both "gws init"
// and "gws refresh" use it so identical problems are reported identically no
// matter which command surfaced them. Nothing is written when errs is empty.
//
// Callers pass os.Stderr: these are warnings, and keeping them off stdout
// leaves piped command output clean.
func ReportScanErrors(w io.Writer, errs []error) {
	if len(errs) == 0 {
		return
	}

	fmt.Fprintf(w, "Warning: %d errors occurred during scanning:\n", len(errs))
	for i, err := range errs {
		if i >= maxReportedScanErrors {
			break
		}
		fmt.Fprintf(w, "  - %v\n", err)
	}
	if len(errs) > maxReportedScanErrors {
		fmt.Fprintf(w, "  ... and %d more errors\n", len(errs)-maxReportedScanErrors)
	}
}
