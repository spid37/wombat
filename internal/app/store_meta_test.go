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

func TestNilStorePublicAPIReturnsError(t *testing.T) {
	a := NewAPI(t.TempDir())

	if _, err := a.GetWorkspaceOptions(); err != errStoreUnavailable {
		t.Fatalf("GetWorkspaceOptions: got %v, want %v", err, errStoreUnavailable)
	}
	if _, err := a.GetReflectMetadata("x"); err != errStoreUnavailable {
		t.Fatalf("GetReflectMetadata: got %v, want %v", err, errStoreUnavailable)
	}
	if _, err := a.GetMetadata("x"); err != errStoreUnavailable {
		t.Fatalf("GetMetadata: got %v, want %v", err, errStoreUnavailable)
	}
	if _, err := a.ListWorkspaces(); err != errStoreUnavailable {
		t.Fatalf("ListWorkspaces: got %v, want %v", err, errStoreUnavailable)
	}
	if err := a.SelectWorkspace("wksp_x"); err != errStoreUnavailable {
		t.Fatalf("SelectWorkspace: got %v, want %v", err, errStoreUnavailable)
	}
	if err := a.DeleteWorkspace("wksp_x"); err != errStoreUnavailable {
		t.Fatalf("DeleteWorkspace: got %v, want %v", err, errStoreUnavailable)
	}
	if _, err := a.GetRawMessageState("/a/b"); err != errStoreUnavailable {
		t.Fatalf("GetRawMessageState: got %v, want %v", err, errStoreUnavailable)
	}
	if err := a.Connect(map[string]interface{}{"addr": "localhost:1"}, nil, true); err != errStoreUnavailable {
		t.Fatalf("Connect(save): got %v, want %v", err, errStoreUnavailable)
	}
	if cmd := a.ExportCommands("/a/b", "{}", nil); cmd != nil {
		t.Fatalf("ExportCommands: got %+v, want nil", cmd)
	}
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
