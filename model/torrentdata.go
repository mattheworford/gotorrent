package torrentdata

import (
	"errors"
	"os"

	"github.com/jackpal/bencode-go"
)

type InfoDictionary struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type MetainfoFile struct {
	Announce string         `bencode:"announce"`
	Info     InfoDictionary `bencode:"info"`
}

type TorrentData struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func Open(path string) (*MetainfoFile, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	metainfoFile := MetainfoFile{}
	err = bencode.Unmarshal(file, &metainfoFile)
	if err != nil {
		return nil, err
	}
	return &metainfoFile, nil
}
