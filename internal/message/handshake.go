package message

import (
	"fmt"
	"io"
)

type Handshake struct {
	ProtocolString string
	InfoHash       [20]byte
	PeerID         [20]byte
}

func NewHandshake(infoHash [20]byte, peerID [20]byte) *Handshake {
	return &Handshake{
		ProtocolString: "BitTorrent protocol",
		InfoHash:       infoHash,
		PeerID:         peerID,
	}
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.ProtocolString)+49)
	buf[0] = byte(len(h.ProtocolString))
	curr := 1
	curr += copy(buf[curr:], h.ProtocolString)
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

func ParseHandshake(r io.Reader) (*Handshake, error) {

	protocolStringLenBuf := make([]byte, 1)
	_, err := r.Read(protocolStringLenBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read ProtocolString length: %v", err)
	}
	protocolStringLen := int(protocolStringLenBuf[0])

	protocolStringBuf := make([]byte, protocolStringLen)
	_, err = r.Read(protocolStringBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read ProtocolString: %v", err)
	}
	protocolString := string(protocolStringBuf)

	reservedBuf := make([]byte, 8)
	_, err = r.Read(reservedBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read reserved bytes: %v", err)
	}

	infoHashBuf := make([]byte, 20)
	_, err = r.Read(infoHashBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read InfoHash: %v", err)
	}
	var infoHash [20]byte
	copy(infoHash[:], infoHashBuf)

	peerIDBuf := make([]byte, 20)
	_, err = r.Read(peerIDBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read PeerID: %v", err)
	}
	var peerID [20]byte
	copy(peerID[:], peerIDBuf)

	return &Handshake{
		ProtocolString: protocolString,
		InfoHash:       infoHash,
		PeerID:         peerID,
	}, nil
}
