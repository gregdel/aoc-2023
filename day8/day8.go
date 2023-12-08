package day

import (
	"bufio"
	"io"
	"strconv"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 8)
}

type game struct {
	directions         string
	startNode, endNode *node
	nodes              map[string]*node
}

func newGame(directions string) *game {
	return &game{
		directions: directions,
		nodes:      map[string]*node{},
	}
}

func (g *game) stepsToEnd() int {
	step := 0
	current := g.nodes["AAA"]
	end := g.nodes["ZZZ"]

	for {
		for _, d := range g.directions {
			current = current.move(d)
			step++
			if current == end {
				return step
			}
		}
	}
}

func (g *game) stepsToEndZ(n *node) int {
	current := n

	steps := 0
	for {
		for _, d := range g.directions {
			current = current.move(d)
			steps++
			if current.endsWithZ() {
				return steps
			}
		}
	}
}

func (g *game) stepsToEndWithZ() int {
	steps := []int{}

	for _, n := range g.nodes {
		if !n.endsWithA() {
			continue
		}

		steps = append(steps, g.stepsToEndZ(n))
	}

	return aoc.LeastCommonMultiple(steps[0], steps[1], steps[2:]...)
}

type node struct {
	right, left             *node
	name, rightStr, leftStr string
}

func (n *node) endsWithA() bool {
	return n.name[2] == 'A'
}

func (n *node) endsWithZ() bool {
	return n.name[2] == 'Z'
}

func (n *node) move(direction rune) *node {
	switch direction {
	case 'L':
		return n.left
	case 'R':
		return n.right
	default:
		return nil
	}
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("6", "16697", "6", "10668805667831").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	game := d.parseInput(r, part)

	var result int
	if part == 1 {
		result = game.stepsToEnd()
	}

	if part == 2 {
		result = game.stepsToEndWithZ()
	}

	return strconv.Itoa(result), nil
}

func (d *day) parseInput(r io.Reader, part int) *game {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	g := newGame(scanner.Text())
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		name := line[0:3]
		left := line[7:10]
		right := line[12:15]
		g.nodes[name] = &node{
			name:     name,
			leftStr:  left,
			rightStr: right,
		}
	}

	for _, node := range g.nodes {
		node.left = g.nodes[node.leftStr]
		node.right = g.nodes[node.rightStr]
	}

	return g
}
