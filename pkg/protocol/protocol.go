/**************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@license	Copyright © 2021-2026 Michael Roberts

/**************************************************************************************/

package protocol

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
