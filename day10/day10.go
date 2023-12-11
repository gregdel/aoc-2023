package day

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fatih/color"
	aoc "github.com/gregdel/aoc2023/lib"
)

type direction int

const (
	directionUp direction = iota
	directionDown
	directionRight
	directionLeft
)

var opposites = []direction{
	directionDown,
	directionUp,
	directionLeft,
	directionRight,
}

var dirStrs = []string{"Up", "Down", "Right", "Left"}

func oppositeDirection(d direction) direction {
	return opposites[d]
}

func dirStr(d direction) string {
	return dirStrs[d]
}

func init() {
	aoc.Register(&day{}, 10)
}

type tiles struct {
	start *tile
	data  [][]*tile
}

func (t *tiles) findPath() int {
	var nextDirection direction
	current := t.start
	current.explored = true
	current.inPath = true
	for i := 0; i < 4; i++ {
		d := direction(i)
		next := t.move(t.start, d)
		if next != nil && next.canComeFrom(oppositeDirection(d)) {
			nextDirection = d
			current = next
			break
		}
	}

	path := []*tile{t.start}

	for {
		current.explored = true
		current.inPath = true

		next, nd := t.next(current, nextDirection)
		if next == nil {
			panic("WTF")
		}
		if next == t.start {
			break
		}
		path = append(path, current)
		current = next
		nextDirection = nd
	}

	return len(path)/2 + 1
}

func expandTile(tiles [][]*tile, t *tile, x, y int, signs ...rune) {
	// Dots in the corner
	tiles[y][x] = newTile(x, y, '.', false, true)
	tiles[y+2][x] = newTile(x, y+2, '.', false, true)
	tiles[y][x+2] = newTile(x+2, y, '.', false, true)
	tiles[y+2][x+2] = newTile(x+2, y+2, '.', false, true)

	// Same tile in the center
	tiles[y+1][x+1] = newTile(x+1, y+1, t.sign, t.explored, false)

	isExplored := func(sign rune) bool {
		if sign == '.' {
			return false
		}

		return t.explored
	}

	up, down, left, right := signs[0], signs[1], signs[2], signs[3]
	tiles[y][x+1] = newTile(x+1, y, up, isExplored(up), true)
	tiles[y+2][x+1] = newTile(x+1, y+2, down, isExplored(down), true)
	tiles[y+1][x] = newTile(x, y+1, left, isExplored(left), true)
	tiles[y+1][x+2] = newTile(x+2, y+1, right, isExplored(right), true)
}

func (t *tiles) expand() {
	newLenY := 3 * len(t.data)
	newTiles := make([][]*tile, newLenY)
	for y := 0; y < len(t.data); y++ {
		idY := y * 3
		newLenX := 3 * len(t.data[0])
		newTiles[idY] = make([]*tile, newLenX)
		newTiles[idY+1] = make([]*tile, newLenX)
		newTiles[idY+2] = make([]*tile, newLenX)
		for x := 0; x < len(t.data[0]); x++ {
			tile := t.data[y][x]
			idX := x * 3

			switch tile.sign {
			case '─':
				expandTile(newTiles, tile, idX, idY, '.', '.', '─', '─')
			case '│':
				expandTile(newTiles, tile, idX, idY, '│', '│', '.', '.')
			case '└':
				expandTile(newTiles, tile, idX, idY, '│', '.', '.', '─')
			case '┌':
				expandTile(newTiles, tile, idX, idY, '.', '│', '.', '─')
			case '┘':
				expandTile(newTiles, tile, idX, idY, '│', '.', '─', '.')
			case '┐':
				expandTile(newTiles, tile, idX, idY, '.', '│', '─', '.')
			case 'S':
				expandTile(newTiles, tile, idX, idY, '│', '│', '─', '─')
			default:
				expandTile(newTiles, tile, idX, idY, '.', '.', '.', '.')
			}
		}
	}
	t.data = newTiles
}

func (t *tiles) explore() int {
	for _, x := range []int{0, len(t.data[0]) - 1} {
		for y := 0; y < len(t.data); y++ {
			t.exploreTile(t.data[y][x])
		}
	}

	for _, y := range []int{0, len(t.data) - 1} {
		for x := 0; x < len(t.data[0]); x++ {
			t.exploreTile(t.data[y][x])
		}
	}

	inside := 0
	for y := 0; y < len(t.data); y++ {
		for x := 0; x < len(t.data[0]); x++ {
			tile := t.data[y][x]
			if tile.explored || tile.fake {
				continue
			}
			inside++
		}
	}

	return inside
}

