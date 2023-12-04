package day4

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day4{})
}

type card struct {
	count int
	wins  int
}

func (c *card) score() int {
	if c.wins == 0 {
		return 0
	}

	return 1 << (c.wins - 1)
}

type day4 struct {
	cards []card
}

func (d *day4) Day() int {
	return 4
}

func (d *day4) Solve(r io.Reader, part int) (string, error) {
	d.parseInput(r)
	if part == 1 {
		return d.solve1(r)
	}

	return d.solve2(r)
}

func (d *day4) Expect(part int, test bool) string {
	return aoc.NewResult("13", "20107", "30", "8172507").Expect(part, test)
}

func (d *day4) parseInput(r io.Reader) {
	d.cards = []card{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, ":")
		line = parts[1]

		parts = strings.Split(line, "|")
		winningNumbers := strings.Fields(parts[0])
		numbers := strings.Fields(parts[1])

		winnings := map[string]struct{}{}
		for _, n := range winningNumbers {
			winnings[n] = struct{}{}
		}

		wins := 0
		for _, n := range numbers {
			_, ok := winnings[n]
			if ok {
				wins++
			}
		}

		d.cards = append(d.cards, card{
			count: 1,
			wins:  wins,
		})
	}
}

func (d *day4) solve1(r io.Reader) (string, error) {
	result := 0
	for _, card := range d.cards {
		result += card.score()
	}

	return strconv.Itoa(result), nil
}

func (d *day4) solve2(r io.Reader) (string, error) {
	result := 0
	for idx, card := range d.cards {
		for i := 1; i <= card.wins; i++ {
			d.cards[idx+i].count += card.count
		}

		result += card.count
	}

	return strconv.Itoa(result), nil
}
