package day

import (
	"bufio"
	"io"
	"strconv"
	"unicode"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 3)
}

type element struct {
	x, y, length int
	isNumber     bool
	isValid      bool
	value        string
}

type day struct{}

func distance(first, last element) int {
	return last.x - (first.x + first.length)
}

func inRange(element, toCheck element) bool {
	min, max := element.x-1, element.x+element.length
	inRange := false
	if toCheck.x >= min && toCheck.x <= max {
		inRange = true
	}
	return inRange
}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("4361", "554003", "467835", "87263515").Expect(part, test)
}

func (d *day) analyse(r io.Reader, part int) [][]element {
	elements := d.parseInput(r, part)

	for y := 0; y < len(elements); y++ {
		for i := range elements[y] {
			element := &elements[y][i]
			if !element.isNumber {
				continue
			}

			// Prev on same line
			if i > 0 {
				prev := elements[y][i-1]
				if distance(prev, *element) == 0 {
					element.isValid = true
					continue
				}
			}

			// Next on same line
			if i < len(elements[y])-1 {
				next := elements[y][i+1]
				if distance(*element, next) == 0 {
					element.isValid = true
					continue
				}
			}

			// Previous line
			if y > 0 {
				for _, e := range elements[y-1] {
					if e.isNumber {
						continue
					}

					if inRange(*element, e) {
						element.isValid = true
						break
					}
				}

				if element.isValid {
					continue
				}
			}

			// Next line
			if y < len(elements)-1 {
				for _, e := range elements[y+1] {
					if e.isNumber {
						continue
					}

					if inRange(*element, e) {
						element.isValid = true
						break
					}
				}

				if element.isValid {
					continue
				}
			}
		}
	}

	return elements
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	if part == 1 {
		return d.solve1(r, part)
	}

	return d.solve2(r, part)
}

func (d *day) solve1(r io.Reader, part int) (string, error) {
	elements := d.analyse(r, part)

	result := 0
	for y := 0; y < len(elements); y++ {
		for _, e := range elements[y] {
			if e.isValid {
				result += aoc.MustGet(strconv.Atoi(e.value))
			}
		}
	}

	return strconv.Itoa(result), nil
}

func (d *day) solve2(r io.Reader, part int) (string, error) {
	elements := d.analyse(r, part)
	newElements := make([][]element, len(elements))

	for y := 0; y < len(elements); y++ {
		newElements[y] = []element{}
		for _, e := range elements[y] {
			if e.isNumber && !e.isValid {
				continue
			}
			newElements[y] = append(newElements[y], e)
		}
	}

	result := 0

	elements = newElements
	for y := 0; y < len(elements); y++ {
		for x := range elements[y] {
			e := elements[y][x]
			if e.isNumber {
				continue
			}

			parts := []*element{}

			// Prev on same line
			if x > 0 {
				prev := elements[y][x-1]
				if prev.isNumber && distance(prev, e) == 0 {
					parts = append(parts, &prev)
				}
			}

			// Next on same line
			if x < len(elements[y])-1 {
				next := elements[y][x+1]
				if next.isNumber && distance(e, next) == 0 {
					parts = append(parts, &next)
				}
			}

			// Previous line
			if y > 0 {
				for i, el := range elements[y-1] {
					if el.isNumber && inRange(el, e) {
						parts = append(parts, &elements[y-1][i])
					}
				}
			}

			// Next line
			if y < len(elements)-1 {
				for i, el := range elements[y+1] {
					if el.isNumber && inRange(el, e) {
						parts = append(parts, &elements[y+1][i])
					}
				}
			}

			if len(parts) == 2 {
				p1 := aoc.MustGet(strconv.Atoi(parts[0].value))
				p2 := aoc.MustGet(strconv.Atoi(parts[1].value))
				result += p1 * p2
			}
		}
	}

	return strconv.Itoa(result), nil
}

func (d *day) parseInput(r io.Reader, part int) [][]element {
	scanner := bufio.NewScanner(r)
	result := [][]element{}

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		elements := []element{}
		current := element{}
		for x, c := range line {
			if c == '.' {
				if current.length != 0 {
					elements = append(elements, current)
					current = element{}
				}
				continue
			}

			if !unicode.IsDigit(c) {
				if part == 2 && c != '*' {
					continue
				}

				if current.length != 0 {
					elements = append(elements, current)
					current = element{}
				}

				elements = append(elements, element{
					x:      x,
					y:      y,
					length: 1,
					value:  string(c),
				})
				continue
			}

			if current.length == 0 {
				current.isNumber = true
				current.x = x
				current.y = y
				current.length = 1
				current.value = string(c)
			} else {
				current.value += string(c)
				current.length++
			}
		}

		if current.length != 0 {
			elements = append(elements, current)
			current = element{}
		}

		y++
		result = append(result, elements)
	}

	return result
}
