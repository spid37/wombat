package app

import (
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func TestRetryConnectionNilContextDoesNotBlock(t *testing.T) {
	// Closed conn is Shutdown; without a Wails ctx, once() cannot register —
	// must return immediately instead of waiting for the 2s timeout.
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	srv := grpc.NewServer()
	go srv.Serve(lis)
	t.Cleanup(func() {
		srv.Stop()
		lis.Close()
	})

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	conn.Connect()
	conn.Close()

	a := NewAPI(t.TempDir())
	a.client = &client{conn: conn}

	done := make(chan struct{})
	go func() {
		a.RetryConnection()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("RetryConnection blocked with nil ctx")
	}
}
