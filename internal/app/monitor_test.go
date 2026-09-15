package app

import (
	"context"
	"testing"
	"time"
)

// Verifies #110: when client is nil, monitor must respect context cancel
// instead of spinning a busy loop.
func TestMonitorStateChangesExitsOnCancelWhenClientNil(t *testing.T) {
	a := NewAPI(t.TempDir())
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		a.monitorStateChanges(ctx)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("monitorStateChanges did not exit after cancel with nil client (busy loop)")
	}
}
