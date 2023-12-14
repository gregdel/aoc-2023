package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 12)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("21", "7939", "525152", "850504257483930").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := parseInput(r, part)
	return strconv.Itoa(result), nil
}

var solvableCache = map[string]bool{}

func solvable(input []byte, hint int) bool {
	k := fmt.Sprintf("%s|%d", input, hint)
	v, ok := solvableCache[k]
	if ok {
		return v
	}

	if len(input) < hint {
		return false
	}

	for _, b := range input {
		if b == '.' {
			solvableCache[k] = false
			return false
		}

		hint--
		if hint == 0 {
			solvableCache[k] = true
			return true
		}
	}

	solvableCache[k] = false
	return false
}

var solveCache = map[string]int{}

func solve(input []byte, hints []int) int {
	k := fmt.Sprintf("%s|%+v", input, hints)
	v, ok := solveCache[k]
	if ok {
		return v
	}

	if len(input) == 0 {
		if len(hints) == 0 {
			solveCache[k] = 1
			return 1
		}
		solveCache[k] = 0
		return 0
	}

	// No more hints means not more #
	if len(hints) == 0 {
		for _, b := range input {
			if b == '#' {
				solveCache[k] = 0
				return 0
			}
		}
		solveCache[k] = 1
		return 1
	}

	minWidth := len(hints) - 1
	for _, h := range hints {
		minWidth += h
	}

	if len(input) < minWidth {
		solveCache[k] = 0
		return 0
	}

	if len(input) == 1 && len(hints) == 1 {
		if !solvable(input, hints[0]) {
			solveCache[k] = 0
			return 0
		}
		solveCache[k] = 1
		return 1
	}

	current := input[0]
	if current == '.' {
		return solve(input[1:], hints)
	}

	ret := 0
	if current == '?' {
		// Do nothing
		ret += solve(input[1:], hints)
	}

	// Try to solve the current hint

	if len(input) >= hints[0] && solvable(input[0:hints[0]], hints[0]) {
		// current hint
		ch := hints[0]
		if len(input) == ch {
			// No more next block
			ret++
		} else if len(input) >= ch+1 {
			if input[ch] == '#' {
				solveCache[k] = ret
				return ret
			}

			// next hint
			nh := []int{}
			if len(hints) >= 2 {
				nh = hints[1:]
			}

			// It it can be solved, move on to the next block + 1
			ret += solve(input[ch+1:], nh)
		}
	}

	solveCache[k] = ret
	return ret
}

func parseInput(r io.Reader, part int) int {
	scanner := bufio.NewScanner(r)
	result := 0

	for scanner.Scan() {
		n, h, ok := strings.Cut(scanner.Text(), " ")
		if !ok {
			panic("WUT ?")
		}

		h = strings.ReplaceAll(h, ",", " ")
		hints := aoc.IntsFromString(h)
		if part == 2 {
			nh := []int{}
			ns := make([]string, 5)
			for i := 0; i < 5; i++ {
				nh = append(nh, hints...)
				ns[i] = n
			}
			n = strings.Join(ns, "?")
			hints = nh
		}

		result += solve([]byte(n), hints)
	}

	return result
}
