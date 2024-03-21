package internal

import (
	"encoding/binary"
	"errors"
	"net"
)

const (
	peerSize   = 6
	portOffset = 4
)

// Peer represents connection information for a peer.
type Peer struct {
	IP   net.IP
	Port uint16
}

// DecodePeers parses peer IP addresses and ports from binary data.
func DecodePeers(peerData []byte) ([]Peer, error) {
	if peerData == nil {
		return nil, errors.New("input data is nil")
	}

	if len(peerData)%peerSize != 0 {
		return nil, errors.New("malformed peer data: incorrect size")
	}

	numPeers := len(peerData) / peerSize
	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		ip := net.IPv4(peerData[offset], peerData[offset+1], peerData[offset+2], peerData[offset+3])
		port := binary.BigEndian.Uint16(peerData[offset+portOffset : offset+peerSize])

		peer := Peer{
			IP:   ip,
			Port: port,
		}
		peers[i] = peer
	}

	return peers, nil
}
