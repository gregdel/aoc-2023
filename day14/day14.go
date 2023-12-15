package day

import (
	"bufio"
	"io"
	"strconv"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 14)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("136", "113424", "64", "96003").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := parseInput(r, part)
	return strconv.Itoa(result), nil
}

type area struct {
	*aoc.Map2D
}

func newArea() *area {
	return &area{
		Map2D: aoc.NewMap2D(),
	}
}

func (a *area) weight(p *aoc.Point) int {
	if p.C != 'O' {
		return 0
	}

	return len(a.Points) - p.Y
}

func (a *area) moveRock(d aoc.Direction, p *aoc.Point) {
	if p.C != 'O' {
		return
	}

	next := p
	for {
		n := a.Next(d, next)
		if n == nil || n.C != '.' {
			break
		}
		next = n
	}

	if next == p {
		return
	}

	p.C = '.'
	next.C = 'O'
}

func (a *area) tilt(d aoc.Direction) {
	var d1, d2 aoc.Direction
	switch d {
	case aoc.DirectionUp:
		d1, d2 = aoc.DirectionDown, aoc.DirectionRight
	case aoc.DirectionDown:
		d1, d2 = aoc.DirectionUp, aoc.DirectionRight
	case aoc.DirectionLeft:
		d1, d2 = aoc.DirectionDown, aoc.DirectionRight
	case aoc.DirectionRight:
		d1, d2 = aoc.DirectionDown, aoc.DirectionLeft
	}

	a.ForAllPoints(func(p *aoc.Point) {
		a.moveRock(d, p)
	}, d1, d2)
}

func (a *area) cycle() {
	for _, d := range []aoc.Direction{
		aoc.DirectionUp,
		aoc.DirectionLeft,
		aoc.DirectionDown,
		aoc.DirectionRight,
	} {
		a.tilt(d)
	}
}

func (a *area) totalWeight() int {
	weight := 0
	a.ForAllPoints(func(p *aoc.Point) {
		weight += a.weight(p)
	})
	return weight
}

func parseInput(r io.Reader, part int) int {
	scanner := bufio.NewScanner(r)

	area := newArea()
	for scanner.Scan() {
		area.AddPointsFromLine(scanner.Text())
	}

	if part == 1 {
		area.tilt(aoc.DirectionUp)
	} else {
		var sums = map[string]int{}

		patternStart, patternEnd := "", ""
		patternStartIdx, patternEndIdx := 0, 0

		sum := area.SHA256Sum()
		sums[sum]++
		loops := 1000000000
		for i := 0; i < loops; i++ {
			area.cycle()
			sum = area.SHA256Sum()

			v, ok := sums[sum]
			if ok {
				if v == 1 && patternStart == "" {
					patternStart = sum
					patternStartIdx = i
				} else if v == 2 && patternEnd == "" {
					patternEnd = sum
					patternEndIdx = i

					patternSize := patternEndIdx - patternStartIdx
					remaining := loops - i + 1
					var skip int = remaining / patternSize
					skip *= patternSize
					i += skip
					continue
				}
			}

			sums[sum]++
		}
	}

	return area.totalWeight()
}
