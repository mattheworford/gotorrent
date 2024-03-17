package torrentfile

import (
	"errors"
	"io"
	"os"
	"testing"
)

type MockFile struct {
	data []byte
}

// Read reads data from the mock file.
func (f *MockFile) Read(p []byte) (n int, err error) {
	return copy(p, f.data), io.EOF
}

var bencodeUnmarshal = func(r io.Reader, v interface{}) error {
	return nil // Default behavior: return no error
}

func TestOpen(t *testing.T) {
	t.Run("ValidFile", func(t *testing.T) {
		torrentFile, err := Open("testdata/archlinux-2019.12.01-x86_64.iso.torrent")
		if err != nil {
			t.Fatalf("Open failed: %v", err)
		}

		expectedAnnounce := "http://tracker.archlinux.org:6969/announce"
		if torrentFile.Announce != expectedAnnounce {
			t.Errorf("Unexpected Announce value: got %q, want %q", torrentFile.Announce, expectedAnnounce)
		}

		expectedComment := "archlinux-2019.12.01-x86_64.iso"
		if torrentFile.Info.Name != expectedComment {
			t.Errorf("Unexpected Comment value: got %q, want %q", torrentFile.Info.Name, expectedComment)
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
