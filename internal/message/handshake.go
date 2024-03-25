package message

import (
	"fmt"
	"io"
)

const (
	InfoHashLength  = 20
	PeerIDLength    = 20
	ReservedBufSize = 8
)

type Handshake struct {
	ProtocolString string
	InfoHash       [InfoHashLength]byte
	PeerID         [PeerIDLength]byte
}

// NewHandshake creates a new Handshake message.
func NewHandshake(infoHash [InfoHashLength]byte, peerID [PeerIDLength]byte) *Handshake {
	return &Handshake{
		ProtocolString: "BitTorrent protocol",
		InfoHash:       infoHash,
		PeerID:         peerID,
	}
}

// Serialize serializes the Handshake message into a byte slice.
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.ProtocolString)+49)
	buf[0] = byte(len(h.ProtocolString))
	curr := 1
	curr += copy(buf[curr:], h.ProtocolString)
	curr += copy(buf[curr:], make([]byte, ReservedBufSize))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

// ParseHandshake parses a Handshake message from an io.Reader.
func ParseHandshake(r io.Reader) (*Handshake, error) {
	protocolStringLenBuf := make([]byte, 1)
	_, err := r.Read(protocolStringLenBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read ProtocolString length: %w", err)
	}
	protocolStringLen := int(protocolStringLenBuf[0])

	protocolStringBuf := make([]byte, protocolStringLen)
	_, err = r.Read(protocolStringBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read ProtocolString: %w", err)
	}
	protocolString := string(protocolStringBuf)

	reservedBuf := make([]byte, ReservedBufSize)
	_, err = r.Read(reservedBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read reserved bytes: %w", err)
	}

	infoHashBuf := make([]byte, InfoHashLength)
	_, err = r.Read(infoHashBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read InfoHash: %w", err)
	}
	var infoHash [InfoHashLength]byte
	copy(infoHash[:], infoHashBuf)

	peerIDBuf := make([]byte, PeerIDLength)
	_, err = r.Read(peerIDBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read PeerID: %w", err)
	}
	var peerID [PeerIDLength]byte
	copy(peerID[:], peerIDBuf)

	return &Handshake{
		ProtocolString: protocolString,
		InfoHash:       infoHash,
		PeerID:         peerID,
	}, nil
}
