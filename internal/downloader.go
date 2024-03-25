package internal

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"net"

	"github.com/mattheworford/gotorrent/internal/handshake"
	"github.com/mattheworford/gotorrent/internal/peer"
)

type Downloader struct {
	Config DownloadConfig
}

// DownloadConfig holds configuration for downloading a torrent from a list of peers
type DownloadConfig struct {
	Peers       []peer.ConnectionInfo
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (downloader *Downloader) startWorker(peer peer.ConnectionInfo, workQueue chan *pieceWork, results chan *pieceResult) {
	config := downloader.Config
	c, err := client.New(peer, config.PeerID, config.InfoHash)
	if err != nil {
		log.Printf("Could not handshake with %s. Disconnecting\n", peer.IP)
		return
	}
	defer c.Conn.Close()
	log.Printf("Completed handshake with %s\n", peer.IP)

	c.SendUnchoke()
	c.SendInterested()

	for pw := range workQueue {
		if !c.Bitfield.HasPiece(pw.index) {
			workQueue <- pw // Put piece back on the queue
			continue
		}

		// Download the piece
		buf, err := attemptDownloadPiece(c, pw)
		if err != nil {
			log.Println("Exiting", err)
			workQueue <- pw // Put piece back on the queue
			return
		}

		err = checkIntegrity(pw, buf)
		if err != nil {
			log.Printf("Piece #%d failed integrity check\n", pw.index)
			workQueue <- pw // Put piece back on the queue
			continue
		}

		c.SendHave(pw.index)
		results <- &pieceResult{pw.index, buf}
	}
}

func CompleteHandshake(peer peer.ConnectionInfo, torrentData TorrentData, conn net.Conn) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return fmt.Errorf("failed to generate peer ID")
	}
	handshake := handshake.NewHandshake(
		torrentData.InfoHash,
		peerID,
	)
	_, err = conn.Write(handshake.Serialize())
	if err != nil {
		return fmt.Errorf("failed to send handshake to peer %s:%d: %v", peer.IP, peer.Port, err)
	}

	res, err := ParseHandshake(conn)
	if err != nil {
		return fmt.Errorf("failed to receive handshake response from peer %s:%d: %v", peer.IP, peer.Port, err)
	}
	if !bytes.Equal(res.InfoHash[:], torrentData.InfoHash[:]) {
		return fmt.Errorf("expected infohash %x but got %x", res.InfoHash, torrentData.InfoHash)
	}
	return nil
}

// DownloadPiece simulates downloading a piece from a peer
func DownloadPiece(peer peer.ConnectionInfo, pieceNumber int, torrentData TorrentData) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peer.IP, peer.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s:%d: %v", peer.IP, peer.Port, err)
	}
	defer conn.Close()

	err = CompleteHandshake(peer, torrentData, conn)
	if err != nil {
		return err
	}

	// Step 3: Exchange messages to download pieces
	request := fmt.Sprintf("Requesting piece #%d", pieceNumber)
	_, err = conn.Write([]byte(request))
	if err != nil {
		return fmt.Errorf("failed to send request to peer %s:%d: %v", peer.IP, peer.Port, err)
	}

	// Simulate receiving the piece data from the peer
	var receivedData [1024]byte // Simulated piece data
	_, err = conn.Read(receivedData[:])
	if err != nil {
		return fmt.Errorf("failed to receive piece data from peer %s:%d: %v", peer.IP, peer.Port, err)
	}

	// Process received piece data...

	return nil
}
