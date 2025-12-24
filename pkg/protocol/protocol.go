/**************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@license	Copyright © 2021-2026 Michael Roberts

/**************************************************************************************/

package protocol

/**************************************************************************************/

import (
	"errors"
	"hash/crc32"

	"github.com/michealroberts/adafruit-feather/internal/crc"
)

/**************************************************************************************/

// SyncByte is the fixed marker that identifies the start of a frame on the wire:
const SyncByte uint8 = 0xAA

/**************************************************************************************/

// Version is the protocol version number:
const Version uint8 = 0x01

/**************************************************************************************/

// Flags used in the protocol frame to indicate specific conditions:
const (
	FlagIsRequest uint8 = 1 << 0
	FlagIsError   uint8 = 1 << 1
	// Bits 2..7 reserved for future use.
)

/**************************************************************************************/

// Size constants define the byte lengths of protocol frame components:
const (
	SyncSize         int = 1
	HeaderSize       int = 8 // version + flags + messageId(2) + size(2) + group + code
	ChecksumSize     int = 4
	MinimumFrameSize int = SyncSize + HeaderSize + ChecksumSize
)

/**************************************************************************************/

// Header represents the metadata for a protocol frame:
type Header struct {
	// Version is the protocol version number:
	Version uint8
	// Flags indicate specific conditions for the frame:
	Flags uint8
	// MessageId is a unique identifier for the message within the protocol:
	MessageID uint16
	// Size is the total size of the frame including header, command, and payload:
	Size uint16
}

/**************************************************************************************/

// Command represents the command details within a protocol frame:
type Command struct {
	// Group represents the command group identifier, (e.g., system, sensor, etc):
	Group uint8
	// Code represents the specific command code within the group:
	Code uint8
}

/**************************************************************************************/

// Frame represents a complete protocol frame including header, command, and payload:
type Frame struct {
	// Header contains the metadata for the frame:
	Header
	// Command represents the command details within the frame:
	Command
	// Payload contains the data associated with the frame (optional; variable length):
	Payload []byte
}

/**************************************************************************************/

// NewFrame creates a new Frame with the given parameters. By default it sets the frame
// as a request and calculates the total packet size.
func NewFrame(id uint16, group, code uint8, payload []byte) *Frame {
	return &Frame{
		Header: Header{
			Version:   Version,
			MessageID: id,
			Flags:     FlagIsRequest,
			Size:      uint16(len(payload)) + uint16(HeaderSize),
		},
		Command: Command{
			Group: group,
			Code:  code,
		},
		Payload: payload,
	}
}

/**************************************************************************************/

// IsRequest returns true if the frame's FlagIsRequest bit is set.
func (f Frame) IsRequest() bool {
	return (f.Flags & FlagIsRequest) != 0
}

/**************************************************************************************/

// SetRequest sets the FlagIsRequest bit in the frame's Flags field.
func (f *Frame) SetRequest() {
	f.Flags |= FlagIsRequest
}

/**************************************************************************************/

// IsResponse returns true if the frame's FlagIsRequest bit is not set.
func (f Frame) IsResponse() bool {
	return (f.Flags & FlagIsRequest) == 0
}

/**************************************************************************************/

// SetResponse clears the FlagIsRequest bit in the frame's Flags field.
func (f *Frame) SetResponse() {
	f.Flags &^= FlagIsRequest
}

/**************************************************************************************/

// IsError returns true if the frame's FlagIsError bit is set.
func (f Frame) IsError() bool {
	return (f.Flags & FlagIsError) != 0
}

/**************************************************************************************/

// SetError sets the FlagIsError bit in the frame's Flags field.
func (f *Frame) SetError() {
	f.Flags |= FlagIsError
}

/**************************************************************************************/

