package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MaxPayloadLength = 1 << 14
	LengthBufSize    = 4
)

// PeerMessageType represents the value of the type bit in a peer message.
type PeerMessageType uint8

const (
	ChokeMessage         PeerMessageType = iota // 0
	UnchokeMessage                              // 1
	InterestedMessage                           // 2
	NotInterestedMessage                        // 3
	HaveMessage                                 // 4
	BitfieldMessage                             // 5
	RequestMessage                              // 6
	PieceMessage                                // 7
	CancelMessage                               // 8
)

// PeerMessage represents a non-keepalive message sent between peers.
type PeerMessage struct {
	Type    PeerMessageType
	Payload []byte
}

// Serialize serializes a PeerMessage into a byte slice.
func (m *PeerMessage) Serialize() ([]byte, error) {
	if m == nil {
		return make([]byte, LengthBufSize), nil
	}

	payloadLength := len(m.Payload)
	if payloadLength > MaxPayloadLength {
		return nil, fmt.Errorf("payload length %d exceeds maximum allowed", payloadLength)
	}

	length := uint32(payloadLength + 1)
	buf := make([]byte, LengthBufSize+length)
	binary.BigEndian.PutUint32(buf[0:LengthBufSize], length)
	buf[LengthBufSize] = byte(m.Type)
	copy(buf[LengthBufSize+1:], m.Payload)
	return buf, nil
}

// ReadPeerMessage reads a PeerMessage from an io.Reader.
func ReadPeerMessage(r io.Reader) (*PeerMessage, error) {
	lengthBuf := make([]byte, LengthBufSize)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		return nil, nil
	}

	if length > MaxPayloadLength {
		return nil, fmt.Errorf("message length %d exceeds maximum allowed", length)
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	m := PeerMessage{
		Type:    PeerMessageType(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}

// ParseHaveMessage parses a have message and returns the index of the piece indicated.
func ParseHaveMessage(msg *PeerMessage) (int, error) {
	const expectedPayloadLength = 4

	if msg.Type != HaveMessage {
		return 0, fmt.Errorf("expected piece message (type %d), but got type %d", HaveMessage, msg.Type)
	}

	if len(msg.Payload) != expectedPayloadLength {
		return 0, fmt.Errorf("expected payload length %d, but got length %d", expectedPayloadLength, len(msg.Payload))
	}

	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

// ParsePieceMessage parses piece data from a peer message.
func ParsePieceMessage(msg *PeerMessage) (*Piece, error) {
	const minPayloadLength = 8

	if msg.Type != PieceMessage {
		return nil, fmt.Errorf("expected piece message (type %d), but got type %d", PieceMessage, msg.Type)
	}
	if len(msg.Payload) < minPayloadLength {
		return nil, fmt.Errorf("minimum payload length %d, but got length %d", minPayloadLength, len(msg.Payload))
	}

	index := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	offset := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	data := msg.Payload[8:]

	return &Piece{Data: data, Index: index, Offset: offset}, nil
}
