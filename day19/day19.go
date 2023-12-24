package day

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"strconv"
	"strings"

	aoc "github.com/gregdel/aoc2023/lib"
)

func init() {
	aoc.Register(&day{}, 19)
}

type ratingRange struct {
	min, max int
}

func newRatingRange() ratingRange {
	return ratingRange{min: 1, max: 4000}
}

func (r ratingRange) possibilities() int {
	return r.max - r.min + 1
}

type ratingsRange map[string]ratingRange

func newRatingsRange() ratingsRange {
	return map[string]ratingRange{
		"x": newRatingRange(),
		"m": newRatingRange(),
		"a": newRatingRange(),
		"s": newRatingRange(),
	}
}

func (rr ratingsRange) possibilities() int {
	return rr["x"].possibilities() *
		rr["m"].possibilities() *
		rr["a"].possibilities() *
		rr["s"].possibilities()
}

var xmas = []string{"x", "m", "a", "s"}

func merge(a, b ratingsRange) ([]ratingsRange, bool) {
	common := map[string]bool{}
	for _, s := range xmas {
		if a[s] == b[s] {
			common[s] = true
		}
	}

	switch len(common) {
	case 0, 1, 2:
		return nil, false
	case 3:
		for _, s := range xmas {
			_, ok := common[s]
			if ok {
				continue
			}

			amin, amax := a[s].min, a[s].max
			bmin, bmax := b[s].min, b[s].max
			if amax < bmin || amin > bmax {
				return nil, false
			}

			min, max := amin, amax
			if bmin < amin {
				min = bmin
			}

			if bmax > amax {
				max = bmax
			}

			c := ratingsRange{}
			maps.Copy(c, a)
			c[s] = ratingRange{min: min, max: max}
			return []ratingsRange{c}, true
		}
		return nil, false
	case 4:
		return []ratingsRange{a}, true
	}

	return nil, false
}

func (rr ratingsRange) String() string {
	return fmt.Sprintf("x:[%d;%d]\tm:[%d;%d]\ta:[%d;%d]\ts:[%d;%d]",
		rr["x"].min, rr["x"].max,
		rr["m"].min, rr["m"].max,
		rr["a"].min, rr["a"].max,
		rr["s"].min, rr["s"].max,
	)
}

type rating map[string]int

func (r rating) value() int {
	ret := 0
	for _, v := range r {
		ret += v
	}
	return ret
}

type workflow struct {
	name         string
	instructions []*ins
}

func (w *workflow) run(r rating) string {
	for _, in := range w.instructions {
		ret := in.run(r)
		if ret != "" {
			return ret
		}
	}
	return ""
}

func newWorkflow(name string, input []string) *workflow {
	w := &workflow{
		name:         name,
		instructions: []*ins{},
	}

	for _, in := range input {
		w.instructions = append(w.instructions, newIns(in))
	}

	return w
}

type ins struct {
	reg, cond   string
	value       int
	dest        string
	asCondition bool
}

func (in *ins) String() string {
	if !in.asCondition {
		return fmt.Sprintf("-> %s", in.dest)
	}
	return fmt.Sprintf("%s %s %d -> %s", in.reg, in.cond, in.value, in.dest)
}

func (in *ins) run(r rating) string {
	if !in.asCondition {
		return in.dest
	}

	reg := r[in.reg]
	if in.cond == ">" && reg > in.value {
		return in.dest
	}

	if in.cond == "<" && reg < in.value {
		return in.dest
	}

	return ""
}

func newIns(input string) *ins {
	ins := &ins{}

	if strings.Contains(input, ":") {
		parts := strings.Split(input, ":")
		ins.reg = string(parts[0][0])
		ins.cond = string(parts[0][1])
		ins.value = aoc.MustGet(strconv.Atoi(parts[0][2:]))
		ins.dest = parts[1]
		ins.asCondition = true
	} else {
		ins.dest = input
	}

	return ins
}

