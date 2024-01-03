package day

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 22)
}

type day struct{}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("5", "446", "7", "60287").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	a := parseInput(r, part)
	a.settle()

	result := 0
	for _, b := range a.bricks {
		canBeDisintegrated := a.canBeDisintegrated(b)
		if part == 1 && canBeDisintegrated {
			result++
		}

		if part == 2 && !canBeDisintegrated {
			result += a.totalFalling(b)
		}
	}

	return strconv.Itoa(result), nil
}

type point struct {
	x, y, z int
}

func newPoint(x, y, z int) point {
	return point{x: x, y: y, z: z}
}

type brick struct {
	name                string
	from, to            point
	sameX, sameY, sameZ bool
	commonAxis          int
	above, under        aoc.Set[*brick]
}

func (b *brick) points() []point {
	if b.commonAxis == 3 {
		return []point{newPoint(b.from.x, b.from.y, b.from.z)}
	}

	points := []point{}

	if !b.sameX {
		for x := b.from.x; x <= b.to.x; x++ {
			points = append(points, newPoint(x, b.from.y, b.from.z))
		}
	}

	if !b.sameY {
		for y := b.from.y; y <= b.to.y; y++ {
			points = append(points, newPoint(b.from.x, y, b.from.z))
		}
	}

	if !b.sameZ {
		for z := b.from.z; z <= b.to.z; z++ {
			points = append(points, newPoint(b.from.x, b.from.y, z))
		}
	}

	return points
}

func newBrick(name string, input []int) *brick {
	b := &brick{
		above: aoc.NewSet[*brick](),
		under: aoc.NewSet[*brick](),
		name:  name,
		from:  point{x: input[0], y: input[1], z: input[2]},
		to:    point{x: input[3], y: input[4], z: input[5]},
	}
	b.sameX = b.from.x == b.to.x
	b.sameY = b.from.y == b.to.y
	b.sameZ = b.from.z == b.to.z
	for _, v := range []bool{b.sameX, b.sameY, b.sameZ} {
		if v {
			b.commonAxis++
		}
	}

	return b
}

type area struct {
	bricks    []*brick
	locations map[point]*brick
}

func newArea() *area {
	return &area{
		bricks:    []*brick{},
		locations: map[point]*brick{},
	}
}

func (a *area) sortBricks() {
	sort.Slice(a.bricks, func(i, j int) bool {
		return a.bricks[i].from.z < a.bricks[j].from.z
	})
}

func (a *area) updateLocations() {
	a.sortBricks()
	for _, b := range a.bricks {
		points := b.points()
		for _, p := range points {
			a.locations[p] = b
		}
	}
}

func (a *area) maxDown(b *brick) int {
	maxDown := 1_000_000

	points := b.points()
	if !b.sameZ {
		// Check only the lowest point
		points = []point{newPoint(b.from.x, b.from.y, b.from.z)}
	}

	for _, p := range points {
		maxDownForPoint := 0
		if p.z != 1 {
			for z := p.z - 1; z >= 1; z++ {
				_, ok := a.locations[newPoint(p.x, p.y, z)]
				if ok {
					break
				}
				maxDownForPoint++
			}
		}

		maxDown = aoc.Min(maxDown, maxDownForPoint)
	}
	return maxDown
}

func (a *area) settleBrick(b *brick) bool {
	moveDown := a.maxDown(b)
	if moveDown == 0 {
		return false
	}

	// Remove the points from their locations
	for _, p := range b.points() {
		delete(a.locations, p)
	}

	// Move the brick down
	b.from.z -= moveDown
	b.to.z -= moveDown

	// Update the points locations
	for _, p := range b.points() {
		a.locations[p] = b
	}

	return true
}

func (a *area) updateBrick(b *brick) {
	for _, p := range b.points() {
		ab, ok := a.locations[newPoint(p.x, p.y, p.z+1)]
		if ok && ab != b {
			b.above.Add(ab)
		}

		if p.z == 1 {
			continue
		}

		ub, ok := a.locations[newPoint(p.x, p.y, p.z-1)]
		if ok && ub != b {
			b.under.Add(ub)
		}
	}
}

func (a *area) canBeDisintegrated(b *brick) bool {
	if len(b.above) == 0 {
		return true
	}

	removable := true
	for ab := range b.above {
		if ab.under.Len() == 1 {
			removable = false
		}
	}

	return removable
}

func (a *area) settle() {
	i := 0
	for {
		i++
		allSettled := true
		for _, b := range a.bricks {
			if a.settleBrick(b) {
				allSettled = false
			}
		}

		if allSettled {
			break
		}
		a.sortBricks()
	}

	a.sortBricks()
	for _, b := range a.bricks {
		a.updateBrick(b)
	}
}

func (a *area) totalFalling(b *brick) int {
	falling := aoc.NewSet[*brick]()

	queue := aoc.NewQueue[*brick]()
	for a := range b.above {
		queue.Push(a)
	}
	falling.Add(b)

	for queue.Len() != 0 {
		e := queue.Pop()
		isFalling := true
		for ub := range e.under {
			if !falling.Has(ub) {
				isFalling = false
				break
			}
		}

		if !isFalling {
			continue
		}

		falling.Add(e)
		for a := range e.above {
			queue.Push(a)
		}
	}

	return falling.Len() - 1
}

func parseInput(r io.Reader, part int) *area {
	scanner := bufio.NewScanner(r)

	a := newArea()
	i := 0
	for scanner.Scan() {
		i++
		input := strings.ReplaceAll(scanner.Text(), ",", " ")
		input = strings.ReplaceAll(input, "~", " ")
		values := aoc.IntsFromString(input)
		a.bricks = append(a.bricks, newBrick(fmt.Sprintf("%4d", i), values))
	}

	a.updateLocations()
	return a
}
