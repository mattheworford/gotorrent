package message

import "errors"

// Bitfield represents the pieces that a peer has
type Bitfield []byte

// HasPiece tells if a bitfield has a particular index set
func (bf Bitfield) HasPiece(index int) bool {
    if index < 0 || index >= len(bf)*8 {
        return false
    }
    byteIndex := index / 8
    offset := index % 8
    return bf[byteIndex]>>(7-offset)&1 != 0
}

// SetPiece sets a bit in the bitfield and returns a new bitfield
func (bf Bitfield) SetPiece(index int) (Bitfield, error) {
    if index < 0 || index >= len(bf)*8 {
        return nil, errors.New("index out of range")
    }
    byteIndex := index / 8
    offset := index % 8
    copyOfBf := make(Bitfield, len(bf))
    copy(copyOfBf, bf)
    copyOfBf[byteIndex] |= 1 << (7 - offset)
    return copyOfBf, nil
}
