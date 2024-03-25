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
