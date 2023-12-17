package day

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/fatih/color"
	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 16)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("46", "6978", "51", "7315").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := parseInput(r, part)
	return strconv.Itoa(result), nil
}

type beam struct {
	commingFrom aoc.Direction
	current     *aoc.Point
}

func newBeam(p *aoc.Point, d aoc.Direction) *beam {
	return &beam{commingFrom: d, current: p}
}

type area struct {
	*aoc.Map2D
	explored map[*aoc.Point]map[aoc.Direction]bool
	beams    []*beam
}

// String implements the stringer interface
func (a *area) String() string {
	exploredColor := color.New(color.FgGreen).SprintFunc()

	var out strings.Builder
	for y := 0; y < len(a.Points); y++ {
		for x := 0; x < len(a.Points[y]); x++ {
			p := a.Points[y][x]
			_, ok := a.explored[p]
			if ok {
				out.WriteString(exploredColor(string(p.C)))
			} else {
				out.WriteRune(p.C)
			}
		}
		out.WriteRune('\n')
	}

	return out.String()
}

func newArea() *area {
	return &area{
		beams:    []*beam{},
		explored: map[*aoc.Point]map[aoc.Direction]bool{},
		Map2D:    aoc.NewMap2D(),
	}
}

func (a *area) move(b *beam, d aoc.Direction) {
	next := a.Next(d, b.current)
	if next == nil {
		return
	}

	a.beams = append(a.beams, newBeam(next, aoc.OppositeDirection(d)))
}

func (a *area) moveAll() bool {
	currentBeams := a.beams
	hasNewPoint := false
	a.beams = []*beam{}
	for _, beam := range currentBeams {
		cf := beam.commingFrom

		_, ok := a.explored[beam.current]
		if !ok {
			a.explored[beam.current] = map[aoc.Direction]bool{}
		}
		_, ok = a.explored[beam.current][cf]
		if ok {
			continue
		}

		hasNewPoint = true
		a.explored[beam.current][cf] = true

		switch beam.current.C {
		case '.':
			a.move(beam, aoc.OppositeDirection(beam.commingFrom))
		case '|':
			if cf == aoc.DirectionLeft || cf == aoc.DirectionRight {
				a.move(beam, aoc.DirectionUp)
				a.move(beam, aoc.DirectionDown)
			} else {
				a.move(beam, aoc.OppositeDirection(beam.commingFrom))
			}
		case '-':
			if cf == aoc.DirectionLeft || cf == aoc.DirectionRight {
				a.move(beam, aoc.OppositeDirection(beam.commingFrom))
			} else {
				a.move(beam, aoc.DirectionRight)
				a.move(beam, aoc.DirectionLeft)
			}
		case '/':
			var nd aoc.Direction
			switch cf {
			case aoc.DirectionDown:
				nd = aoc.DirectionRight
			case aoc.DirectionUp:
				nd = aoc.DirectionLeft
			case aoc.DirectionRight:
				nd = aoc.DirectionDown
			case aoc.DirectionLeft:
				nd = aoc.DirectionUp
			}
			a.move(beam, nd)
		case '\\':
			var nd aoc.Direction
			switch cf {
			case aoc.DirectionDown:
				nd = aoc.DirectionLeft
			case aoc.DirectionUp:
				nd = aoc.DirectionRight
			case aoc.DirectionRight:
				nd = aoc.DirectionUp
			case aoc.DirectionLeft:
				nd = aoc.DirectionDown
			}
			a.move(beam, nd)
		}
	}

	return hasNewPoint
}

func (a *area) explore(x, y int, d aoc.Direction) int {
	a.explored = map[*aoc.Point]map[aoc.Direction]bool{}
	a.beams = []*beam{newBeam(a.Points[y][x], d)}

	for a.moveAll() {
	}

	result := 0
	a.ForAllPoints(func(p *aoc.Point) {
		_, ok := a.explored[p]
		if ok {
			result++
		}
	})
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseInput(r io.Reader, part int) int {
	scanner := bufio.NewScanner(r)

	area := newArea()
	for scanner.Scan() {
		area.AddPointsFromLine(scanner.Text())
	}

	if part == 1 {
		return area.explore(0, 0, aoc.DirectionLeft)
	}

	result := 0

	// Top row down and bottom row up
	maxY := len(area.Points) - 1
	for x := 0; x < len(area.Points[0]); x++ {
		v := area.explore(x, 0, aoc.DirectionUp)
		result = max(v, result)
		v = area.explore(x, maxY, aoc.DirectionDown)
		result = max(v, result)
	}

	// Top row down and bottom row up
	maxX := len(area.Points[0]) - 1
	for y := 0; y < len(area.Points); y++ {
		v := area.explore(0, y, aoc.DirectionLeft)
		result = max(v, result)
		v = area.explore(maxX, y, aoc.DirectionRight)
		result = max(v, result)
	}

	return result
}
