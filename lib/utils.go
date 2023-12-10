package aoc

import (
	"strconv"
	"strings"
)

// IntsFromString returns a slice of ints from a string
func IntsFromString(input string) []int {
	fields := strings.Fields(input)
	output := make([]int, len(fields))
	for i := 0; i < len(fields); i++ {
		output[i] = MustGet(strconv.Atoi(fields[i]))
	}
	return output
}
