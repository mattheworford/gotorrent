package status

import (
	"fmt"

	"github.com/mattheworford/gotorrent/internal/message"
	"github.com/mattheworford/gotorrent/internal/peer"
)

type CurrentStatus struct {
	Index      int
	Client     *peer.Client
	Buf        []byte
	Downloaded int
	Requested  int
	Backlog    int
}

// Update updates the current status based on a parsed piece.
func (cs *CurrentStatus) Update(piece *message.Piece) error {
	if piece.Index != cs.Index {
		return fmt.Errorf("expected index %d, but got index %d", cs.Index, piece.Index)
	}

	if piece.Offset >= len(cs.Buf) {
		return fmt.Errorf("begin offset exceeds buffer length. Offset: %d, Buffer Length: %d", piece.Offset, len(cs.Buf))
	}

	if piece.Offset+len(piece.Data) > len(cs.Buf) {
		return fmt.Errorf("data exceeds buffer capacity. Offset: %d, Data Length: %d, Buffer Length: %d", piece.Offset, len(piece.Data), len(cs.Buf))
	}

	copy(cs.Buf[piece.Offset:], piece.Data)

	cs.Downloaded += len(piece.Data)
	cs.Backlog--
	return nil
}

func (cs *CurrentStatus) readMessage() error {
	msg, err := cs.Client.Read() // this call blocks
	if err != nil {
		return err
	}
	switch msg.Type {
	case message.UnchokeMessage:
		cs.Client.Choked = false
	case message.ChokeMessage:
		cs.Client.Choked = true
	case message.HaveMessage:
		index, err := message.ParseHaveMessage(msg)
		if err != nil {
			return err
		}
		cs.Client.Bitfield.SetPiece(index)
	case message.PieceMessage:
		piece, err := message.ParsePieceMessage(msg)
		if err != nil {
			return err
		}
		if err := cs.Update(piece); err != nil {
			return err
		}
	}
	return nil
}
