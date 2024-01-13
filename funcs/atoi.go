package funcs

import "fmt"

// Atoi converts a string to an integer. It handles positive and negative integers.
// It returns the converted integer and an error if the conversion fails.
func Atoi(s string) (int, error) {
	// Convert the string to a rune slice for better character handling
	st := []rune(s)

	// Initialize variables for the integer value, and sign
	n := 0
	sign := 1

	// Check if the string has more than one character
	if len(st) > 1 {
		// Handle positive and negative signs
		if st[0] == '+' {
			st = st[1:]
		} else if st[0] == '-' {
			sign = -1
			st = st[1:]
		}
	}

	// Iterate through the rune slice to convert characters to integer
	for i := 0; i < len(st); i++ {
		// Check for invalid characters, spaces are not allowed
		if st[i] < '0' || st[i] > '9' || st[i] == ' ' {
			return 0, fmt.Errorf("invalid argument for Atoi")
		}

		// Update the integer value
		n = (n*10 + int(st[i]) - '0')
	}

	// Return the final converted integer with the correct sign
	return sign * n, nil
}
