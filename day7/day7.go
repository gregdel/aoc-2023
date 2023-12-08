package day

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{})
}

var scorePart1 = map[string]int{
	"2": 0, "3": 1, "4": 2, "5": 3, "6": 4,
	"7": 5, "8": 6, "9": 7, "T": 8, "J": 9,
	"Q": 10, "K": 11, "A": 12,
}

var scorePart2 = map[string]int{
	"J": 0, "2": 1, "3": 2, "4": 3, "5": 4,
	"6": 5, "7": 6, "8": 7, "9": 8, "T": 9,
	"Q": 10, "K": 11, "A": 12,
}

type gameType int

const (
	typeUndefined gameType = iota
	typeHighCard
	typeOnePair
	typeTwoPair
	typeThreeOfAKind
	typeFullHouse
	typeFourOfAKind
	typeFiveOfAKind
)

type hand struct {
	cards string
	bid   int
	Type  gameType
}

func (h *hand) less(other hand, part int) bool {
	ht := h.gameType(part)
	ot := other.gameType(part)
	if ht != ot {
		return ht < ot
	}

	for i := 0; i < len(h.cards); i++ {
		m := h.cards[i]
		o := other.cards[i]
		if m == o {
			continue
		}

		if part == 1 {
			return scorePart1[string(m)] < scorePart1[string(o)]
		}

		return scorePart2[string(m)] < scorePart2[string(o)]
	}

	panic("WUT ?")
}

func (h *hand) gameType(part int) gameType {
	if h.Type != typeUndefined {
		return h.Type
	}

	count := map[string]int{}
	for _, c := range h.cards {
		count[string(c)]++
	}

	numberOfJ, ok := count["J"]
	if ok && part == 2 {
		delete(count, "J")
		bestCard := ""
		bestCount := 0
		for card, nb := range count {
			if bestCard == "" ||
				nb > bestCount ||
				(nb == bestCount && scorePart2[card] > scorePart2[bestCard]) {
				bestCard = card
				bestCount = nb
			}
		}
		count[bestCard] += numberOfJ
	}

	type oc struct {
		card  string
		count int
	}
	cards := []oc{}
	for card, c := range count {
		cards = append(cards, oc{card: card, count: c})
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].count > cards[j].count })

	switch len(cards) {
	case 1:
		h.Type = typeFiveOfAKind
	case 2:
		if cards[0].count == 4 {
			h.Type = typeFourOfAKind
			break
		}

		if cards[0].count == 3 {
			h.Type = typeFullHouse
			break
		}
	case 3:
		if cards[0].count == 3 {
			h.Type = typeThreeOfAKind
			break
		}

		if cards[0].count == 2 && cards[1].count == 2 {
			h.Type = typeTwoPair
			break
		}
	case 4:
		if cards[0].count == 2 {
			h.Type = typeOnePair
		}
	case 5:
		h.Type = typeHighCard
	}

	return h.Type
}

type game struct {
	hands []hand
}

func (g *game) sort(part int) {
	sort.Slice(g.hands, func(i, j int) bool {
		return g.hands[i].less(g.hands[j], part)
	})
}

func newGame() *game {
	return &game{
		hands: []hand{},
	}
}

type day struct{}

func (d *day) Day() int {
	return 7
}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("6440", "246795406", "5905", "249356515").Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	game := d.parseInput(r, part)
	game.sort(part)

	result := 0
	for i, hand := range game.hands {
		rank := i + 1
		result += (rank * hand.bid)
	}

	return strconv.Itoa(result), nil
}

func (d *day) parseInput(r io.Reader, part int) *game {
	scanner := bufio.NewScanner(r)

	g := newGame()
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		g.hands = append(g.hands, hand{
			cards: fields[0],
			bid:   aoc.MustGet(strconv.Atoi(fields[1])),
		})
	}

	return g
}
