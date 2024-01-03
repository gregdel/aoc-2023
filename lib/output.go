package aoc

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	redColor   = color.New(color.FgRed).SprintFunc()
	greenColor = color.New(color.FgGreen).SprintFunc()

	successIcon = greenColor("✓")
	failureIcon = redColor("✗")
)

// RunResult represents the result of a run
type RunResult struct {
	day, part        int
	test             bool
	output, expected string
	start, stop      time.Time
}

func newRunResult(day, part int, test bool) *RunResult {
	return &RunResult{
		day:  day,
		part: part,
		test: test,
	}
}

// Show shows the result of the test
func (r *RunResult) Show() {
	success := r.output == r.expected

	var icon, msg string
	if !success {
		icon = failureIcon
		msg = fmt.Sprintf("expected %q, got %q", r.expected, r.output)
	} else {
		icon = successIcon
	}

	var out strings.Builder
	fmt.Fprintf(&out, "%s day:%d part:%d test:%t duration:%s",
		icon, r.day, r.part, r.test, r.stop.Sub(r.start))
	fmt.Fprintf(&out, " %s", msg)

	fmt.Println(out.String())
}
