package message

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestReadPeerMessage(t *testing.T) {
	testCases := []struct {
		name          string
		reader        io.Reader
		expected      *PeerMessage
		expectErr     bool
		expectedError error
	}{
		{
			name:   "ValidMessage",
			reader: bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x05, 0x01, 0x02, 0x03, 0x04, 0x05}),
			expected: &PeerMessage{
				Type:    PeerMessageType(0x01),
				Payload: []byte{0x02, 0x03, 0x04, 0x05},
			},
			expectErr:     false,
			expectedError: nil,
		},
		{
			name:          "EmptyMessage",
			reader:        bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00}),
			expected:      nil,
			expectErr:     false,
			expectedError: nil,
		},
		{
			name:          "InvalidMessageLength",
			reader:        bytes.NewReader(bytes.Repeat([]byte{0x0A}, 4)),
			expected:      nil,
			expectErr:     true,
			expectedError: errors.New("message length 168430090 exceeds maximum allowed"),
		},
		{
			name:          "ReadError",
			reader:        &errorReader{},
			expected:      nil,
			expectErr:     true,
			expectedError: errors.New("failed to read message length: read error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			message, err := ReadPeerMessage(tc.reader)

			if tc.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("Unexpected error: got %q, want %q", err.Error(), tc.expectedError.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(message, tc.expected) {
					t.Errorf("Unexpected result. Expected: %v, Got: %v", tc.expected, message)
				}
			}
		})
	}
}

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func TestSerializePeerMessage(t *testing.T) {
	testCases := []struct {
		name          string
		message       *PeerMessage
		expected      []byte
		expectErr     bool
		expectedError error
	}{
		{
			name:          "NilMessage",
			message:       nil,
			expected:      []byte{0x00, 0x00, 0x00, 0x00},
			expectErr:     false,
			expectedError: nil,
		},
		{
			name: "ValidMessage",
			message: &PeerMessage{
				Type:    PeerMessageType(0x01),
				Payload: []byte{0x02, 0x03, 0x04, 0x05},
			},
			expected:      []byte{0x00, 0x00, 0x00, 0x05, 0x01, 0x02, 0x03, 0x04, 0x05},
			expectErr:     false,
			expectedError: nil,
		},
		{
			name: "InvalidPayloadLength",
			message: &PeerMessage{
				Type:    PeerMessageType(0x01),
				Payload: make([]byte, MaxPayloadLength+1),
			},
			expected:      nil,
			expectErr:     true,
			expectedError: errors.New("payload length 16385 exceeds maximum allowed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			serialized, err := tc.message.Serialize()

			if tc.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("Unexpected error: got %q, want %q", err.Error(), tc.expectedError.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !bytes.Equal(serialized, tc.expected) {
					t.Errorf("Unexpected result. Expected: %v, Got: %v", tc.expected, serialized)
				}
			}
		})
	}
}
