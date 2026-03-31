package fracidx

import (
	"sort"
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	s := Start()
	if s == "" {
		t.Fatal("Start() returned empty string")
	}
}

func TestEnd(t *testing.T) {
	s := Start()
	e := End(s)
	if e <= s {
		t.Fatalf("End(%q) = %q, expected > %q", s, e, s)
	}
}

func TestGenerateBetween_BothEmpty(t *testing.T) {
	result := GenerateBetween("", "")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestGenerateBetween_BeforeFirst(t *testing.T) {
	b := Start()
	a := GenerateBetween("", b)
	if a >= b {
		t.Fatalf("expected %q < %q", a, b)
	}
}

func TestGenerateBetween_AfterLast(t *testing.T) {
	a := Start()
	b := GenerateBetween(a, "")
	if b <= a {
		t.Fatalf("expected %q > %q", b, a)
	}
}

func TestGenerateBetween_Midpoint(t *testing.T) {
	a := "A"
	b := "Z"
	mid := GenerateBetween(a, b)
	if mid <= a || mid >= b {
		t.Fatalf("expected %q < %q < %q", a, mid, b)
	}
}

func TestGenerateBetween_Adjacent(t *testing.T) {
	a := "a"
	b := "b"
	mid := GenerateBetween(a, b)
	if mid <= a || mid >= b {
		t.Fatalf("expected %q < %q < %q", a, mid, b)
	}
}

func TestGenerateBetween_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when a >= b")
		}
	}()
	GenerateBetween("z", "a")
}

func TestLexOrder_Sequential(t *testing.T) {
	keys := make([]string, 20)
	keys[0] = Start()
	for i := 1; i < len(keys); i++ {
		keys[i] = End(keys[i-1])
	}
	if !sort.StringsAreSorted(keys) {
		t.Fatalf("keys are not sorted: %v", keys)
	}
}

func TestLexOrder_InsertBetween(t *testing.T) {
	a := "M"
	b := "N"
	keys := []string{a}
	for i := 0; i < 10; i++ {
		mid := GenerateBetween(a, b)
		keys = append(keys, mid)
		b = mid
	}
	keys = append(keys, "N")
	sort.Strings(keys)
	for i := 1; i < len(keys); i++ {
		if keys[i] <= keys[i-1] {
			t.Fatalf("duplicate or unsorted at index %d: %v", i, keys)
		}
	}
}

func TestLexOrder_PrependMany(t *testing.T) {
	keys := make([]string, 15)
	keys[14] = Start()
	for i := 13; i >= 0; i-- {
		keys[i] = GenerateBetween("", keys[i+1])
	}
	if !sort.StringsAreSorted(keys) {
		t.Fatalf("prepended keys not sorted: %v", keys)
	}
}

func TestRebalance(t *testing.T) {
	keys := Rebalance(5)
	if len(keys) != 5 {
		t.Fatalf("expected 5 keys, got %d", len(keys))
	}
	if !sort.StringsAreSorted(keys) {
		t.Fatalf("rebalanced keys not sorted: %v", keys)
	}
}

func TestRebalance_Zero(t *testing.T) {
	keys := Rebalance(0)
	if keys != nil {
		t.Fatalf("expected nil, got %v", keys)
	}
}

func TestNeedsRebalance(t *testing.T) {
	short := []string{"A", "M", "Z"}
	if NeedsRebalance(short) {
		t.Fatal("short keys should not need rebalance")
	}

	long := make([]byte, 51)
	for i := range long {
		long[i] = 'M'
	}
	if !NeedsRebalance([]string{string(long)}) {
		t.Fatal("long key should need rebalance")
	}
}

func TestGenerateBetween_VeryClose(t *testing.T) {
	// Force midpoint fallback path by using very close keys
	a := "Ma"
	b := "Mb"
	mid := GenerateBetween(a, b)
	if mid <= a || mid >= b {
		t.Fatalf("expected %q < %q < %q", a, mid, b)
	}
}

func TestGenerateBetween_EqualPrefix(t *testing.T) {
	// Keys with equal prefix chars to exercise the carry loop
	a := "aaa"
	b := "aab"
	mid := GenerateBetween(a, b)
	if mid <= a || mid >= b {
		t.Fatalf("expected %q < %q < %q", a, mid, b)
	}
}

func TestGenerateBetween_LongEqualPrefix(t *testing.T) {
	// All equal chars — forces fallback
	a := "MMMM"
	b := "MMMMM"
	mid := GenerateBetween(a, b)
	if mid <= a || mid >= b {
		t.Fatalf("expected %q < %q < %q", a, mid, b)
	}
}

func TestGenerateBetween_AdjacentAllPositions(t *testing.T) {
	// Adjacent at every position to force midpoint fallback path
	a := "aa"
	b := "ab"
	mid := GenerateBetween(a, b)
	if mid <= a || mid >= b {
		t.Fatalf("expected %q < %q < %q", a, mid, b)
	}
	// Single char adjacent
	a2 := "b"
	b2 := "c"
	mid2 := GenerateBetween(a2, b2)
	if mid2 <= a2 || mid2 >= b2 {
		t.Fatalf("expected %q < %q < %q", a2, mid2, b2)
	}
}

func TestRebalance_LargeN(t *testing.T) {
	// Rebalance with n that fits in char range (charEnd - charStart = 74)
	keys := Rebalance(50)
	if len(keys) != 50 {
		t.Fatalf("expected 50 keys, got %d", len(keys))
	}
	if !sort.StringsAreSorted(keys) {
		t.Fatalf("rebalanced keys not sorted")
	}
	// Verify uniqueness
	seen := make(map[string]bool)
	for _, k := range keys {
		if seen[k] {
			t.Fatalf("duplicate key in rebalance: %q", k)
		}
		seen[k] = true
	}
}

func TestRebalance_Single(t *testing.T) {
	keys := Rebalance(1)
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
}

func TestNeedsRebalance_EmptySlice(t *testing.T) {
	if NeedsRebalance(nil) {
		t.Fatal("nil slice should not need rebalance")
	}
	if NeedsRebalance([]string{}) {
		t.Fatal("empty slice should not need rebalance")
	}
}

func TestGenerateBetween_DeepNesting(t *testing.T) {
	// Insert 50 items between two very close keys to test key growth
	a := "M"
	b := "N"
	for i := 0; i < 50; i++ {
		mid := GenerateBetween(a, b)
		if mid <= a || mid >= b {
			t.Fatalf("iteration %d: expected %q < %q < %q", i, a, mid, b)
		}
		b = mid
	}
}

func TestGenerateBefore_MinChars(t *testing.T) {
	// Generate before the smallest possible key
	small := string(charStart + 1)
	before := GenerateBetween("", small)
	if before >= small {
		t.Fatalf("expected %q < %q", before, small)
	}
}

func TestGenerateAfter_MaxChars(t *testing.T) {
	// Generate after a key near the maximum character
	big := string(charEnd - 1)
	after := GenerateBetween(big, "")
	if after <= big {
		t.Fatalf("expected %q > %q", after, big)
	}
}

func TestGenerateBetween_Concurrent(t *testing.T) {
	var wg sync.WaitGroup
	results := make([]string, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = GenerateBetween("A", "z")
		}(i)
	}
	wg.Wait()
	for i, r := range results {
		if r == "" {
			t.Fatalf("empty result at index %d", i)
		}
		if r <= "A" || r >= "z" {
			t.Fatalf("result %q at index %d out of range", r, i)
		}
	}
}
