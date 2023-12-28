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

// Min returns the min of two numbers
func Min[T constraints.Integer](a, b T) T {
	if a < b {
		return a
	}

	return b
}

// Max returns the max of two numbers
func Max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}

	return b
}

// Abs returns the absolute value of v
func Abs[T constraints.Integer](v T) T {
	if v < 0 {
		return -v
	}

	return v
}
