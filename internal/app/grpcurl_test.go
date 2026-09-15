package app

import (
	"os"
	"os/exec"
	"testing"
)

func TestParseGrpcurlValidCommand(t *testing.T) {
	got, err := parseGrpcurlCommand(`grpcurl -d '{"x":1}' -rpc-header 'a:b' localhost:50051 pkg.Svc/Method`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Target != "localhost:50051" || got.Method != "pkg.Svc/Method" {
		t.Fatalf("unexpected parse result: %+v", got)
	}
	if got.Data != `{"x":1}` {
		t.Fatalf("unexpected data: %q", got.Data)
	}
	if len(got.Metadata) != 1 || got.Metadata[0].Key != "a" || got.Metadata[0].Val != "b" {
		t.Fatalf("unexpected metadata: %+v", got.Metadata)
	}
}

// Unknown flags must return an error to the caller — not os.Exit the process.
func TestParseGrpcurlUnknownFlagReturnsError(t *testing.T) {
	if os.Getenv("WOMBAT_GRPCURL_SUBPROC") == "1" {
		_, err := parseGrpcurlCommand(`grpcurl -not-a-real-flag localhost:1 Svc/Method`)
		if err != nil {
			os.Exit(0) // returned error instead of exiting — desired behavior
		}
		os.Exit(3) // no error when one was expected
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestParseGrpcurlUnknownFlagReturnsError", "-test.v")
	cmd.Env = append(os.Environ(), "WOMBAT_GRPCURL_SUBPROC=1")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("subprocess exited with error (likely flag.ExitOnError called os.Exit): %v", err)
	}
}
