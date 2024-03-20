package torrentdata

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	t.Run("ValidFile", func(t *testing.T) {
		metainfoFile, err := Open("../../test/data/archlinux-2019.12.01-x86_64.iso.torrent")
		if err != nil {
			t.Fatalf("Open failed: %v", err)
		}

		expectedAnnounce := "http://tracker.archlinux.org:6969/announce"
		if metainfoFile.Announce != expectedAnnounce {
			t.Errorf("Unexpected Announce value: got %q, want %q", metainfoFile.Announce, expectedAnnounce)
		}

		expectedInfoName := "archlinux-2019.12.01-x86_64.iso"
		if metainfoFile.Info.Name != expectedInfoName {
			t.Errorf("Unexpected Comment value: got %q, want %q", metainfoFile.Info.Name, expectedInfoName)
		}
	})

	t.Run("EmptyPath", func(t *testing.T) {
		_, err := Open("")
		if err == nil {
			t.Error("Expected error, got nil")
		} else if err.Error() != "torrentdata: path cannot be empty" {
			t.Errorf("Unexpected error message: got %q, want %q", err.Error(), "torrentdata: path cannot be empty")
		}
	})

	t.Run("InvalidFile", func(t *testing.T) {
		_, err := Open("testdata/nonexistent.torrent")
		if err == nil {
			t.Error("Expected error, got nil")
		} else if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("Unexpected error: got %v, want %v", err, os.ErrNotExist)
		}
	})
}

func TestToTorrentData(t *testing.T) {
	t.Run("ValidMetainfoFile", func(t *testing.T) {
		metainfoFile := MetainfoFile{
			Announce: "http://tracker.example.com",
			Info: InfoDictionary{
				Pieces:      "1234567890abcdefghij",
				PieceLength: 256,
				Length:      1024,
				Name:        "example.torrent",
			},
		}

		expectedTorrentData := TorrentData{
			Announce:    "http://tracker.example.com",
			InfoHash:    [20]byte{},
			PieceHashes: [][20]byte{},
			PieceLength: 256,
			Length:      1024,
			Name:        "example.torrent",
		}

		actualTorrentData, err := metainfoFile.toTorrentData()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		mockInfoHash := [20]byte{26, 96, 139, 207, 103, 107, 192, 195, 176, 60, 164, 43, 162, 89, 18, 65, 96, 50, 130, 221}
		mockPieceHashes := [][20]byte{
			{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106},
		}

		if !reflect.DeepEqual(actualTorrentData.Announce, expectedTorrentData.Announce) {
			t.Errorf("Unexpected Announce value. Expected: %s, Got: %s", expectedTorrentData.Announce, actualTorrentData.Announce)
		}
		if !bytes.Equal(actualTorrentData.InfoHash[:], mockInfoHash[:]) {
			t.Errorf("Unexpected InfoHash value. Expected: %v, Got: %v", mockInfoHash, actualTorrentData.InfoHash)
		}
		if !reflect.DeepEqual(actualTorrentData.PieceHashes, mockPieceHashes) {
			t.Errorf("Unexpected PieceHashes value. Expected: %v, Got: %v", mockPieceHashes, actualTorrentData.PieceHashes)
		}
		if actualTorrentData.PieceLength != expectedTorrentData.PieceLength {
			t.Errorf("Unexpected PieceLength value. Expected: %d, Got: %d", expectedTorrentData.PieceLength, actualTorrentData.PieceLength)
		}
		if actualTorrentData.Length != expectedTorrentData.Length {
			t.Errorf("Unexpected Length value. Expected: %d, Got: %d", expectedTorrentData.Length, actualTorrentData.Length)
		}
		if actualTorrentData.Name != expectedTorrentData.Name {
			t.Errorf("Unexpected Name value. Expected: %s, Got: %s", expectedTorrentData.Name, actualTorrentData.Name)
		}
	})
	t.Run("MalformedPieces", func(t *testing.T) {
		metainfoFile := MetainfoFile{
			Announce: "http://tracker.example.com",
			Info: InfoDictionary{
				Pieces:      "1",
				PieceLength: 256,
				Length:      1024,
				Name:        "example.torrent",
			},
		}

		_, err := metainfoFile.toTorrentData()
		if err == nil {
			t.Error("Expected error, got nil")
		} else if err.Error() != "torrentdata: failed to split piece hashes: torrentdata: malformed pieces" {
			t.Errorf("Unexpected error message: got %q, want %q", err.Error(), "torrentdata: failed to split piece hashes: torrentdata: malformed pieces")
		}
	})
}
func TestBuildTrackerURL(t *testing.T) {
	t.Run("ValidAnnounceUrl", func(t *testing.T) {
		torrentData := TorrentData{
			Announce:    "http://tracker.example.com",
			InfoHash:    [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			PieceHashes: [][20]byte{},
			PieceLength: 256,
			Length:      1024,
			Name:        "example.torrent",
		}

		peerID := [20]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}
		port := uint16(8080)

		expectedURL := "http://tracker.example.com?compact=1&downloaded=0&info_hash=%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14&left=1024&peer_id=%15%16%17%18%19%1A%1B%1C%1D%1E%1F+%21%22%23%24%25%26%27%28&port=8080&uploaded=0"

		url, err := torrentData.buildTrackerURL(peerID, port)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if url != expectedURL {
			t.Errorf("Unexpected URL. Expected: %s, Got: %s", expectedURL, url)
		}
	})

	t.Run("MalformedAnnounceUrl", func(t *testing.T) {
		torrentData := TorrentData{
			Announce:    ":)",
			InfoHash:    [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			PieceHashes: [][20]byte{},
			PieceLength: 256,
			Length:      1024,
			Name:        "example.torrent",
		}

		peerID := [20]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}
		port := uint16(8080)

		_, err := torrentData.buildTrackerURL(peerID, port)
		if err == nil {
			t.Error("Expected error, got nil")
		} else if err.Error() != "failed to parse announce URL: parse \":)\": missing protocol scheme" {
			t.Errorf("Unexpected error message: got %q, want %q", err.Error(), "failed to parse announce URL: parse \":)\": missing protocol scheme")
		}
	})
}
