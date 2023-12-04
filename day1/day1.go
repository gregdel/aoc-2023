package day1

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"unicode"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day1{})
}

var words = []string{
	"zero", "one", "two", "three", "four",
	"five", "six", "seven", "eight", "nine",
}

type day1 struct{}

func (d *day1) Day() int {
	return 1
}

func (d *day1) Solve(r io.Reader, part int) (string, error) {
	return solve(r, part == 2)
}

func (d *day1) Expect(part int, test bool) string {
	return aoc.NewResult("142", "53921", "281", "54676").Expect(part, test)
}

func transformLine(input string) string {
	line := input
	for i := 1; i < len(words); i++ {
		o := words[i]
		n := words[i] + strconv.Itoa(i) + words[i]
		line = strings.ReplaceAll(line, o, n)
	}

	return line
}

func solve(reader io.Reader, transform bool) (string, error) {
	result := 0
	var err error

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var first, last int
		var foundFirst bool

		line := scanner.Text()
		if transform {
			line = transformLine(line)
		}

		for _, r := range line {
			if !unicode.IsDigit(r) {
				continue
			}

			if !foundFirst {
				first, err = strconv.Atoi(string(r))
				if err != nil {
					return "", err
				}
				last = first
				foundFirst = true
			}

			last, err = strconv.Atoi(string(r))
			if err != nil {
				return "", err
			}
		}

		value := first*10 + last
		result += value
	}

	return strconv.Itoa(result), nil
}