func (t *tiles) exploreTile(tile *tile) {
	// if tile.explored && tile.sign != '|' {
	if tile.explored {
		return
	}

	tile.explored = true
	for i := 0; i < 4; i++ {
		direction := direction(i)
		next := t.move(tile, direction)
		// if next == nil || (next.explored && next.sign != '│') {
		if next == nil || next.explored {
			continue
		}

		t.exploreTile(next)
	}
}

func (t *tiles) next(start *tile, d direction) (*tile, direction) {
	from := oppositeDirection(d)
	var other direction
	for i, v := range start.directions {
		d := direction(i)
		if !v || d == from {
			continue
		}

		other = d
		break
	}

	return t.move(start, other), other
}

func (t *tiles) move(start *tile, d direction) *tile {
	switch d {
	case directionUp:
		return t.up(start.x, start.y)
	case directionDown:
		return t.down(start.x, start.y)
	case directionRight:
		return t.right(start.x, start.y)
	case directionLeft:
		return t.left(start.x, start.y)
	}
	return nil
}

func (t *tiles) up(x, y int) *tile {
	if y == 0 {
		return nil
	}

	return t.data[y-1][x]
}

func (t *tiles) down(x, y int) *tile {
	if y == len(t.data)-1 {
		return nil
	}

	return t.data[y+1][x]
}

func (t *tiles) left(x, y int) *tile {
	if x == 0 {
		return nil
	}

	return t.data[y][x-1]
}

func (t *tiles) right(x, y int) *tile {
	if x == len(t.data[y])-1 {
		return nil
	}

	return t.data[y][x+1]
}

func (t *tiles) String() string {
	exploredColor := color.New(color.FgRed).SprintFunc()
	inPathColor := color.New(color.FgGreen).SprintFunc()
	var out strings.Builder
	for y := 0; y < len(t.data); y++ {
		for x := 0; x < len(t.data[y]); x++ {
			tile := t.data[y][x]
			if tile.inPath && tile.explored {
				out.WriteString(inPathColor(string(tile.sign)))
			} else if tile.explored {
				out.WriteString(exploredColor(string(tile.sign)))
			} else {
				out.WriteRune(tile.sign)
			}
		}
		out.WriteRune('\n')
	}
	return out.String()
}

func newTiles() *tiles {
	return &tiles{
		data: [][]*tile{},
	}
}

type tile struct {
	x, y              int
	sign              rune
	isStart, explored bool
	inPath            bool
	fake              bool
	directions        []bool
}

func (t *tile) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[%d;%d] %c", t.x, t.y, t.sign)
	for i, v := range t.directions {
		fmt.Fprintf(&b, " %s:%t", dirStr(direction(i)), v)
	}
	return b.String()
}

func (t *tile) canComeFrom(from direction) bool {
	return t.directions[from]
}

func newTile(x, y int, c rune, explored, fake bool) *tile {
	tile := &tile{
		x: x, y: y,
		sign:       c,
		explored:   explored,
		fake:       fake,
		directions: make([]bool, 4),
	}

	switch c {
	case 'S':
		tile.isStart = true
	case '|', '│':
		tile.directions[directionUp] = true
		tile.directions[directionDown] = true
		tile.sign = '│'
	case '-', '─':
		tile.directions[directionLeft] = true
		tile.directions[directionRight] = true
		tile.sign = '─'
	case 'L', '└':
		tile.directions[directionUp] = true
		tile.directions[directionRight] = true
		tile.sign = '└'
	case 'J', '┘':
		tile.directions[directionUp] = true
		tile.directions[directionLeft] = true
		tile.sign = '┘'
	case '7', '┐':
		tile.directions[directionDown] = true
		tile.directions[directionLeft] = true
		tile.sign = '┐'
	case 'F', '┌':
		tile.directions[directionDown] = true
		tile.directions[directionRight] = true
		tile.sign = '┌'
	}
	return tile
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("8", "7173", "10", "291").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	tiles := parseInput(r, part)
	result := tiles.findPath()
	if part == 2 {
		tiles.expand()
		result = tiles.explore()
	}
	return strconv.Itoa(result), nil
}

func parseInput(r io.Reader, part int) *tiles {
	scanner := bufio.NewScanner(r)

	tiles := newTiles()

	y := 0
	for scanner.Scan() {
		line := []*tile{}
		for x, c := range scanner.Text() {
			nt := newTile(x, y, c, false, false)
			if nt.isStart {
				tiles.start = nt
			}
			line = append(line, nt)
		}
		tiles.data = append(tiles.data, line)
		y++
	}

	return tiles
}
