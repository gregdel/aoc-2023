package main

import (
	"flag"
	"fmt"
	"os"

	aoc "github.com/gregdel/aoc2023/lib"

	_ "github.com/gregdel/aoc2023/day01"
	_ "github.com/gregdel/aoc2023/day02"
	_ "github.com/gregdel/aoc2023/day03"
	_ "github.com/gregdel/aoc2023/day04"
	_ "github.com/gregdel/aoc2023/day05"
	_ "github.com/gregdel/aoc2023/day06"
	_ "github.com/gregdel/aoc2023/day07"
	_ "github.com/gregdel/aoc2023/day08"
	_ "github.com/gregdel/aoc2023/day09"
	_ "github.com/gregdel/aoc2023/day10"
	_ "github.com/gregdel/aoc2023/day11"
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
