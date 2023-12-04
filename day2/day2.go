package day2

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day2{})
}

type day2 struct{}

func (d *day2) Day() int {
	return 2
}

func (d *day2) Solve(r io.Reader, part int) (string, error) {
	games, err := parseInput(r)
	if err != nil {
		return "", err
	}

	if part == 1 {
		return solve1(games), nil
	}

	return solve2(games), nil
}

func (d *day2) Expect(part int, test bool) string {
	return aoc.NewResult("8", "1734", "2286", "70387").Expect(part, test)
}

func solve1(games []game) string {
	max := [3]int{}
	max[red] = 12
	max[green] = 13
	max[blue] = 14

	result := 0
	for _, game := range games {
		if game.isPossible(max) {
			result += game.id
		}
	}
	return strconv.Itoa(result)
}

func solve2(games []game) string {
	result := 0
	for _, game := range games {
		result += game.power()
	}
	return strconv.Itoa(result)
}

const (
	green = iota
	blue
	red
)

type game struct {
	id   int
	sets [][3]int
}

func (g *game) isPossible(max [3]int) bool {
	for _, s := range g.sets {
		for i := 0; i < 3; i++ {
			if s[i] > max[i] {
				return false
			}
		}
	}

	return true
}

func (g *game) power() int {
	max := [3]int{}
	for _, s := range g.sets {
		for i := 0; i < 3; i++ {
			if s[i] > max[i] {
				max[i] = s[i]
			}
		}
	}

	power := 1
	for i := 0; i < 3; i++ {
		power *= max[i]
	}

	return power
}

func parseInput(reader io.Reader) ([]game, error) {
	games := []game{}
	scanner := bufio.NewScanner(reader)
	i := 0
	for scanner.Scan() {
		i++
		game := game{
			id:   i,
			sets: [][3]int{},
		}

		line := scanner.Text()
		idx := strings.Index(line, ":")

		sets := strings.Split(line[idx+1:], ";")
		for _, set := range sets {
			currentSet := [3]int{0, 0, 0}
			cubes := strings.Split(set, ",")
			for _, cube := range cubes {
				parts := strings.Split(cube, " ")
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid parts")
				}

				value, err := strconv.Atoi(parts[1])
				if err != nil {
					return nil, err
				}

				switch parts[2] {
				case "red":
					currentSet[red] = value
				case "green":
					currentSet[green] = value
				case "blue":
					currentSet[blue] = value
				}
			}
			game.sets = append(game.sets, currentSet)
		}

		games = append(games, game)
	}

	return games, nil
}
