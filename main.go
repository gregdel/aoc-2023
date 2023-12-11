package main

import (
	"flag"
	"fmt"
	"os"

	aoc "github.com/gregdel/aoc2023/lib"

	_ "github.com/gregdel/aoc2023/day1"
	_ "github.com/gregdel/aoc2023/day2"
	_ "github.com/gregdel/aoc2023/day3"
	_ "github.com/gregdel/aoc2023/day4"
	_ "github.com/gregdel/aoc2023/day5"
	_ "github.com/gregdel/aoc2023/day6"
	_ "github.com/gregdel/aoc2023/day7"
	_ "github.com/gregdel/aoc2023/day8"
	_ "github.com/gregdel/aoc2023/day9"

	_ "github.com/gregdel/aoc2023/day10"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var day, part int
	var all bool
	var test *bool

	flag.IntVar(&day, "day", 0, "day to run")
	flag.IntVar(&part, "part", 0, "part to run")
	flag.BoolVar(&all, "all", false, "run test and non test")
	flag.BoolFunc("test", "run on test input", func(s string) error {
		var testValue bool
		if s == "true" || s == "1" {
			testValue = true
		}
		test = &testValue
		return nil
	})
	flag.Parse()

	tests := []bool{true, false}
	if test != nil {
		tests = []bool{*test}
	}

	days := []int{}
	if all {
		days = aoc.AllDays()
	} else if day != 0 {
		days = []int{day}
	}

	if len(days) == 0 {
		return fmt.Errorf("Missing day")
	}

	parts := []int{1, 2}
	if part != 0 {
		parts = []int{part}
	}

	for _, day := range days {
		for _, part := range parts {
			for _, test := range tests {
				result, err := aoc.Run(day, part, test)
				if err != nil {
					return err
				}

				result.Show()
			}
		}
	}

	return nil
}
