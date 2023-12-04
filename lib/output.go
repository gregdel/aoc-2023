package aoc

import (
	"fmt"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/text"
)

// RunResult represents the result of a run
type RunResult struct {
	day, part        int
	test, ok         bool
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
	sign := text.Colors{text.FgGreen}.Sprint("✓")
	msg := ""
	if !success {
		sign = text.Colors{text.FgRed}.Sprint("✗")
		msg = fmt.Sprintf("expected %q, got %q", r.expected, r.output)
	}

	var out strings.Builder
	fmt.Fprintf(&out, "%s day:%d part:%d test:%t duration:%s",
		sign, r.day, r.part, r.test, r.stop.Sub(r.start))
	fmt.Fprintf(&out, " %s", msg)

	fmt.Println(out.String())
}
