package internal

import (
	"net"
	"reflect"
	"testing"
)

func TestDecodePeers(t *testing.T) {
	testCases := []struct {
		name       string
		peerData   []byte
		expected   []Peer
		expectErr  bool
		errMessage string
	}{
		{
			name:       "ValidPeerData",
			peerData:   []byte{192, 168, 0, 1, 0, 80, 192, 168, 0, 2, 0, 81},
			expected:   []Peer{{IP: net.IPv4(192, 168, 0, 1), Port: 80}, {IP: net.IPv4(192, 168, 0, 2), Port: 81}},
			expectErr:  false,
			errMessage: "",
		},
		{
			name:       "NilPeerData",
			peerData:   nil,
			expected:   nil,
			expectErr:  true,
			errMessage: "input data is nil",
		},
		{
			name:       "MalformedPeerData",
			peerData:   []byte{192, 168, 0, 1, 0, 80, 192, 168, 0, 2},
			expected:   nil,
			expectErr:  true,
			errMessage: "malformed peer data: incorrect size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			peers, err := DecodePeers(tc.peerData)

			if tc.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tc.errMessage {
					t.Errorf("Unexpected error message: got %q, want %q", err.Error(), tc.errMessage)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(peers, tc.expected) {
					t.Errorf("Unexpected result. Expected: %v, Got: %v", tc.expected, peers)
				}
			}
		})
	}
}
