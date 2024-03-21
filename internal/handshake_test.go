package internal

import (
	"bytes"
	"reflect"
	"testing"
)

func TestHandshake_Serialize(t *testing.T) {
	tests := []struct {
		name      string
		handshake Handshake
		want      []byte
	}{
		{
			name: "SerializeHandshake",
			handshake: Handshake{
				ProtocolString: "BitTorrent protocol",
				InfoHash:       [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				PeerID:         [20]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
			},
			want: []byte{
				19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0,
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handshake.Serialize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handshake.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRead(t *testing.T) {
	mockReader := bytes.NewReader([]byte{
		19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
	})

	wantHandshake := Handshake{
		ProtocolString: "BitTorrent protocol",
		InfoHash:       [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		PeerID:         [20]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
	}

	gotHandshake, err := ParseHandshake(mockReader)
	if err != nil {
		t.Fatalf("Read() error = %v, want nil", err)
	}

	if !reflect.DeepEqual(gotHandshake, &wantHandshake) {
		t.Errorf("Read() got = %v, want %v", gotHandshake, &wantHandshake)
	}
}

func TestHandshake_SerializeAndRead(t *testing.T) {
	testCases := []struct {
		name      string
		handshake *Handshake
	}{
		{
			name: "SerializeAndRead",
			handshake: &Handshake{
				ProtocolString: "BitTorrent protocol",
				InfoHash:       [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
				PeerID:         [20]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serializedHandshake := tc.handshake.Serialize()

			reader := bytes.NewReader(serializedHandshake)

			readHandshake, err := ParseHandshake(reader)
			if err != nil {
				t.Fatalf("Read() error = %v, want nil", err)
			}

			if !reflect.DeepEqual(readHandshake, tc.handshake) {
				t.Errorf("Read() got = %v, want %v", readHandshake, tc.handshake)
			}
		})
	}
}
