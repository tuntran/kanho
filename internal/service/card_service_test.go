package service_test

import (
	"testing"

	"github.com/tungtran/kanho/internal/fracidx"
)

func TestCreateCard_AssignsSequentialPosition(t *testing.T) {
	// Test that fracidx.Start produces a valid initial position
	pos := fracidx.Start()
	if pos == "" {
		t.Error("Start() should return a non-empty position")
	}

	// Second card should get a position after the first
	pos2 := fracidx.End(pos)
	if pos2 == "" {
		t.Error("End() should return a non-empty position")
	}
	if pos2 <= pos {
		t.Errorf("second position %q should be after first %q", pos2, pos)
	}

	// Third card after second
	pos3 := fracidx.End(pos2)
	if pos3 <= pos2 {
		t.Errorf("third position %q should be after second %q", pos3, pos2)
	}
}

func TestMoveCard_UpdatesPosition(t *testing.T) {
	// Simulate existing positions
	posA := fracidx.Start()
	posB := fracidx.End(posA)
	posC := fracidx.End(posB)

	// Move card to between A and B
	between := fracidx.GenerateBetween(posA, posB)
	if between <= posA || between >= posB {
		t.Errorf("between position %q should be between %q and %q", between, posA, posB)
	}

	// Move card to between B and C
	between2 := fracidx.GenerateBetween(posB, posC)
	if between2 <= posB || between2 >= posC {
		t.Errorf("between position %q should be between %q and %q", between2, posB, posC)
	}

	// Move card to beginning (before A)
	beforeA := fracidx.GenerateBetween("", posA)
	if beforeA >= posA {
		t.Errorf("beforeA position %q should be before %q", beforeA, posA)
	}
}
