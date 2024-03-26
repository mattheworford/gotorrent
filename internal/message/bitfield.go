package message

import "errors"

const BitsPerByte = 8

// Bitfield represents the pieces that a peer has.
type Bitfield []byte

// HasPiece tells if the bit at a given index is set.
func (bf Bitfield) HasPiece(index int) bool {
	if index < 0 || index >= len(bf)*BitsPerByte {
		return false
	}
	byteIndex := index / BitsPerByte
	offset := index % BitsPerByte
	return bf[byteIndex]>>(7-offset)&1 != 0
}

// SetPiece returns a copy of the bitfield with the bit at the given index set.
func (bf Bitfield) SetPiece(index int) (Bitfield, error) {
	if index < 0 || index >= len(bf)*BitsPerByte {
		return nil, errors.New("index out of range")
	}
	byteIndex := index / BitsPerByte
	offset := index % BitsPerByte
	copyOfBf := make(Bitfield, len(bf))
	copy(copyOfBf, bf)
	copyOfBf[byteIndex] |= 1 << (7 - offset)
	return copyOfBf, nil
}
