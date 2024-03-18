package torrentdata

import (
	"errors"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	t.Run("ValidFile", func(t *testing.T) {
		metainfoFile, err := Open("testdata/archlinux-2019.12.01-x86_64.iso.torrent")
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
		} else if err.Error() != "path cannot be empty" {
			t.Errorf("Unexpected error message: got %q, want %q", err.Error(), "path cannot be empty")
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
