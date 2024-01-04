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
	return aoc.NewResult("94", "2010", "", "").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	a := parseInput(r)
	a.findPath(a.start, aoc.DirectionDown, part)
	a.printGraph(a.start)

	result := 0
	if part == 1 {
		result, _ = a.findMax(a.start)
	} else {
		result, _ = a.findMax(a.start)
	}

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
}

func newArea() *area {
	return &area{
		Map2D:        aoc.NewMap2D(),
		nodes:        map[*aoc.Point]*node{},
		printedNodes: aoc.NewSet[*node](),
		visitedNodes: aoc.NewSet[*node](),
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

func (a *area) printGraph(node *node) {
	if a.printedNodes.Has(node) {
		return
	}
	a.printedNodes.Add(node)

	// fmt.Printf("Node: %+v\n", node)
	for _, e := range node.edges {
		fmt.Printf("%03d%03d -> %03d%03d [label=\"%d\"]\n",
			e.from.point.X, e.from.point.Y,
			e.to.point.X, e.to.point.Y, e.distance)
		a.printGraph(e.to)
	}
}

func (a *area) findMax(node *node) (int, bool) {
	if node.maxDistance != 0 {
		return node.maxDistance, true
	}

	allOk := true
	for _, e := range node.edges {
		m, ok := a.findMax(e.to)
		if !ok {
			allOk = false
			continue
		}
		node.maxDistance = aoc.Max(node.maxDistance, m+e.distance)
	}

	return node.maxDistance, allOk
}

func (a *area) findPath(node *node, id aoc.Direction, part int) {
	fmt.Println("Searching for path", node.point, id)
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
				fmt.Println("Found a tile", p)
				next = p
				prevDirection = d
				if p != a.end.point {
					break
				}
			}

			if p != a.end.point {
				distance++
				p = a.Next(slopes[p.C], p)
			}

			to, explored := a.getNode(p)
			// check if we visited to, skip if part 2
			if part == 2 && a.visitedNodes.Has(to) {
				return
			}
			edge := &edge{
				distance: distance,
				from:     node,
				to:       to,
			}
			node.edges = append(node.edges, edge)

			if explored || to == a.end {
				return
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
					fmt.Printf("skipping direction sd:%s nd:%s\n", sd, nd)
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

	a.start = newNode(a.Points[0][1])
	a.end = newNode(a.Points[a.Height()-1][a.Width()-2])

	return a
}
