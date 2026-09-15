package app

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestCloseSendIdempotent(t *testing.T) {
	a := NewAPI(t.TempDir())
	a.streamReq = make(chan proto.Message)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CloseSend panicked on second call: %v", r)
		}
	}()

	a.CloseSend()
	a.CloseSend()
}

func TestCloseSendAfterInternalCloseDoesNotPanic(t *testing.T) {
	a := NewAPI(t.TempDir())
	a.streamReq = make(chan proto.Message)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CloseSend panicked after internal close: %v", r)
		}
	}()

	// Simulate send-error path and UI CloseSend racing through the same helper.
	a.closeStreamReq()
	a.CloseSend()
}
