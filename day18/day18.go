package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 18)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("62", "33491", "952408144115", "87716969654406").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := parseInput(r, part)
	return strconv.Itoa(result), nil
}

func parseInput(r io.Reader, part int) int {
	scanner := bufio.NewScanner(r)

	border := 0
	area := 0
	x, y := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		sx, sy := x, y
		var d rune
		var cnt, color int
		fmt.Sscanf(line, "%c %d (#%x)", &d, &cnt, &color)

		if part == 2 {
			switch color & 0x0F {
			case 0:
				d = 'R'
			case 1:
				d = 'D'
			case 2:
				d = 'L'
			case 3:
				d = 'U'
			}
			cnt = color >> 4
		}

		switch d {
		case 'U':
			y -= cnt
		case 'D':
			y += cnt
		case 'R':
			x += cnt
		case 'L':
			x -= cnt
		}
		border += cnt
		area += (y + sy) * (x - sx)
	}

	if area < 0 {
		area = -area
	}
	area /= 2

	inside := area - (border / 2) + 1
	result := border + inside
	return result
}
