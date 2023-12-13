package day

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 13)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("405", "30518", "400", "36735").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := parseInput(r, part)
	return strconv.Itoa(result), nil
}

type area struct {
	points                   [][]*point
	reflectionX, reflectionY int
}

func (a *area) String() string {
	var out strings.Builder
	for y := 0; y < len(a.points); y++ {
		for x := 0; x < len(a.points[0]); x++ {
			c := '.'
			if a.points[y][x].isRock {
				c = '#'
			}
			out.WriteRune(c)
		}
		out.WriteRune('\n')
	}

	return out.String()
}

func newArea() *area {
	return &area{
		points:      [][]*point{},
		reflectionX: -1,
		reflectionY: -1,
	}
}

func (a *area) reflectionRightOfX(x int) bool {
	lineLen := len(a.points[0])
	reflection := true
	for y := 0; y < len(a.points); y++ {
		offset := 1
		for {
			x1, x2 := x-offset+1, x+offset
			if x1 < 0 || x2 > lineLen-1 {
				break
			}

			p1, p2 := a.points[y][x1], a.points[y][x2]
			if p1.isRock != p2.isRock {
				return false
			}

			offset++
		}
	}

	return reflection
}

func (a *area) reflectionUnderY(y int) bool {
	colLen := len(a.points)
	reflection := true

	offset := 1
	for {
		y1, y2 := y-offset+1, y+offset
		if y1 < 0 || y2 > colLen-1 {
			break
		}

		for x := 0; x < len(a.points[0]); x++ {
			p1, p2 := a.points[y1][x], a.points[y2][x]
			if p1.isRock != p2.isRock {
				return false
			}
		}

		offset++
	}

	return reflection
}

func (a *area) findReflection(save bool) int {
	for x := 0; x < len(a.points[0])-1; x++ {
		if x == a.reflectionX {
			continue
		}

		if a.reflectionRightOfX(x) {
			if save {
				a.reflectionX = x
			}

			return x + 1
		}
	}

	for y := 0; y < len(a.points)-1; y++ {
		if y == a.reflectionY {
			continue
		}

		if a.reflectionUnderY(y) {
			if save {
				a.reflectionY = y
			}
			return (y + 1) * 100
		}
	}
	return 0
}

func (a *area) findSmudge() int {
	a.findReflection(true)

	for y := 0; y < len(a.points); y++ {
		for x := 0; x < len(a.points[0]); x++ {
			point := a.points[y][x]
			point.isRock = !point.isRock

			v := a.findReflection(false)
			point.isRock = !point.isRock

			if v != 0 {
				return v
			}
		}
	}

	return 0
}

type point struct {
	x, y   int
	isRock bool
}

func newPoint(x, y int, isRock bool) *point {
	return &point{
		x: x, y: y, isRock: isRock,
	}
}

func parseInput(r io.Reader, part int) int {
	scanner := bufio.NewScanner(r)
	result := 0

	area := newArea()
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			if part == 1 {
				result += area.findReflection(false)
			} else {
				result += area.findSmudge()
			}

			area = newArea()
			y = 0
			continue
		}

		points := []*point{}
		for x, c := range line {
			points = append(points, newPoint(x, y, c == '#'))
		}
		area.points = append(area.points, points)
		y++
	}

	if part == 1 {
		result += area.findReflection(false)
	} else {
		result += area.findSmudge()
	}
	return result
}