// Serialize converts the Frame into a byte slice suitable for transmission over the
// wire according to the protocol specification.
func (f Frame) Serialize() ([]byte, error) {
	if f.Version == 0 {
		return nil, errors.New("frame version is not set")
	}

	size := uint16(len(f.Payload)) + uint16(8)

	if f.Size != 0 && f.Size != size {
		return nil, errors.New("frame size field does not match actual payload length + header size")
	}

	// The buffer contains enough bytes for the total size of the frame, plus one for the
	// sync byte and four for the CRC32 checksum:
	buffer := make([]byte, SyncSize+int(size)+ChecksumSize)

	// Add the sync byte to the buffer:
	buffer[0] = SyncByte

	// Add the version to the header:
	buffer[1] = f.Version

	// Add the flags to the header:
	buffer[2] = f.Flags

	// Add the message ID to the header:
	buffer[3] = byte(f.MessageID >> 8)
	buffer[4] = byte(f.MessageID)

	// Add the size to the header:
	buffer[5] = byte(size >> 8)
	buffer[6] = byte(size)

	// Add the command to the buffer:
	buffer[7] = f.Group
	buffer[8] = f.Code

	// Add the payload to the buffer:
	copy(buffer[9:], f.Payload)

	// Calculate the CRC32 checksum for the frame (excluding the sync byte and checksum
	// itself):
	checksum := crc32.Checksum(buffer[:size+1], crc.Table)

	// Calculate the offset for the checksum based on the sync size and frame size:
	offset := SyncSize + int(size)

	// Add the 4-byte checksum to the buffer (big-endian):
	buffer[offset] = byte(checksum >> 24)
	buffer[offset+1] = byte(checksum >> 16)
	buffer[offset+2] = byte(checksum >> 8)
	buffer[offset+3] = byte(checksum)

	return buffer, nil
}

/**************************************************************************************/

// Parse takes a byte slice containing a serialized frame and reconstructs the Frame
// structure. It returns an error if the frame is invalid or cannot be parsed.
func Parse(buffer []byte) (*Frame, uint16, error) {
	size := uint16(0)

	// Check if the buffer length is less than the minimum frame size:
	if len(buffer) < MinimumFrameSize {
		return nil, 0, errors.New("frame is too short")
	}

	// Ensure that the first byte in the buffer is the sync byte:
	if buffer[0] != SyncByte {
		return nil, 0, errors.New("invalid sync byte")
	}

	version := buffer[1]

	// Check if the version is correct:
	if version != Version {
		return nil, 0, errors.New("unsupported protocol version")
	}

	// Extract size (header-without-sync + payload), big-endian:
	size = uint16(buffer[5])<<8 | uint16(buffer[6])

	// Check if the buffer length is less than the declared size:
	if len(buffer) < SyncSize+int(size)+ChecksumSize {
		return nil, 0, errors.New("incomplete frame, buffer shorter than declared size")
	}

	checksum := crc32.Checksum(buffer[:SyncSize+int(size)], crc.Table)

	// Calculate the offset for the checksum based on the sync size and frame size:
	offset := SyncSize + int(size)

	// Verify the CRC32 checksum for the frame (including the sync byte but excluding the
	// checksum field itself):
	crc := uint32(buffer[offset])<<24 | uint32(buffer[offset+1])<<16 | uint32(buffer[offset+2])<<8 | uint32(buffer[offset+3])

	if checksum != crc {
		return nil, 0, errors.New("invalid CRC32 checksum")
	}

	// Extract the message ID from the buffer:
	id := uint16(buffer[3])<<8 | uint16(buffer[4])

	// Extract the flags from the buffer:
	flags := buffer[2]

	// Extract the group from the buffer:
	group := buffer[7]

	// Extract the code from the buffer:
	code := buffer[8]

	payloadSize := int(size) - HeaderSize

	if payloadSize < 0 {
		return nil, 0, errors.New("invalid payload size, frame size is too small")
	}

	// Calculate the start index for the payload in the buffer:
	start := SyncSize + HeaderSize

	// Calculate the end index for the payload in the buffer:
	end := start + payloadSize

	payload := make([]byte, payloadSize)
	copy(payload, buffer[start:end])

	// Re-construct the header from the buffer:
	header := Header{
		Version:   version,
		Flags:     flags,
		MessageID: id,
		Size:      size,
	}

	// Re-construct the command from the buffer:
	command := Command{
		Group: group,
		Code:  code,
	}

	return &Frame{
		Header:  header,
		Command: command,
		Payload: payload,
	}, size, nil
}

/**************************************************************************************/
