package aoc

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

// Custom errors
var (
	ErrMissingDay = errors.New("missing day")
)

// Challenges list
var Challenges map[int]Challenge = map[int]Challenge{}

// Register registers a challenge
func Register(c Challenge) {
	Challenges[c.Day()] = c
}

// AllDays returns a list of all the days regitered
func AllDays() []int {
	values := []int{}
	for d := range Challenges {
		values = append(values, d)
	}
	sort.Ints(values)
	return values
}

// Challenge represents a challenge
type Challenge interface {
	Day() int
	Solve(r io.Reader, part int) (string, error)
	Expect(part int, test bool) string
}

// Open opens the input for a given day
func Open(day, part int, test bool) (io.ReadCloser, error) {
	filename := "input"
	if test {
		filename = "input-test-" + strconv.Itoa(part)
	}

	path := filepath.Join("day"+strconv.Itoa(day), filename)
	return os.Open(path)
}

// Run run a the challenge
func Run(day, part int, test bool) (*RunResult, error) {
	challenge, ok := Challenges[day]
	if !ok {
		return nil, ErrMissingDay
	}

	input, err := Open(day, part, test)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	result := newRunResult(day, part, test)
	result.expected = challenge.Expect(part, test)
	result.start = time.Now()
	result.output, err = challenge.Solve(input, part)
	result.stop = time.Now()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MustGet returns v as is. It panics if err is non-nil.
func MustGet[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
