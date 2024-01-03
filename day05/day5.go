package day

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 5)
}

type interval struct {
	source, destination, length int64
}

type seed struct {
	source, length int64
}

func (i *interval) split(s seed) ([]seed, bool) {
	start := i.source
	end := i.source + i.length - 1

	seedStart := s.source
	seedEnd := s.source + s.length - 1

	// Outside of this interval
	if seedEnd < start || seedStart > end {
		return nil, false
	}

	// Around this interval
	if seedStart < start && seedEnd > end {
		seeds := []seed{
			{source: seedStart, length: start - seedStart},
			{source: start, length: i.length},
			{source: end + 1, length: seedEnd - end},
		}
		return seeds, true
	}

	var removed int64 = 0
	seeds := []seed{}
	newStart := seedStart
	if seedStart < start {
		removed = start - seedStart
		newStart = start
		seeds = append(seeds, seed{
			source: seedStart, length: removed,
		})
	}

	if seedEnd > end {
		removed = seedEnd - end
		seeds = append(seeds, seed{
			source: end + 1, length: removed,
		})
	}

	seeds = append(seeds, seed{
		source: newStart, length: s.length - removed,
	})

	updated := (len(seeds) != 1 || seeds[0].source != s.source || seeds[0].length != s.length)

	return seeds, updated
}

func (i *interval) translate(s seed) (seed, bool) {
	start := i.source
	end := i.source + i.length - 1

	seedStart := s.source
	seedEnd := s.source + s.length - 1

	inRange := false
	if seedStart >= start && seedEnd <= end {
		inRange = true
	}

	if !inRange {
		return seed{}, false
	}

	diff := i.destination - i.source
	return seed{
		source: s.source + diff,
		length: s.length,
	}, true
}

// Map represents a map
type Map struct {
	source, destination string
	ranges              []interval
}

func (m *Map) split(seeds []seed) []seed {
	toSplit := seeds

	for {
		allGood := true
		result := []seed{}
		for _, seed := range toSplit {
			splitted := false
			for _, r := range m.ranges {
				parts, ok := r.split(seed)
				if ok {
					allGood = false
					splitted = true
					result = append(result, parts...)
				}
			}

			if !splitted {
				result = append(result, seed)
			}
		}

		toSplit = result
		if allGood {
			break
		}
	}

	tmp := map[seed]struct{}{}
	for _, s := range toSplit {
		tmp[s] = struct{}{}
	}

	result := []seed{}
	for s := range tmp {
		result = append(result, s)
	}
	return result
}

func (m *Map) translate(seeds []seed) []seed {
	newSeeds := []seed{}
	for _, seed := range seeds {
		translated := false
		for _, r := range m.ranges {
			s, ok := r.translate(seed)
			if ok {
				translated = true
				newSeeds = append(newSeeds, s)
			}
		}

		if !translated {
			newSeeds = append(newSeeds, seed)
		}
	}
	return newSeeds
}

func (m *Map) newValues(seeds []seed) []seed {
	seeds = m.split(seeds)
	seeds = m.translate(seeds)
	return seeds
}

type day struct {
	seeds []seed
	index map[string]*Map
	maps  []*Map
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	d.parseInput(r, part)
	seeds := d.translateAll(d.seeds)

	var result int64 = 0
	found := false
	for _, seed := range seeds {
		value := seed.source
		if !found {
			result = value
			found = true
		} else {
			if value < result {
				result = value
			}
		}
	}

	return strconv.FormatInt(result, 10), nil
}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult("35", "331445006", "46", "6472060").Expect(part, test)
}

func (d *day) translateAll(seeds []seed) []seed {
	values := seeds
	current := "seed"

	for {
		m := d.index[current]
		values = m.newValues(values)
		if current == "humidity" {
			break
		}
		current = m.destination
	}

	return values
}

func parseNumbers(input string) []int64 {
	output := []int64{}
	for _, str := range strings.Fields(input) {
		v := aoc.MustGet(strconv.ParseInt(str, 10, 64))
		output = append(output, v)
	}

	return output
}

func (d *day) parseInput(r io.Reader, part int) {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	line := scanner.Text()
	parts := strings.Split(line, ":")

	d.seeds = []seed{}
	seeds := parseNumbers(parts[1])
	if part == 1 {
		for _, s := range seeds {
			d.seeds = append(d.seeds, seed{source: s, length: 1})
		}
	} else {
		for i := 0; i < len(seeds); i += 2 {
			d.seeds = append(d.seeds, seed{
				source: seeds[i],
				length: seeds[i+1],
			})
		}
	}

	d.index = map[string]*Map{}
	d.maps = []*Map{}
	var m *Map
	for scanner.Scan() {
		line := scanner.Text()
		// New line new map
		if len(line) == 0 {
			m = &Map{
				ranges: []interval{},
			}
			continue
		}

		if m.source == "" {
			words := strings.Fields(line)
			parts = strings.Split(words[0], "-")
			m.source = parts[0]
			m.destination = parts[2]
			d.index[m.source] = m
			d.maps = append(d.maps, m)
			continue
		}

		numbers := parseNumbers(line)
		m.ranges = append(m.ranges, interval{
			destination: numbers[0],
			source:      numbers[1],
			length:      numbers[2],
		})
	}
}
