package goBoom

func reverse(input string) string {
	runes := make([]rune, len(input))
	n := 0
	for _, r := range input {
		runes[n] = r
		n++
	}
	runes = runes[0:n]

	// Reverse
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}

	// Convert back to UTF-8.
	return string(runes)
}
