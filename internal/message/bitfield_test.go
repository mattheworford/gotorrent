package message

import (
	"testing"
)

func TestBitfield_HasPiece(t *testing.T) {
	bf := Bitfield{0b10000000, 0b01000000, 0b00100000}

	tests := []struct {
		index          int
		expectedResult bool
	}{
		{0, true},
		{1, false},
		{8, false},
		{9, true},
		{15, false},
		{-1, false},
		{24, false},
	}

	for _, test := range tests {
		result := bf.HasPiece(test.index)
		if result != test.expectedResult {
			t.Errorf("Expected HasPiece(%d) to be %v, but got %v", test.index, test.expectedResult, result)
		}
	}
}

func TestBitfield_SetPiece(t *testing.T) {
	bf := Bitfield{0x00, 0x00, 0x00}

	tests := []struct {
		index            int
		expectedBitfield Bitfield
		shouldError      bool
	}{
		{0, Bitfield{0x80, 0x00, 0x00}, false},
		{1, Bitfield{0x40, 0x00, 0x00}, false},
		{7, Bitfield{0x01, 0x00, 0x00}, false},
		{-1, nil, true},
		{24, nil, true},
	}

	for _, test := range tests {
		result, err := bf.SetPiece(test.index)
		if (err != nil) != test.shouldError {
			t.Errorf("Unexpected error state for SetPiece(%d): %v", test.index, err)
			continue
		}
		if err == nil && !equalBitfields(result, test.expectedBitfield) {
			t.Errorf("Expected SetPiece(%d) to return %v, but got %v", test.index, test.expectedBitfield, result)
		}
	}
}

func equalBitfields(bf1, bf2 Bitfield) bool {
	if len(bf1) != len(bf2) {
		return false
	}
	for i := range bf1 {
		if bf1[i] != bf2[i] {
			return false
		}
	}
	return true
}
