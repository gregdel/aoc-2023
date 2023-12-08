package aoc

import "golang.org/x/exp/constraints"

// GreatestCommonDivisor returns the greatest common divisor of a and b.
func GreatestCommonDivisor[T constraints.Integer](a, b T) T {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LeastCommonMultiple returns the least common multiple of integers
func LeastCommonMultiple[T constraints.Integer](a, b T, integers ...T) T {
	result := a * b / GreatestCommonDivisor(a, b)

	for i := 0; i < len(integers); i++ {
		result = LeastCommonMultiple(result, integers[i])
	}

	return result
}
