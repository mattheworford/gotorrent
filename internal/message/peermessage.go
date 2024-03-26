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

type PeerMessageType uint8

const (
	MsgChoke         PeerMessageType = iota // 0
	MsgUnchoke                              // 1
	MsgInterested                           // 2
	MsgNotInterested                        // 3
	MsgHave                                 // 4
	MsgBitfield                             // 5
	MsgRequest                              // 6
	MsgPiece                                // 7
	MsgCancel                               // 8
)

type PeerMessage struct {
	Type    PeerMessageType
	Payload []byte
}

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

func ParsePeerMessage(r io.Reader) (*PeerMessage, error) {
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

func ParsePiece(index int, buf []byte, msg *PeerMessage) (int, error) {
	if msg.Type != MsgPiece {
		return 0, fmt.Errorf("Expected PIECE (Type %d), got Type %d", MsgPiece, msg.Type)
	}
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	copy(buf[begin:], data)
	return len(data), nil
}

// ParseHave parses a HAVE message
func ParseHave(msg *PeerMessage) (int, error) {
	if msg.Type != MsgHave {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", MsgHave, msg.Type)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}
