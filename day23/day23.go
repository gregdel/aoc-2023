package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 23)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("94", "2010", "154", "6318").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	a := parseInput(r)
	a.findPath(a.start, aoc.DirectionDown, part)
	a.visitedNodes.Reset()
	a.findMax(a.start, 0)
	result := a.end.maxDistance
	return strconv.Itoa(result), nil
}

type edge struct {
	distance int
	from, to *node
}

type node struct {
	edges       []*edge
	point       *aoc.Point
	maxDistance int
}

func (n *node) String() string {
	return fmt.Sprintf("%s edges:%d", n.point, len(n.edges))
}

func newNode(p *aoc.Point) *node {
	return &node{
		edges: []*edge{},
		point: p,
	}
}

type area struct {
	*aoc.Map2D
	start, end   *node
	nodes        map[*aoc.Point]*node
	printedNodes aoc.Set[*node]
	visitedNodes aoc.Set[*node]
	canReachEnd  map[*node]bool
}

func newArea() *area {
	return &area{
		Map2D:        aoc.NewMap2D(),
		nodes:        map[*aoc.Point]*node{},
		printedNodes: aoc.NewSet[*node](),
		visitedNodes: aoc.NewSet[*node](),
		canReachEnd:  map[*node]bool{},
	}
}

var slopes = map[rune]aoc.Direction{
	'>': aoc.DirectionRight,
	'<': aoc.DirectionLeft,
	'v': aoc.DirectionDown,
}

func (a *area) getNode(p *aoc.Point) (*node, bool) {
	n, ok := a.nodes[p]
	if !ok {
		n = newNode(p)
		a.nodes[p] = n
	}
	return n, ok
}

func (a *area) findMax(node *node, currentMax int) {
	if a.visitedNodes.Has(node) {
		return
	}
	a.visitedNodes.Add(node)
	defer a.visitedNodes.Remove(node)

	if node.maxDistance < currentMax {
		node.maxDistance = currentMax
	}

	for _, e := range node.edges {
		a.findMax(e.to, currentMax+e.distance)
	}
}

func (a *area) findPath(node *node, id aoc.Direction, part int) {
	a.visitedNodes.Add(node)
	defer a.visitedNodes.Remove(node)
	distance := 1
	prevDirection := id
	next := a.Next(id, node.point)
	for {
		od := aoc.OppositeDirection(prevDirection)
		for _, d := range aoc.AllDirection {
			if d == od {
				continue
			}

			p := a.Next(d, next)
			if p == nil || p.C == '#' {
				continue
			}

			distance++
			if p.C == '.' {
				next = p
				prevDirection = d
				if p != a.end.point {
					break
				}
			}

			if p != a.end.point {
				distance++
				p = a.Next(d, p)
			}

			to, explored := a.getNode(p)
			node.edges = append(node.edges, &edge{
				distance: distance,
				from:     node,
				to:       to,
			})

			if explored || to == a.end {
				return
			}

			if part == 2 {
				to.edges = append(to.edges, &edge{
					distance: distance,
					from:     to,
					to:       node,
				})

				if a.visitedNodes.Has(to) {
					return
				}
			}

			nod := aoc.OppositeDirection(d)
			for _, nd := range aoc.AllDirection {
				if nd == nod {
					continue
				}

				nn := a.Next(nd, p)
				if nn == nil {
					continue
				}

				if nn.C == '#' {
					continue

				}

				sd := slopes[nn.C]
				if part == 1 && sd != nd {
					continue
				}
				a.findPath(to, nd, part)
			}

			return
		}
	}
}

func parseInput(r io.Reader) *area {
	scanner := bufio.NewScanner(r)

	a := newArea()
	for scanner.Scan() {
		a.AddPointsFromLine(scanner.Text())
	}

	a.start, _ = a.getNode(a.Points[0][1])
	a.end, _ = a.getNode(a.Points[a.Height()-1][a.Width()-2])

	return a
}
