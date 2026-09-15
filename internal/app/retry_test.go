package app

import (
	"testing"
)

func TestRetryConnectionNilClientDoesNotPanic(t *testing.T) {
	a := NewAPI(t.TempDir())
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RetryConnection panicked with nil client: %v", r)
		}
	}()
	a.RetryConnection()
}

func TestRetryConnectionNilConnDoesNotPanic(t *testing.T) {
	a := NewAPI(t.TempDir())
	a.client = &client{}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RetryConnection panicked with nil conn: %v", r)
		}
	}()
	a.RetryConnection()
}
