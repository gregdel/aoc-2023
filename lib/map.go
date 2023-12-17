package aoc

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// Direction represents a direction
type Direction uint8

// Available directions
const (
	DirectionUp Direction = iota
	DirectionDown
	DirectionRight
	DirectionLeft
)

// OppositeDirection returns the opposite direction of a direction.
func OppositeDirection(d Direction) Direction {
	switch d {
	case DirectionUp:
		return DirectionDown
	case DirectionDown:
		return DirectionUp
	case DirectionRight:
		return DirectionLeft
	case DirectionLeft:
		return DirectionRight
	}
	panic("Invalid direction")
}

// Point represents a 2d point
type Point struct {
	X, Y int
	C    rune
}

// NewPoint returns a new point
func NewPoint(x, y int, c rune) *Point {
	return &Point{X: x, Y: y, C: c}
}

// Map2D represents a 2DMap
type Map2D struct {
	Points [][]*Point
}

// NewMap2D returns a new 2D map
func NewMap2D() *Map2D {
	return &Map2D{
		Points: [][]*Point{},
	}
}

// AddPointsFromLine adds points to a 2d map from a string
func (m *Map2D) AddPointsFromLine(line string) {
	y := len(m.Points)
	points := make([]*Point, len(line))
	m.Points = append(m.Points, points)
	for x, c := range line {
		m.Points[y][x] = NewPoint(x, y, c)
	}
}

// ForAllPoints calls a function for each point in a map
func (m *Map2D) ForAllPoints(f func(p *Point), directions ...Direction) {
	d1, d2 := DirectionDown, DirectionRight
	if len(directions) == 2 {
		d1 = directions[0]
		d2 = directions[1]
	}

	switch d1 {
	case DirectionDown:
		switch d2 {
		case DirectionRight:
			for y := 0; y < len(m.Points); y++ {
				for x := 0; x < len(m.Points[y]); x++ {
					f(m.Points[y][x])
				}
			}
		case DirectionLeft:
			for y := 0; y < len(m.Points); y++ {
				for x := len(m.Points[y]) - 1; x >= 0; x-- {
					f(m.Points[y][x])
				}
			}
		default:
			panic("Invalid direction d2")
		}
	case DirectionUp:
		switch d2 {
		case DirectionRight:
			for y := len(m.Points) - 1; y >= 0; y-- {
				for x := 0; x < len(m.Points[y]); x++ {
					f(m.Points[y][x])
				}
			}
		case DirectionLeft:
			for y := len(m.Points) - 1; y >= 0; y-- {
				for x := len(m.Points[y]) - 1; x >= 0; x-- {
					f(m.Points[y][x])
				}
			}
		default:
			panic("Invalid direction d2")
		}
	default:
		panic("Invalid direction d1")
	}
}

// Next returns the next point in the given direction
func (m *Map2D) Next(d Direction, p *Point) *Point {
	switch d {
	case DirectionUp:
		if p.Y > 0 {
			return m.Points[p.Y-1][p.X]
		}
	case DirectionDown:
		if p.Y < len(m.Points)-1 {
			return m.Points[p.Y+1][p.X]
		}
	case DirectionLeft:
		if p.X > 0 {
			return m.Points[p.Y][p.X-1]
		}
	case DirectionRight:
		if p.X < len(m.Points[p.Y])-1 {
			return m.Points[p.Y][p.X+1]
		}
	}
	return nil
}

// SHA256Sum return the checksum of the map
func (m *Map2D) SHA256Sum() string {
	h := sha256.New()
	m.ForAllPoints(func(p *Point) {
		h.Write([]byte{byte(p.C)})
	})
	return fmt.Sprintf("%x", h.Sum(nil))
}

// String implements the stringer interface
func (m *Map2D) String() string {
	var out strings.Builder
	for y := 0; y < len(m.Points); y++ {
		for x := 0; x < len(m.Points[y]); x++ {
			p := m.Points[y][x]
			out.WriteRune(p.C)
		}
		out.WriteRune('\n')
	}
	return out.String()
}
