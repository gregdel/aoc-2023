package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 17)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("102", "953", "94", "1180").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := parseInput(r, part)
	return strconv.Itoa(result), nil
}

func heat(p *aoc.Point) int {
	return int(p.C) - int('0')
}

type area struct {
	*aoc.Map2D
	start, target *aoc.Point
	heats         map[state]int
}

func newArea() *area {
	return &area{
		Map2D: aoc.NewMap2D(),
		heats: map[state]int{},
	}
}

type state struct {
	p *aoc.Point
	d aoc.Direction
	s int // steps
}

func (s state) String() string {
	return fmt.Sprintf("%s d:%s s:%d", s.p, s.d, s.s)
}

func (a *area) toTarget(p *aoc.Point) int {
	return aoc.ManhattanDistance(p, a.target)
}

func (a *area) explore(min, max int) int {
	explored := aoc.NewSet[state]()

	queue := aoc.NewPriorityQueue[state]()
	queue.Push(state{p: a.start, d: aoc.DirectionRight}, a.toTarget(a.start))
	queue.Push(state{p: a.start, d: aoc.DirectionDown}, a.toTarget(a.start))

	for queue.Len() != 0 {
		s, prio := queue.Pop()

		if explored.Has(s) {
			continue
		}

		h := prio - a.toTarget(s.p)
		if s.p == a.target {
			return h
		}

		opposite := aoc.OppositeDirection(s.d)
		for _, d := range aoc.AllDirection {
			if d == opposite {
				continue
			}

			steps := 1
			if d == s.d {
				if s.s >= max {
					continue
				}

				steps = s.s + 1
			} else {
				if s.s < min {
					continue
				}
			}

			next := a.Next(d, s.p)
			if next == nil {
				continue
			}

			ns := state{p: next, d: d, s: steps}
			nh := h + heat(next)

			updated := false
			oh, ok := a.heats[ns]
			if !ok || nh < oh {
				updated = true
				a.heats[ns] = nh
			}

			if !updated {
				continue
			}

			np := a.toTarget(next) + nh
			queue.Push(ns, np)
		}

		explored.Add(s)
	}

	return 0
}

func parseInput(r io.Reader, part int) int {
	scanner := bufio.NewScanner(r)

	a := newArea()
	for scanner.Scan() {
		a.AddPointsFromLine(scanner.Text())
	}
	a.start = a.Points[0][0]
	a.target = a.Points[len(a.Points)-1][len(a.Points[0])-1]

	if part == 1 {
		return a.explore(1, 3)
	}

	return a.explore(4, 10)
}
