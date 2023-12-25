package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

type nodeType int

const (
	typeNop nodeType = iota
	typeBroadcast
	typeFlipFlop
	typeConjunction
)

var typeNames = []string{
	"output", "broadcaster", "flip-flop", "conjuction",
}

type node struct {
	name    string
	t       nodeType
	inputs  []*node
	outputs []*node
	state   bool
}

func (n *node) String() string {
	ins, outs := []string{}, []string{}
	for _, i := range n.inputs {
		ins = append(ins, i.name)
	}
	for _, o := range n.outputs {
		outs = append(outs, o.name)
	}

	i, o := "none", "none"
	if len(ins) > 0 {
		i = strings.Join(ins, ",")
	}
	if len(outs) > 0 {
		o = strings.Join(outs, ",")
	}

	return fmt.Sprintf("%s type:%s state:%t inputs:%s outputs:%s",
		n.name, typeNames[int(n.t)], n.state, i, o)
}

func newNode(name string) *node {
	return &node{
		name:    name,
		inputs:  []*node{},
		outputs: []*node{},
	}
}

type fromTo struct {
	from, to *node
}

type pulses struct {
	low, high int
}

func (p *pulses) reset() {
	p.low, p.high = 0, 0
}

func (p pulses) value() int {
	return p.low * p.high
}

type stateMachine struct {
	nodes         map[string]*node
	toExplore     []fromTo
	pulses, total pulses
	pushes        int
	minPushes     map[string]int
}

func (sm *stateMachine) addPulse(high bool) {
	if high {
		sm.pulses.high++
	} else {
		sm.pulses.low++
	}
}

func (sm *stateMachine) runNode(from, to *node) {
	switch to.t {
	case typeBroadcast:
		// Nothing to do
	case typeFlipFlop:
		if from.state {
			return
		}
		to.state = !to.state
	case typeConjunction:
		allHigh := true
		for _, i := range to.inputs {
			if !i.state {
				allHigh = false
			}
		}
		to.state = !allHigh
	default:
		return
	}

	endNodes := []string{"vz", "bq", "lt", "qh"}
	for _, endNode := range endNodes {
		if to.name == endNode && to.state {
			_, ok := sm.minPushes[endNode]
			if !ok {
				sm.minPushes[endNode] = sm.pushes
			}
		}
	}

	for _, o := range to.outputs {
		sm.addPulse(to.state)
		sm.toExplore = append(sm.toExplore, fromTo{from: to, to: o})
	}
}

func (sm *stateMachine) pushButton() {
	sm.toExplore = []fromTo{{to: sm.nodes["broadcaster"]}}
	sm.pulses.reset()
	sm.pushes++
	sm.addPulse(false)

	for len(sm.toExplore) != 0 {
		ft := sm.toExplore[0]
		sm.runNode(ft.from, ft.to)
		if len(sm.toExplore) == 1 {
			sm.toExplore = []fromTo{}
		} else {
			sm.toExplore = sm.toExplore[1:]
		}
	}

	sm.total.low += sm.pulses.low
	sm.total.high += sm.pulses.high
}

func (sm *stateMachine) run() int {
	for i := 0; i < 1000; i++ {
		sm.pushButton()
	}

	return sm.total.value()
}

func (sm *stateMachine) run2() int {
	_, ok := sm.nodes["rx"]
	if !ok {
		return 0
	}

	for {
		sm.pushButton()
		if len(sm.minPushes) == 4 {
			break
		}
	}

	return aoc.LeastCommonMultiple(
		sm.minPushes["vz"],
		sm.minPushes["bq"],
		sm.minPushes["lt"],
		sm.minPushes["qh"],
	)
}

func newStateMachine(nodes map[string]*node) *stateMachine {
	return &stateMachine{
		nodes:     nodes,
		total:     pulses{},
		minPushes: map[string]int{},
	}
}

func init() {
	aoc.Register(&day{}, 20)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult(
		"11687500", "737679780", "0", "227411378431763",
	).Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	sm := parseInput(r, part)
	ret := 0
	if part == 1 {
		ret = sm.run()
	} else {
		ret = sm.run2()
	}
	return strconv.Itoa(ret), nil
}

func parseInput(r io.Reader, part int) *stateMachine {
	scanner := bufio.NewScanner(r)

	nodes := map[string]*node{}
	for scanner.Scan() {
		line := scanner.Text()
		part0, part1, _ := strings.Cut(line, " -> ")

		var name string
		var t nodeType
		if part0 == "broadcaster" {
			name = part0
			t = typeBroadcast
		} else if part0 == "output" {
			name = part0
			t = typeNop
		} else {
			name = part0[1:]
			switch part0[0:1] {
			case "%":
				t = typeFlipFlop
			case "&":
				t = typeConjunction
			}
		}

		n, ok := nodes[name]
		if !ok {
			nodes[name] = newNode(name)
			n = nodes[name]
		}
		n.t = t

		output := strings.Split(part1, ", ")
		for _, o := range output {
			on, ok := nodes[o]
			if !ok {
				nodes[o] = newNode(o)
				on = nodes[o]
			}
			n.outputs = append(n.outputs, on)
			on.inputs = append(on.inputs, n)
		}
	}

	return newStateMachine(nodes)
}
