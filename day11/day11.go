package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fatih/color"
	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 11)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("374", "10422930", "82000210", "699909023130").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	universe := parseInput(r, part)
	result := universe.computeDistances()
	return strconv.Itoa(result), nil
}

type point struct {
	x, y     int
	c, p     rune
	isGalaxy bool
}

func abs(x int) int {
	if x < 0 {
		return x * -1
	}

	return x
}

func distance(a, b *point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func between(x1, x2, x int) bool {
	if x1 > x2 {
		x2, x1 = x1, x2
	}
	return x > x1 && x < x2
}

func newPoint(x, y int, c rune) *point {
	isGalaxy := false
	p := '·'
	if c == '#' {
		isGalaxy = true
		p = ''
	}

	return &point{
		x: x, y: y,
		c: c, p: p,
		isGalaxy: isGalaxy,
	}
}

type universe struct {
	galaxies       []*point
	data           [][]*point
	emptyX, emptyY []int
	factor         int
}

func (u *universe) distance(a, b *point) int {
	distance := distance(a, b)
	for _, x := range u.emptyX {
		if between(a.x, b.x, x) {
			distance += u.factor - 1
		}
	}

	for _, y := range u.emptyY {
		if between(a.y, b.y, y) {
			distance += u.factor - 1
		}
	}

	return distance
}

func (u *universe) computeDistances() int {
	result := 0

	for i := 0; i < len(u.galaxies)-1; i++ {
		for j := i + 1; j < len(u.galaxies); j++ {
			result += u.distance(u.galaxies[i], u.galaxies[j])
		}
	}

	return result
}

func (u *universe) border(out *strings.Builder, start bool) {
	opening, closing := '┌', '┐'
	if !start {
		opening, closing = '└', '┘'
	}

	out.WriteRune(opening)
	for x := 0; x < len(u.data[0]); x++ {
		out.WriteRune('─')
	}
	out.WriteRune(closing)
	out.WriteRune('\n')
}

func (u *universe) String() string {
	galaxyColor := color.New(color.FgGreen).SprintFunc()
	var out strings.Builder
	u.border(&out, true)
	for y := 0; y < len(u.data); y++ {
		out.WriteRune('│')
		for x := 0; x < len(u.data[0]); x++ {
			pt := u.data[y][x]
			if pt.isGalaxy {
				out.WriteString(galaxyColor(string(pt.p)))
			} else {
				out.WriteRune(pt.p)
			}
		}
		out.WriteRune('│')
		out.WriteRune('\n')
	}
	u.border(&out, false)

	out.WriteString(fmt.Sprintf("emptyX:%+v\n", u.emptyX))
	out.WriteString(fmt.Sprintf("emptyY:%+v\n", u.emptyY))

	return out.String()
}

func newUniverse() *universe {
	return &universe{
		data:     [][]*point{},
		galaxies: []*point{},
		emptyX:   []int{},
		emptyY:   []int{},
	}
}

func parseInput(r io.Reader, part int) *universe {
	scanner := bufio.NewScanner(r)

	universe := newUniverse()
	universe.factor = 2
	if part == 2 {
		universe.factor = 1000000
	}

	y := 0
	for scanner.Scan() {
		allDots := true
		points := []*point{}
		for x, c := range scanner.Text() {
			point := newPoint(x, y, c)
			points = append(points, point)
			if point.isGalaxy {
				universe.galaxies = append(universe.galaxies, point)
				allDots = false
			}
		}

		if allDots {
			universe.emptyY = append(universe.emptyY, y)
		}

		universe.data = append(universe.data, points)
		y++
	}

	for x := 0; x < len(universe.data); x++ {
		allDots := true
		for y := 0; y < len(universe.data); y++ {
			if universe.data[y][x].isGalaxy {
				allDots = false
			}
		}

		if allDots {
			universe.emptyX = append(universe.emptyX, x)
		}
	}

	return universe
}
