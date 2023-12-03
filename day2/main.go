package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

func solve1(games []game) int {
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

	return result
}

func solve2(games []game) int {
	result := 0
	for _, game := range games {
		result += game.power()
	}

	return result
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Missing filename")
	}

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	games, err := parseInput(file)
	if err != nil {
		return err
	}

	result := solve2(games)
	fmt.Printf("Result %d\n", result)
	return nil
}
