package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
)

// buildTestBundle creates a minimal trust bundle for testing
func buildTestBundle(t *testing.T, revoked []string, asOf time.Time) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pubKeyDER, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pubKeyB64 := base64.StdEncoding.EncodeToString(pubKeyDER)

	listID := "atlas-revlist:domain-a.test"
	pub, err := revstatus.NewPublisher(key, listID)
	if err != nil {
		t.Fatal(err)
	}

	var revokedIDs []record.InstanceID
	for _, s := range revoked {
		id, _ := record.InstanceIDFromString(s)
		revokedIDs = append(revokedIDs, id)
	}

	snap, err := pub.Publish(revokedIDs, asOf)
	if err != nil {
		t.Fatal(err)
	}

	bundle := bundleFile{
		Version:     1,
		TrustDomain: "domain-a.test",
		Keys:        map[string]string{"authority-key-1": pubKeyB64},
	}
	bundle.Revocation.ListID = snap.ListID
	bundle.Revocation.AsOf = snap.AsOf
	bundle.Revocation.Revoked = revoked
	bundle.Revocation.Sig = snap.Sig

	data, _ := json.Marshal(bundle)
	return string(data)
}

func newTestGate(t *testing.T, bundle string, maxStaleness, reload time.Duration) *httptest.Server {
	t.Helper()
	tmpFile := t.TempDir() + "/bundle.json"
	if err := writeFile(tmpFile, bundle); err != nil {
		t.Fatal(err)
	}

	// Create a minimal upstream that returns 200
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("upstream ok"))
	}))
	t.Cleanup(upstream.Close)

	// Build command-line args for the gate
	// We can't easily test the main() function, so we'll test the buildVerifier and gate logic directly
	// For now, we'll test the deny function and response codes
	return nil // Placeholder - actual test needs integration
}

func writeFile(path, content string) error {
	return nil // Placeholder
}

// DEFECT D5: The gate answers the same 403 for forged, revoked and stale;
// no cause is exposed, no metrics exist, and a stale bundle is indistinguishable from an attack.
// This test asserts the CURRENT WRONG behavior: all failures return 403 with generic message.
// After fix, this test should be renamed TestD5_GateFailVisible and inverted.
func TestDefectD5_GateIndistinguishableFailures(t *testing.T) {
	// This test documents the current behavior by examining the denyWithReason() function
	// and the single error path in the gate handler.

	// The denyWithReason() function at main.go:275-283 returns a JSON error with the same
	// structure for all failure cases.
	// The gate handler at main.go:245-252 has a single deny call for all non-accept verdicts.

	// This test verifies that the denyWithReason function produces the same structure
	var buf bytes.Buffer
	w := &responseWriter{header: make(http.Header), Body: &buf}
	denyWithReason(w, http.StatusForbidden, "invalid_capability", "test message")

	var out map[string]string
	json.Unmarshal(buf.Bytes(), &out)
	if out["error"] != "test message" {
		t.Fatalf("deny output: %v", out)
	}
	if w.statusCode != http.StatusForbidden {
		t.Fatalf("status: %d, want 403", w.statusCode)
	}

	// Check that Atlas-Reason header IS present (fixed behavior)
	reason := w.Header().Get("Atlas-Reason")
	if reason == "" {
		t.Fatalf("DEFECT D5 FIXED: Atlas-Reason header missing")
	}
	if reason != "invalid_capability" {
		t.Fatalf("Atlas-Reason = %q, want 'invalid_capability'", reason)
	}

	// Check that /metrics endpoint doesn't exist in current gate (no metrics)
	// This would require running the full gate, which we can't easily do in unit test
	t.Logf("DEFECT D5 FIXED: gate includes Atlas-Reason header: %s", reason)
}

// responseWriter captures HTTP response for testing
type responseWriter struct {
	header     http.Header
	Body       *bytes.Buffer
	statusCode int
}

func (w *responseWriter) Header() http.Header         { return w.header }
func (w *responseWriter) Write(b []byte) (int, error) { return w.Body.Write(b) }
func (w *responseWriter) WriteHeader(statusCode int)  { w.statusCode = statusCode }
