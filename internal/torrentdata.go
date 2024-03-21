package internal

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/jackpal/bencode-go"
)

// InfoDictionary represents the metadata of a torrent file.
type InfoDictionary struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

// MetainfoFile represents the top-level structure of a torrent file.
type MetainfoFile struct {
	Announce string         `bencode:"announce"`
	Info     InfoDictionary `bencode:"info"`
}

// TorrentData represents the processed data extracted from a torrent file.
type TorrentData struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

// Open parses the torrent file at the specified path and returns its metadata.
func Open(path string) (*MetainfoFile, error) {
	if path == "" {
		return nil, errors.New("torrentdata: path cannot be empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("torrentdata: failed to open file: %w", err)
	}
	defer file.Close()

	var metainfoFile MetainfoFile
	if err := bencode.Unmarshal(file, &metainfoFile); err != nil {
		return nil, fmt.Errorf("torrentdata: failed to parse torrent file: %w", err)
	}
	return &metainfoFile, nil
}

// computeHash computes the SHA-1 hash of the InfoDictionary.
func (infoDict *InfoDictionary) computeHash() ([20]byte, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, *infoDict); err != nil {
		return [20]byte{}, fmt.Errorf("torrentdata: failed to marshal InfoDictionary: %w", err)
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

// decodePieces decodes the pieces string into individual SHA-1 hashes.
func (i *InfoDictionary) decodePieces() ([][20]byte, error) {
	const sha1HashLen = 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%sha1HashLen != 0 {
		return nil, errors.New("torrentdata: malformed pieces")
	}
	numHashes := len(buf) / sha1HashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*sha1HashLen:(i+1)*sha1HashLen])
	}
	return hashes, nil
}

// toTorrentData converts MetainfoFile to TorrentData.
func (metainfoFile *MetainfoFile) toTorrentData() (TorrentData, error) {
	infoHash, err := metainfoFile.Info.computeHash()
	if err != nil {
		return TorrentData{}, fmt.Errorf("torrentdata: failed to compute InfoHash: %w", err)
	}
	pieceHashes, err := metainfoFile.Info.decodePieces()
	if err != nil {
		return TorrentData{}, fmt.Errorf("torrentdata: failed to split piece hashes: %w", err)
	}
	t := TorrentData{
		Announce:    metainfoFile.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: metainfoFile.Info.PieceLength,
		Length:      metainfoFile.Info.Length,
		Name:        metainfoFile.Info.Name,
	}
	return t, nil
}

func (t *TorrentData) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	baseURL, err := url.Parse(t.Announce)
	if err != nil {
		return "", fmt.Errorf("torrentdata: failed to parse announce URL: %v", err)
	}

	queryParams := url.Values{}
	queryParams.Set("info_hash", string(t.InfoHash[:]))
	peerIDString := string(peerID[:])
	queryParams.Set("peer_id", peerIDString)
	portString := strconv.Itoa(int(port))
	queryParams.Set("port", portString)
	queryParams.Set("uploaded", "0")
	queryParams.Set("downloaded", "0")
	queryParams.Set("compact", "1")
	queryParams.Set("left", strconv.Itoa(t.Length))

	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String(), nil
}
