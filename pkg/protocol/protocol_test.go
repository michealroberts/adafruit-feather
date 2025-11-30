/**************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@license	Copyright © 2021-2026 Michael Roberts

/**************************************************************************************/

package protocol

/**************************************************************************************/

import (
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
