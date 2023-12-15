package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 15)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("1320", "515210", "145", "246762").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	result := parseInput(r, part)
	return strconv.Itoa(result), nil
}

func hash(input string) int {
	r := 0
	for _, c := range input {
		r += int(c)
		r *= 17
		r %= 256
	}
	return r
}

func part1(tokens []string) int {
	result := 0
	for _, token := range tokens {
		result += hash(token)
	}
	return result
}

type lens struct {
	label string
	value int

	prev, next *lens
}

func newLens(label string, value int) *lens {
	return &lens{label: label, value: value}
}

type box struct {
	number      int
	index       map[string]*lens
	first, last *lens
}

func newBox(number int) *box {
	return &box{
		number: number,
		index:  map[string]*lens{},
	}
}

func (b *box) String() string {
	var out strings.Builder
	fmt.Fprintf(&out, "Box %d:", b.number)
	cur := b.first
	for cur != nil {
		fmt.Fprintf(&out, " [%s %d]", cur.label, cur.value)
		cur = cur.next
	}
	return out.String()
}

func (b *box) score() int {
	if b.first == nil {
		return 0
	}

	i := 1
	ret := 0
	base := 1 + b.number

	cur := b.first
	for cur != nil {
		ret += base * i * cur.value
		cur = cur.next
		i++
	}

	return ret
}

func (b *box) add(l *lens) {
	cur, ok := b.index[l.label]
	if ok {
		cur.value = l.value
		return
	}
	b.index[l.label] = l

	if b.first == nil {
		b.first = l
		b.last = l
		return
	}

	last := b.last
	last.next = l
	l.prev = last
	b.last = l
}

func (b *box) remove(l *lens) {
	cur, ok := b.index[l.label]
	if !ok {
		return
	}

	prev := cur.prev
	next := cur.next
	delete(b.index, l.label)

	// Only element
	if prev == nil && next == nil {
		b.first = nil
		b.last = nil
		return
	}

	// Last element
	if next == nil {
		prev.next = nil
		b.last = prev
		return
	}

	// First element
	if prev == nil {
		next.prev = nil
		b.first = next
		return
	}

	prev.next = next
	next.prev = prev
}

func part2(tokens []string) int {
	boxes := map[int]*box{}
	for _, token := range tokens {
		label := ""
		var sign rune
		var value int
		for _, c := range token {
			if unicode.IsLetter(c) {
				label += string(c)
			} else if unicode.IsNumber(c) {
				value = aoc.MustGet(strconv.Atoi(string(c)))
			} else {
				sign = c
			}
		}

		idx := hash(label)
		lens := newLens(label, value)

		box, ok := boxes[idx]
		if !ok {
			box = newBox(idx)
			boxes[idx] = box
		}

		if sign == '=' {
			box.add(lens)
		} else {
			box.remove(lens)
		}
	}

	ret := 0
	for _, b := range boxes {
		score := b.score()
		if score == 0 {
			continue
		}

		ret += b.score()
	}

	return ret
}

func parseInput(r io.Reader, part int) int {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	tokens := strings.Split(scanner.Text(), ",")

	if part == 1 {
		return part1(tokens)
	}

	return part2(tokens)
}
