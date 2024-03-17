package torrentfile

import (
	"errors"
	"os"

	"github.com/jackpal/bencode-go"
)

type TorrentInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type TorrentFile struct {
	Announce string      `bencode:"announce"`
	Info     TorrentInfo `bencode:"info"`
}

func Open(path string) (*TorrentFile, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	torrentFile := TorrentFile{}
	err = bencode.Unmarshal(file, &torrentFile)
	if err != nil {
		return nil, err
	}
	return &torrentFile, nil
}