type day struct {
	workflows    map[string]*workflow
	ratings      []rating
	validRatings []ratingsRange
}

func (d *day) Expect(part int, test bool) string {
	return aoc.NewResult(
		"19114", "263678", "167409079868000", "125455345557345",
	).Expect(part, test)
}

func (d *day) Solve(r io.Reader, part int) (string, error) {
	d.parseInput(r, part)

	ret := 0
	if part == 1 {
		for _, rat := range d.ratings {
			w := d.workflows["in"]
			for {
				result := w.run(rat)
				if result == "A" {
					ret += rat.value()
					break
				}

				if result == "R" {
					break
				}

				w = d.workflows[result]
			}
		}
	} else {
		d.explore("in", newRatingsRange(), 1)
		for {
			ratings := []ratingsRange{}
			copy(ratings, d.validRatings)
			allGood := true

			for _, r := range ratings {
				if d.addRating(r) {
					allGood = false
				}
			}

			if allGood {
				break
			}
		}

		ret = 0
		for _, vr := range d.validRatings {
			ret += vr.possibilities()
		}
	}

	return strconv.Itoa(ret), nil
}

func (d *day) addRating(rr ratingsRange) bool {
	if len(d.validRatings) == 0 {
		d.validRatings = []ratingsRange{}
	}

	if len(d.validRatings) == 0 {
		d.validRatings = append(d.validRatings, rr)
		return true
	}

	nr := []ratingsRange{}
	wasMerged := false
	for _, vr := range d.validRatings {
		merged, ok := merge(vr, rr)
		if ok {
			nr = append(nr, merged...)
			wasMerged = true
		} else {
			nr = append(nr, vr)
		}
	}

	if !wasMerged {
		nr = append(nr, rr)
	}

	d.validRatings = nr
	return wasMerged
}

func (d *day) explore(name string, irr ratingsRange, depth int) {
	rr := ratingsRange{}
	maps.Copy(rr, irr)
	if name == "R" {
		return
	}

	if name == "A" {
		d.addRating(rr)
		return
	}

	w := d.workflows[name]
	nr := ratingsRange{}
	maps.Copy(nr, rr)

	for _, ins := range w.instructions {
		if !ins.asCondition {
			d.explore(ins.dest, nr, depth+1)
			return
		}

		r := nr[ins.reg]
		switch ins.cond {
		case "<":
			if r.min >= ins.value {
				continue
			}
			nr[ins.reg] = ratingRange{min: r.min, max: ins.value - 1}
			d.explore(ins.dest, nr, depth+1)
			nr[ins.reg] = ratingRange{min: ins.value, max: r.max}
		case ">":
			if r.max <= ins.value {
				continue
			}
			nr[ins.reg] = ratingRange{min: ins.value + 1, max: r.max}
			d.explore(ins.dest, nr, depth+1)
			nr[ins.reg] = ratingRange{min: r.min, max: ins.value}
		case "":
			d.explore(ins.dest, nr, depth+1)
		default:
			panic(fmt.Sprintf("invalid condition: %s", ins))
		}
	}
}

func (d *day) parseInput(r io.Reader, part int) {
	scanner := bufio.NewScanner(r)

	d.workflows = map[string]*workflow{}
	d.ratings = []rating{}
	d.validRatings = []ratingsRange{}
	workflowsDone := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			workflowsDone = true
			continue
		}

		if !workflowsDone {
			start := strings.Index(line, "{")
			name := line[0:start]
			instructions := line[start+1:]
			instructions = instructions[:len(instructions)-1]
			ins := strings.Split(instructions, ",")
			d.workflows[name] = newWorkflow(name, ins)
		} else {
			r := rating{}
			line = line[1 : len(line)-1]
			parts := strings.Split(line, ",")
			for _, p := range parts {
				k := string(p[0])
				v := aoc.MustGet(strconv.Atoi(p[2:]))
				r[k] = v
			}
			d.ratings = append(d.ratings, r)
		}
	}
}
