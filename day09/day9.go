package day

import (
	"bufio"
	"io"
	"strconv"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 9)
}

type suite struct {
	values     []int
	prev, next *suite
}

func (s *suite) predictNext() {
	s.createSubSuites()
	last := s.getLastSubSuite()
	last.addToLast(0)
}

func (s *suite) predictFirst() {
	s.createSubSuites()
	last := s.getLastSubSuite()
	last.addToFirst(0)
}

func (s *suite) getLastSubSuite() *suite {
	last := s.next
	for last.next != nil {
		last = last.next
	}
	return last
}

func (s *suite) lastValue() int {
	return s.values[len(s.values)-1]
}

func (s *suite) firstValue() int {
	return s.values[0]
}

func (s *suite) addToFirst(value int) {
	s.values = s.values[0:1]
	firstValue := s.firstValue()
	newValue := firstValue - value
	s.values[0] = newValue

	if s.prev != nil {
		s.prev.addToFirst(newValue)
	}
}

func (s *suite) addToLast(value int) {
	newValue := s.lastValue() + value
	s.values = append(s.values, newValue)

	if s.prev != nil {
		s.prev.addToLast(newValue)
	}
}

func (s *suite) createSubSuites() {
	if len(s.values) == 0 {
		panic("no more values")
	}

	s.next = &suite{
		values: make([]int, len(s.values)-1),
		prev:   s,
	}

	allZeros := true
	for i := 0; i < len(s.values)-1; i++ {
		v := s.values[i+1] - s.values[i]
		s.next.values[i] = v
		if allZeros && v != 0 {
			allZeros = false
		}
	}

	if !allZeros {
		s.next.createSubSuites()
	}
}

func newSuite(input string) *suite {
	return &suite{
		values: aoc.IntsFromString(input),
	}
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("114", "2043183816", "2", "1118").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := 0
	suites := parseInput(r, part)
	for _, suite := range suites {
		if part == 1 {
			suite.predictNext()
			result += suite.lastValue()
		} else {
			suite.predictFirst()
			result += suite.firstValue()
		}
	}

	return strconv.Itoa(result), nil
}

func parseInput(r io.Reader, part int) []*suite {
	scanner := bufio.NewScanner(r)

	suites := []*suite{}
	for scanner.Scan() {
		line := scanner.Text()
		suites = append(suites, newSuite(line))
	}

	return suites
}
