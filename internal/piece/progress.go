package piece

import (
	"github.com/mattheworford/gotorrent/internal/message"
	"github.com/mattheworford/gotorrent/internal/peer"
)

type Progress struct {
	Index      int
	Client     *peer.Client
	Buf        []byte
	Downloaded int
	Requested  int
	Backlog    int
}

func (state *Progress) readMessage() error {
	msg, err := state.Client.Read() // this call blocks
	if err != nil {
		return err
	}
	switch msg.Type {
	case message.MsgUnchoke:
		state.Client.Choked = false
	case message.MsgChoke:
		state.Client.Choked = true
	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.Client.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := message.ParsePiece(state.Index, state.Buf, msg)
		if err != nil {
			return err
		}
		state.Downloaded += n
		state.Backlog--
	}
	return nil
}
