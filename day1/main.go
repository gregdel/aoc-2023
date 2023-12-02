package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var words = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func transformLine(input string) string {
	line := input
	for i := 1; i < len(words); i++ {
		o := words[i]
		n := words[i] + strconv.Itoa(i) + words[i]
		line = strings.ReplaceAll(line, o, n)
	}

	return line
}

func solve(reader io.Reader, transform bool) (int, error) {
	result := 0
	var err error

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var first, last int
		var foundFirst bool

		line := scanner.Text()
		if transform {
			line = transformLine(line)
		}

		for _, r := range line {
			if !unicode.IsDigit(r) {
				continue
			}

			if !foundFirst {
				first, err = strconv.Atoi(string(r))
				if err != nil {
					return 0, err
				}
				last = first
				foundFirst = true
			}

			last, err = strconv.Atoi(string(r))
			if err != nil {
				return 0, err
			}
		}

		value := first*10 + last
		result += value
	}

	return result, nil
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

	result, err := solve(file, true)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %d\n", result)

	return nil
}
