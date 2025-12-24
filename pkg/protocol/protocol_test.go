/**************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@license	Copyright © 2021-2026 Michael Roberts

/**************************************************************************************/

package protocol

/**************************************************************************************/

import (
	"bytes"
	"testing"
)

/**************************************************************************************/

func TestFrameIsRequest(t *testing.T) {
	t.Run("IsRequest", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: FlagIsRequest,
			},
		}

		want := true

		if frame.IsRequest() != want {
			t.Fatalf("IsRequest = false, want true")
		}
	})

	t.Run("IsNotRequest", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: 0,
			},
		}

		want := false

		if frame.IsRequest() != want {
			t.Fatalf("IsRequest = true, want false")
		}

	})
}

/**************************************************************************************/

func TestSetRequestFlag(t *testing.T) {
	t.Run("SetRequest", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: 0,
			},
		}

		frame.SetRequest()

		want := true

		if frame.IsRequest() != want {
			t.Fatalf("IsRequest = false, want true")
		}
	})

	t.Run("SetRequestIdempotent", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: 0,
			},
		}

		frame.SetRequest()

		want := true

		if frame.IsRequest() != want {
			t.Fatalf("IsRequest = false, want true")
		}

		frame.SetRequest()

		if frame.IsRequest() != want {
			t.Fatalf("IsRequest = false, want true")
		}
	})
}

/**************************************************************************************/

func TestFrameIsResponse(t *testing.T) {
	t.Run("IsResponse", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: 0,
			},
		}

		want := true

		if frame.IsResponse() != want {
			t.Fatalf("IsResponse = false, want true")
		}
	})

	t.Run("IsNotResponse", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: FlagIsRequest,
			},
		}

		want := false

		if frame.IsResponse() != want {
			t.Fatalf("IsResponse = true, want false")
		}
	})
}

/**************************************************************************************/

func TestSetResponseFlag(t *testing.T) {
	t.Run("SetResponse", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: FlagIsRequest,
			},
		}

		frame.SetResponse()

		want := true

		if frame.IsResponse() != want {
			t.Fatalf("IsResponse = false, want true")
		}
	})

	t.Run("SetResponseIdempotent", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: FlagIsRequest,
			},
		}

		frame.SetResponse()

		want := true

		if frame.IsResponse() != want {
			t.Fatalf("IsResponse = false, want true")
		}

		frame.SetResponse()

		if frame.IsResponse() != want {
			t.Fatalf("IsResponse = false, want true")
		}
	})
}

/**************************************************************************************/

func TestFrameIsError(t *testing.T) {
	t.Run("IsError", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: FlagIsError,
			},
		}

		want := true

		if frame.IsError() != want {
			t.Fatalf("IsError = false, want true")
		}
	})

	t.Run("IsNotError", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: 0,
			},
		}

		want := false

		if frame.IsError() != want {
			t.Fatalf("IsError = true, want false")
		}
	})
}

/**************************************************************************************/

func TestSetErrorFlag(t *testing.T) {
	t.Run("SetError", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: 0,
			},
		}

		frame.SetError()

		want := true

		if frame.IsError() != want {
			t.Fatalf("IsError = false, want true")
		}
	})

	t.Run("SetErrorIdempotent", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Flags: 0,
			},
		}

		frame.SetError()

		want := true

		if frame.IsError() != want {
			t.Fatalf("IsError = false, want true")
		}

		frame.SetError()

		if frame.IsError() != want {
			t.Fatalf("IsError = false, want true")
		}
	})
}

/**************************************************************************************/

