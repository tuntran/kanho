package fracidx

import (
	"strings"
)

const (
	// Character range for fractional index keys: digits + uppercase + lowercase
	charStart = '0' // ASCII 48
	charEnd   = 'z' // ASCII 122
	maxKeyLen = 50
)

// midChar returns the midpoint character between a and b.
func midChar(a, b byte) byte {
	return byte((int(a) + int(b)) / 2)
}

// GenerateBetween returns a key lexicographically between a and b.
// If a is empty, generates a key before b (start position).
// If b is empty, generates a key after a (end position).
// Both empty: returns a midpoint key.
func GenerateBetween(a, b string) string {
	if a != "" && b != "" && a >= b {
		panic("fracidx: a must be less than b")
	}

	if a == "" && b == "" {
		return string(midChar(charStart, charEnd))
	}

	if a == "" {
		return generateBefore(b)
	}

	if b == "" {
		return generateAfter(a)
	}

	return generateMidpoint(a, b)
}

// Start returns a key suitable as the first position.
func Start() string {
	return GenerateBetween("", "")
}

// End returns a key after the given last key.
func End(last string) string {
	return GenerateBetween(last, "")
}

// generateBefore returns a key before b.
func generateBefore(b string) string {
	var result strings.Builder
	for i := 0; i < len(b); i++ {
		c := b[i]
		if c > charStart+1 {
			result.WriteByte(midChar(charStart, c))
			return result.String()
		}
		// c is charStart or charStart+1: carry to next position
		result.WriteByte(charStart)
	}
	// All characters were at minimum, append midpoint
	result.WriteByte(midChar(charStart, charEnd))
	return result.String()
}

// generateAfter returns a key after a.
func generateAfter(a string) string {
	var result strings.Builder
	for i := 0; i < len(a); i++ {
		c := a[i]
		if c < charEnd-1 {
			result.WriteByte(midChar(c, charEnd+1))
			return result.String()
		}
		// c is at or near max: carry to next position
		result.WriteByte(c)
	}
	// All characters were at maximum, append midpoint
	result.WriteByte(midChar(charStart, charEnd))
	return result.String()
}

// generateMidpoint returns a key between a and b.
func generateMidpoint(a, b string) string {
	var result strings.Builder
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	for i := 0; i <= maxLen; i++ {
		ca := charStart
		if i < len(a) {
			ca = rune(a[i])
		}
		cb := charEnd + 1
		if i < len(b) {
			cb = rune(b[i])
		}

		if ca+1 < cb {
			result.WriteByte(midChar(byte(ca), byte(cb)))
			return result.String()
		}
		// Characters are adjacent or equal: carry to next level
		result.WriteByte(byte(ca))
	}

	// Fallback: append midpoint character
	result.WriteByte(midChar(charStart, charEnd))
	return result.String()
}

// Rebalance generates evenly-spaced keys for n items.
// Use when keys grow too long (> maxKeyLen).
func Rebalance(n int) []string {
	if n <= 0 {
		return nil
	}

	keys := make([]string, n)
	step := float64(charEnd-charStart) / float64(n+1)
	for i := 0; i < n; i++ {
		pos := charStart + byte(step*float64(i+1))
		keys[i] = string(pos)
	}
	return keys
}

// NeedsRebalance checks if any key exceeds the max length threshold.
func NeedsRebalance(keys []string) bool {
	for _, k := range keys {
		if len(k) > maxKeyLen {
			return true
		}
	}
	return false
}
