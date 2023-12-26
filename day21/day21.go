package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 21)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("42", "3585", "1594", "597102953699891").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	area := parseInput(r, part)
	result := 0
	if part == 1 {
		result = area.run(part, 64)
	} else {
		result = area.run(part, 0)
	}
	return strconv.Itoa(result), nil
}

type point struct {
	*aoc.Point
	mapOffsetX, mapOffsetY int
}

func (p point) String() string {
	return fmt.Sprintf("[%2d;%2d] {%2d;%2d}", p.X, p.Y, p.mapOffsetX, p.mapOffsetY)
}

type area struct {
	*aoc.Map2D
	start point
}

func (a *area) next(d aoc.Direction, p point) point {
	mapOffsetX := p.mapOffsetX
	mapOffsetY := p.mapOffsetY

	next := a.Next(d, p.Point)
	if next == nil {
		switch d {
		case aoc.DirectionUp:
			next = a.Points[len(a.Points)-1][p.X]
			mapOffsetY--
		case aoc.DirectionDown:
			next = a.Points[0][p.X]
			mapOffsetY++
		case aoc.DirectionLeft:
			next = a.Points[p.Y][len(a.Points[0])-1]
			mapOffsetX--
		case aoc.DirectionRight:
			next = a.Points[p.Y][0]
			mapOffsetX++
		}
	}

	return point{
		Point:      next,
		mapOffsetX: mapOffsetX,
		mapOffsetY: mapOffsetY,
	}
}

func newArea() *area {
	return &area{
		Map2D: aoc.NewMap2D(),
	}
}

func (a *area) run(part, steps int) int {
	toExplore := map[point]struct{}{
		a.start: struct{}{},
	}

	inputSize := len(a.Points)
	realInput := inputSize == 131
	if part == 2 {
		if !realInput {
			steps = 50
		} else {
			steps = inputSize*2 + 65 + 1
		}
	}

	// For part 2 we need to find the steps
	// 26501365 = 202300 * 131 + 65 (131 is the input size)
	// Assuming Y = 131*n + 65, f(202300) = 26501365
	// f(0) = 3585
	// f(1) = 32657
	// f(2) = 90909
	// f(n) = XXXXX
	// f(202300) = expected steps
	// We're gonna store values for x={0,1,2}, to find x=202300
	keyValues := []int{}

	tmp := map[point]struct{}{}
	for i := 0; i < steps; i++ {
		for p := range toExplore {
			for _, d := range aoc.AllDirection {
				next := a.next(d, p)
				if part == 1 && (next.mapOffsetX != 0 || next.mapOffsetY != 0) {
					continue
				}
				if next.C == '#' {
					continue
				}
				tmp[next] = struct{}{}
			}
		}

		if part == 2 && realInput && ((i-65)%131) == 0 {
			keyValues = append(keyValues, len(toExplore))
		}

		toExplore = tmp
		tmp = make(map[point]struct{}, len(toExplore)*2)
	}

	if part == 2 && realInput {
		v0, v1, v2 := keyValues[0], keyValues[1], keyValues[2]
		// f(x) = ax² + bx + c
		// f(0) = v0
		// f(1) = v1
		// f(2) = v2
		c := v0
		a := (v2 + v0 - 2*v1) / 2
		b := v1 - v0 - a
		x := 202300
		return a*x*x + b*x + c
	}

	return len(toExplore)
}

func parseInput(r io.Reader, part int) *area {
	scanner := bufio.NewScanner(r)

	area := newArea()
	for scanner.Scan() {
		area.AddPointsFromLine(scanner.Text())
	}

	area.ForAllPoints(func(p *aoc.Point) {
		if p.C == 'S' {
			area.start = point{Point: p}
			p.C = '.'
		}
	})

	return area
}
