package main

import (
	"flag"
	"fmt"
	"os"

	aoc "github.com/gregdel/aoc2023/lib"

	_ "github.com/gregdel/aoc2023/day1"
	_ "github.com/gregdel/aoc2023/day2"
	_ "github.com/gregdel/aoc2023/day4"
	_ "github.com/gregdel/aoc2023/day5"
	_ "github.com/gregdel/aoc2023/day6"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var day, part int
	var test, all bool

	flag.IntVar(&day, "day", 0, "day to run")
	flag.IntVar(&part, "part", 0, "part to run")
	flag.BoolVar(&test, "test", false, "run on test input")
	flag.BoolVar(&all, "all", false, "run test and non test")
	flag.Parse()

	tests := []bool{true, false}
	if !all && test {
		tests = []bool{test}
	}

	if day == 0 {
		return fmt.Errorf("Missing day")
	}

	parts := []int{1, 2}
	if part != 0 {
		parts = []int{part}
	}

	for _, part := range parts {
		for _, test := range tests {
			result, err := aoc.Run(day, part, test)
			if err != nil {
				return err
			}

			result.Show()
		}
	}

	return nil
}
