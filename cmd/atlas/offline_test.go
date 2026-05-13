package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/issuance"
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

type testClock struct{ t time.Time }

func (c testClock) Now() time.Time { return c.t }

type allPerms struct{ set issuance.PermissionSet }

func (a allPerms) PermissionsOf(spiffeid.ID) (issuance.PermissionSet, bool) { return a.set, true }

// buildBundle creates a real key, issues a real record, publishes a real
// signed snapshot, and writes a bundle file — the full offline fixture.
func buildBundle(t *testing.T, revokeIt bool, asOf time.Time) (bundlePath, rec string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	clock := testClock{t: asOf}
	authority, err := issuance.NewAuthority(
		record.Signer{Key: key, KeyID: "authority-key-1"},
		allPerms{set: issuance.NewPermissionSet("read:orders", "write:audit")},
		issuance.NoBinding{}, issuance.RandomMinter{}, clock)
	if err != nil {
		t.Fatal(err)
	}
	res, err := authority.Issue(issuance.Request{
		Principal: spiffeid.RequireFromString("spiffe://domain-a.test/workload/api"),
		Delegate:  spiffeid.RequireFromString("spiffe://domain-b.test/agent/worker"),
		Scope:     []string{"read:orders"}, Expiration: asOf.Add(time.Hour),
	})
	if err != nil || res.Outcome != issuance.Issued {
		t.Fatalf("issue: %v %v", err, res.Outcome)
	}
	var revoked []record.InstanceID
	if revokeIt {
		revoked = []record.InstanceID{res.Record.Read().Instance}
	}
	pub, err := revstatus.NewPublisher(key, "test-list")
	if err != nil {
		t.Fatal(err)
	}
	snap, err := pub.Publish(revoked, asOf)
	if err != nil {
		t.Fatal(err)
	}
	der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	revStrs := make([]string, 0, len(snap.Revoked))
	for _, r := range snap.Revoked {
		revStrs = append(revStrs, r.String())
	}
	b := map[string]any{
		"version": 1, "trustDomain": "domain-a.test",
		"keys": map[string]string{"authority-key-1": base64.StdEncoding.EncodeToString(der)},
		"revocation": map[string]any{
			"listId": snap.ListID, "asOf": snap.AsOf, "revoked": revStrs, "sig": snap.Sig,
		},
		"exportedAt": asOf,
	}
	raw, _ := json.Marshal(b)
	path := filepath.Join(t.TempDir(), "bundle.json")
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	return path, string(res.Record.Presented())
}

func TestOfflineVerifyAccept(t *testing.T) {
	path, rec := buildBundle(t, false, time.Now())
	var out bytes.Buffer
	if code := verifyOffline(&out, path, rec, 10*time.Minute, ""); code != exitOK {
		t.Fatalf("exit=%d out=%s", code, out.String())
	}
	if !strings.Contains(out.String(), "ACCEPT") || !strings.Contains(out.String(), "offline") {
		t.Fatalf("output: %s", out.String())
	}
}

func TestOfflineVerifyRevokedRejects(t *testing.T) {
	path, rec := buildBundle(t, true, time.Now())
	var out bytes.Buffer
	if code := verifyOffline(&out, path, rec, 10*time.Minute, ""); code != exitNonPass {
		t.Fatalf("revoked should exit %d, got %d: %s", exitNonPass, code, out.String())
	}
	if !strings.Contains(out.String(), "RevokedObservable") {
		t.Fatalf("output: %s", out.String())
	}
}

func TestOfflineStaleSnapshotFailsClosed(t *testing.T) {
	// snapshot signed 10s in the past; budget 1s → fail closed (inconclusive).
	path, rec := buildBundle(t, false, time.Now().Add(-10*time.Second))
	var out bytes.Buffer
	if code := verifyOffline(&out, path, rec, time.Second, ""); code != exitNonPass {
		t.Fatalf("stale should exit %d, got %d: %s", exitNonPass, code, out.String())
	}
	if !strings.Contains(out.String(), "RevocationKnowledgeStale") {
		t.Fatalf("expected fail-closed staleness, got: %s", out.String())
	}
}

func TestOfflineScopeGate(t *testing.T) {
	// the record grants read:orders (see buildBundle).
	path, rec := buildBundle(t, false, time.Now())

	// authorized for a granted scope
	var out bytes.Buffer
	if code := verifyOffline(&out, path, rec, 10*time.Minute, "read:orders"); code != exitOK {
		t.Fatalf("granted scope should authorize (exit 0), got %d: %s", code, out.String())
	}
	if !strings.Contains(out.String(), "authorized: YES") {
		t.Fatalf("output: %s", out.String())
	}

	// NOT authorized for a scope the capability doesn't grant — even though the
	// record itself is perfectly valid.
	out.Reset()
	if code := verifyOffline(&out, path, rec, 10*time.Minute, "vercel:deploy:prod"); code != exitNonPass {
		t.Fatalf("ungranted scope should be denied (exit %d), got %d: %s", exitNonPass, code, out.String())
	}
	if !strings.Contains(out.String(), "authorized: NO") || !strings.Contains(out.String(), "does not grant") {
		t.Fatalf("output: %s", out.String())
	}
}

func TestOfflineScopeGateDeniedWhenRevoked(t *testing.T) {
	// even a granted scope must be denied if the capability is revoked.
	path, rec := buildBundle(t, true, time.Now())
	var out bytes.Buffer
	if code := verifyOffline(&out, path, rec, 10*time.Minute, "read:orders"); code != exitNonPass {
		t.Fatalf("revoked capability must not authorize, got %d: %s", code, out.String())
	}
	if !strings.Contains(out.String(), "authorized: NO") {
		t.Fatalf("output: %s", out.String())
	}
}

func TestOfflineTamperedBundleRefused(t *testing.T) {
	path, rec := buildBundle(t, true, time.Now())
	// strip the revocation (the attack) — signature must no longer verify.
	raw, _ := os.ReadFile(path)
	var b map[string]any
	json.Unmarshal(raw, &b)
	b["revocation"].(map[string]any)["revoked"] = []string{}
	tampered, _ := json.Marshal(b)
	os.WriteFile(path, tampered, 0o600)

	var out bytes.Buffer
	if code := verifyOffline(&out, path, rec, 10*time.Minute, ""); code != exitError {
		t.Fatalf("tampered bundle should exit %d, got %d: %s", exitError, code, out.String())
	}
	if !strings.Contains(out.String(), "REFUSING bundle") {
		t.Fatalf("output: %s", out.String())
	}
}