func TestFrame(t *testing.T) {
	t.Run("RoundTripEmptyPayload", func(t *testing.T) {
		id := uint16(0x1234)
		group := uint8(0x01)
		code := uint8(0x02)

		frame := NewFrame(id, group, code, nil)

		buffer, err := frame.Serialize()
		if err != nil {
			t.Fatalf("Serialize error = %v", err)
		}

		frame, size, err := Parse(buffer)
		if err != nil {
			t.Fatalf("Parse error = %v", err)
		}

		if size != frame.Size {
			t.Fatalf("size = %d, want %d", size, frame.Size)
		}

		if frame.Version != Version {
			t.Fatalf("version = %d, want %d", frame.Version, Version)
		}

		if frame.MessageID != id {
			t.Fatalf("id = 0x%04X, want 0x%04X", frame.MessageID, id)
		}

		if frame.Group != group {
			t.Fatalf("group = 0x%02X, want 0x%02X", frame.Group, group)
		}

		if frame.Code != code {
			t.Fatalf("code = 0x%02X, want 0x%02X", frame.Code, code)
		}

		if !frame.IsRequest() {
			t.Fatalf("IsRequest = false, want true")
		}

		if len(frame.Payload) != 0 {
			t.Fatalf("payload length = %d, want 0", len(frame.Payload))
		}
	})

	t.Run("RoundTripWithPayload", func(t *testing.T) {
		id := uint16(0xABCD)

		group := uint8(0x10)

		code := uint8(0x20)

		payload := []byte{0xDE, 0xAD, 0xBE, 0xEF}

		frame := NewFrame(id, group, code, payload)

		buffer, err := frame.Serialize()
		if err != nil {
			t.Fatalf("Serialize error = %v", err)
		}

		frame, size, err := Parse(buffer)
		if err != nil {
			t.Fatalf("Parse error = %v", err)
		}

		if size != frame.Size {
			t.Fatalf("size = %d, want %d", size, frame.Size)
		}

		if frame.MessageID != id {
			t.Fatalf("id = 0x%04X, want 0x%04X", frame.MessageID, id)
		}

		if frame.Group != group {
			t.Fatalf("group = 0x%02X, want 0x%02X", frame.Group, group)
		}

		if frame.Code != code {
			t.Fatalf("code = 0x%02X, want 0x%02X", frame.Code, code)
		}

		if !bytes.Equal(frame.Payload, payload) {
			t.Fatalf("payload = % X, want % X", frame.Payload, payload)
		}
	})

	t.Run("SerializeVersionZero", func(t *testing.T) {
		frame := Frame{
			Header: Header{
				Version:   0,
				MessageID: 1,
			},
			Command: Command{
				Group: 0x01,
				Code:  0x01,
			},
		}

		_, err := frame.Serialize()

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("ParseTooShort", func(t *testing.T) {
		buf := make([]byte, MinimumFrameSize-1)

		buf[0] = SyncByte

		_, _, err := Parse(buf)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("ParseTooShort", func(t *testing.T) {
		buf := make([]byte, MinimumFrameSize-1)

		buf[0] = SyncByte

		_, _, err := Parse(buf)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("ParseBadSync", func(t *testing.T) {
		frame := NewFrame(1, 0x01, 0x01, []byte{0x01})

		data, err := frame.Serialize()
		if err != nil {
			t.Fatalf("Serialize error = %v", err)
		}

		data[0] = 0x00

		_, _, err = Parse(data)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("ParseBadChecksum", func(t *testing.T) {
		frame := NewFrame(1, 0x01, 0x01, []byte{0x01, 0x02, 0x03})

		data, err := frame.Serialize()
		if err != nil {
			t.Fatalf("Serialize error = %v", err)
		}

		data[9] ^= 0xFF

		_, _, err = Parse(data)
		if err == nil {
			t.Fatalf("expected checksum error, got nil")
		}
	})

	t.Run("ParseBadVersion", func(t *testing.T) {
		frame := NewFrame(1, 0x01, 0x01, []byte{0x01})

		data, err := frame.Serialize()
		if err != nil {
			t.Fatalf("Serialize error = %v", err)
		}

		data[1] = Version + 1

		_, _, err = Parse(data)
		if err == nil {
			t.Fatalf("expected version error, got nil")
		}
	})
}

/**************************************************************************************/
