package peer

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/mattheworford/gotorrent/internal/message"
)

const (
	peerSize   = 6
	portOffset = 4
)

// ConnectionInfo represents connection information for a peer.
type ConnectionInfo struct {
	IP   net.IP
	Port uint16
}

// Client represents a peer client.
type Client struct {
	Conn           net.Conn
	Choked         bool
	Bitfield       message.Bitfield
	ConnectionInfo ConnectionInfo
	InfoHash       [20]byte
	PeerID         [20]byte
}

// DecodeConnectionInfo parses peer IP addresses and ports from binary data.
func DecodeConnectionInfo(peerData []byte) ([]ConnectionInfo, error) {
	if peerData == nil {
		return nil, errors.New("input data is nil")
	}

	if len(peerData)%peerSize != 0 {
		return nil, errors.New("malformed connection data: incorrect size")
	}

	numPeers := len(peerData) / peerSize
	peers := make([]ConnectionInfo, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		ip := net.IPv4(peerData[offset], peerData[offset+1], peerData[offset+2], peerData[offset+3])
		port := binary.BigEndian.Uint16(peerData[offset+portOffset : offset+peerSize])

		connectionInfo := ConnectionInfo{
			IP:   ip,
			Port: port,
		}
		peers[i] = connectionInfo
	}

	return peers, nil
}
