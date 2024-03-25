package message

import (
	"encoding/binary"
	"errors"
	"io"
)

// PeerMessageType represents the type of peer message
type PeerMessageType uint8

const (
	MaxPayloadLength             = 1 << 14
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

// PeerMessage stores ID and payload of a message.
type PeerMessage struct {
	Type    PeerMessageType
	Payload []byte
}

// Serialize serializes a message into a buffer of the form.
func (m *PeerMessage) Serialize() ([]byte, error) {
	if m == nil {
		return make([]byte, 4), nil
	}

	payloadLength := len(m.Payload)
	if payloadLength > MaxPayloadLength {
		return nil, errors.New("payload length exceeds maximum allowed")
	}

	length := uint32(payloadLength + 1)
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.Type)
	copy(buf[5:], m.Payload)
	return buf, nil
}

// ParsePeerMessage reads a message from a stream
func ParsePeerMessage(r io.Reader) (*PeerMessage, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		return nil, nil
	}

	if length > MaxPayloadLength {
		return nil, errors.New("message length exceeds maximum allowed")
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := PeerMessage{
		Type:    PeerMessageType(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}
