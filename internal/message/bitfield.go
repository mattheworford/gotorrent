package message

import "errors"

// Bitfield represents the pieces that a peer has
type Bitfield []byte

const bitsPerByte = 8

// HasPiece tells if a bitfield has a particular index set
func (bf Bitfield) HasPiece(index int) bool {
	if index < 0 || index >= len(bf)*bitsPerByte {
		return false
	}
	byteIndex := index / bitsPerByte
	offset := index % bitsPerByte
	return bf[byteIndex]>>(7-offset)&1 != 0
}

// SetPiece sets a bit in the bitfield and returns a new bitfield
func (bf Bitfield) SetPiece(index int) (Bitfield, error) {
	if index < 0 || index >= len(bf)*bitsPerByte {
		return nil, errors.New("index out of range")
	}
	byteIndex := index / bitsPerByte
	offset := index % bitsPerByte
	copyOfBf := make(Bitfield, len(bf))
	copy(copyOfBf, bf)
	copyOfBf[byteIndex] |= 1 << (7 - offset)
	return copyOfBf, nil
}
