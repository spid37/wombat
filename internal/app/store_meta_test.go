package app

import (
	"testing"
)

func testAPIWithStore(t *testing.T) *API {
	t.Helper()
	a := NewAPI(t.TempDir())
	var err error
	a.store, err = newStore(a.appData, newAppLogger(nil, "DB"))
	if err != nil {
		t.Fatalf("newStore: %v", err)
	}
	a.state = &workspaceState{CurrentID: defaultWorkspaceKey}
	t.Cleanup(func() { a.store.close() })
	return a
}

func TestGetReflectMetadataMissingKeyReturnsEmpty(t *testing.T) {
	a := testAPIWithStore(t)
	hds, err := a.GetReflectMetadata("localhost:5001")
	if err != nil {
		t.Fatalf("expected nil error for missing key, got %v", err)
	}
	if hds == nil {
		t.Fatal("expected non-nil empty headers")
	}
	if len(hds) != 0 {
		t.Fatalf("expected empty headers, got %+v", hds)
	}
}

func TestGetMetadataMissingKeyReturnsEmpty(t *testing.T) {
	a := testAPIWithStore(t)
	hds, err := a.GetMetadata("localhost:5001")
	if err != nil {
		t.Fatalf("expected nil error for missing key, got %v", err)
	}
	if len(hds) != 0 {
		t.Fatalf("expected empty headers, got %+v", hds)
	}
}

func TestGetRawMessageStateMissingKeyReturnsEmpty(t *testing.T) {
	a := testAPIWithStore(t)
	val, err := a.GetRawMessageState("/pkg.Svc/Method")
	if err != nil {
		t.Fatalf("expected nil error for missing key, got %v", err)
	}
	if val != "" {
		t.Fatalf("expected empty string, got %q", val)
	}
}

func TestWailsReadyNilStoreDoesNotPanic(t *testing.T) {
	a := NewAPI(t.TempDir())
	// store intentionally left nil (startup failure path)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("wailsReady panicked with nil store: %v", r)
		}
	}()
	a.wailsReady()
}

func TestShutdownNilStoreDoesNotPanic(t *testing.T) {
	a := NewAPI(t.TempDir())
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("shutdown panicked with nil store: %v", r)
		}
	}()
	a.shutdown(nil)
}
