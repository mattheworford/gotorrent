package torrentfile

import (
	"testing"
)

func TestOpen(t *testing.T) {
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
}
