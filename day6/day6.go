package day

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{})
}

type race struct {
	time, distance int
}

func (r *race) waysToWin() int {
	wins := 0

	for i := 0; i < r.time; i++ {
		pressed := i
		remainingTime := r.time - i
		distance := pressed * remainingTime
		if distance > r.distance {
			wins++
		}
	}

	return wins
}

type day struct{}

func (d *day) Day() int {
	return 6
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := 1
	races := d.parseInput(r, part)
	for _, r := range races {
		result *= r.waysToWin()
	}
	return strconv.Itoa(result), nil
}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("288", "316800", "71503", "45647654").Expect(part, test)
}

func parseNumbers(input string) []int {
	output := []int{}
	for _, str := range strings.Fields(input) {
		v := aoc.MustGet(strconv.Atoi(str))
		output = append(output, v)
	}

	return output
}

func (d *day) parseInput(r io.Reader, part int) []race {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	parts := strings.Split(scanner.Text(), ":")
	if part == 2 {
		parts[1] = strings.ReplaceAll(parts[1], " ", "")
	}
	times := parseNumbers(parts[1])
	scanner.Scan()
	parts = strings.Split(scanner.Text(), ":")
	if part == 2 {
		parts[1] = strings.ReplaceAll(parts[1], " ", "")
	}
	distances := parseNumbers(parts[1])

	races := make([]race, len(times))
	for i := 0; i < len(times); i++ {
		races[i] = race{
			time:     times[i],
			distance: distances[i],
		}
	}
	return races
}
