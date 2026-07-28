package vectors_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/Raj-glitch-max/atlas/tests/conformance"
	jose "github.com/go-jose/go-jose/v3"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

const revocationPath = "revocation-vectors.json"

// TestRevocationVectorsRegenerate rewrites the committed revocation vectors.
// Run deliberately with -update; commit the result.
func TestRevocationVectorsRegenerate(t *testing.T) {
	if !*update {
		t.Skip("run with -update to regenerate committed vectors")
	}
	vf, err := conformance.EmitRevocationVectors()
	if err != nil {
		t.Fatalf("emit revocation vectors: %v", err)
	}
	data, err := json.MarshalIndent(vf, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(revocationPath, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %d revocation vectors to %s", len(vf.Vectors), revocationPath)
}

// TestRevocationVectorsReplay is the drift guard and the executable template a
// foreign implementation mirrors.
//
// Every input is reconstructed FROM THE JSON using only public APIs. Crucially,
// the revocation state is NOT read from the vector — it is DERIVED by ingesting
// the published snapshots and querying the provider, which is precisely the
// work the schema-1 vectors skipped.
func TestRevocationVectorsReplay(t *testing.T) {
	raw, err := os.ReadFile(revocationPath)
	if err != nil {
		t.Fatalf("read %s: %v", revocationPath, err)
	}
	var vf conformance.RevocationVectorFile
	if err := json.Unmarshal(raw, &vf); err != nil {
		t.Fatalf("parse %s: %v", revocationPath, err)
	}
	if vf.Schema != conformance.RevocationSchemaVersion {
		t.Fatalf("schema %d, want %d", vf.Schema, conformance.RevocationSchemaVersion)
	}
	if len(vf.Vectors) == 0 {
		t.Fatal("no revocation vectors committed")
	}

	for _, v := range vf.Vectors {
		t.Run(v.Name, func(t *testing.T) {
			now, err := time.Parse(time.RFC3339Nano, v.Now)
			if err != nil {
				t.Fatalf("now: %v", err)
			}

			// --- trust material for the RECORD ---
			keys := map[string]*ecdsa.PublicKey{}
			for _, kraw := range v.Trust.Keys {
				var jwk jose.JSONWebKey
				if err := jwk.UnmarshalJSON(kraw); err != nil {
					t.Fatalf("trust key: %v", err)
				}
				pk, ok := jwk.Key.(*ecdsa.PublicKey)
				if !ok {
					t.Fatalf("trust key %q is not an EC public key", jwk.KeyID)
				}
				keys[jwk.KeyID] = pk
			}
			domain, err := spiffeid.TrustDomainFromString(v.Trust.Domain)
			if err != nil {
				t.Fatalf("trust domain: %v", err)
			}
			tm, err := record.NewTrustMaterial(domain, keys)
			if err != nil {
				t.Fatalf("trust material: %v", err)
			}
			ts := &staticTrust{domain: domain, mat: tm}

			// --- the key the RP holds for the REVOCATION ORIGIN ---
			var revJWK jose.JSONWebKey
			if err := revJWK.UnmarshalJSON(v.SnapshotKey); err != nil {
				t.Fatalf("snapshot key: %v", err)
			}
			revPub, ok := revJWK.Key.(*ecdsa.PublicKey)
			if !ok {
				t.Fatal("snapshot key is not an EC public key")
			}

			// --- ingest the published snapshots, in order ---
			//
			// This is the work the schema-1 vectors skipped entirely: signature
			// verification, list-id binding, and monotonic adoption all happen
			// here, and each snapshot's adoption is asserted individually so a
			// verifier that adopts a forged or rolled-back snapshot is caught at
			// the point of the mistake rather than only via the final verdict.
			provider := revstatus.NewSignedSetProvider(revPub, v.ExpectList)
			for i, snap := range v.Ingest {
				set, err := conformance.SnapshotFromVector(snap)
				if err != nil {
					t.Fatalf("ingest[%d] rebuild: %v", i, err)
				}
				adopted, _ := provider.Ingest(set)
				if i < len(v.ExpectAdopt) && adopted != v.ExpectAdopt[i] {
					t.Fatalf("ingest[%d] adopted=%v, want %v — %s",
						i, adopted, v.ExpectAdopt[i], v.Description)
				}
			}

			// --- verify ---
			policy, err := verify.NewPolicy(
				time.Duration(v.Policy.RSeconds*float64(time.Second)),
				time.Duration(v.Policy.SkewSeconds*float64(time.Second)),
			)
			if err != nil {
				t.Fatalf("policy: %v", err)
			}
			ver, err := verify.NewVerifier(policy, ts, revAdapter{p: provider}, fixedClock{t: now})
			if err != nil {
				t.Fatalf("verifier: %v", err)
			}
			verdict, trace := ver.Verify([]byte(v.Record))

			if got := decisionToken(verdict.Decision); got != v.Expect.Decision {
				t.Fatalf("decision = %s, want %s\n  %s", got, v.Expect.Decision, v.Description)
			}
			if !sameSet(causeNames(verdict.Causes), v.Expect.Causes) {
				t.Fatalf("causes = %v, want %v\n  %s", causeNames(verdict.Causes), v.Expect.Causes, v.Description)
			}
			if len(trace.Entries) != 5 {
				t.Fatalf("trace has %d entries, want 5 unconditionally", len(trace.Entries))
			}
		})
	}
}

// --- ports, reconstructed exactly as a foreign implementation would ---

type staticTrust struct {
	domain spiffeid.TrustDomain
	mat    record.TrustMaterial
}

func (s *staticTrust) TrustMaterialFor(d spiffeid.TrustDomain) (record.TrustMaterial, bool) {
	if d == s.domain {
		return s.mat, true
	}
	return record.TrustMaterial{}, false
}

type revAdapter struct{ p revstatus.Provider }

func (a revAdapter) StatusOf(i record.InstanceID) verify.RevocationStatus {
	ans := a.p.StatusOf(i)
	st := verify.Indeterminate
	switch ans.State {
	case revstatus.NotObservedRevoked:
		st = verify.NotObservedRevoked
	case revstatus.ObservablyRevoked:
		st = verify.ObservablyRevoked
	}
	return verify.RevocationStatus{State: st, AsOf: ans.AsOf}
}

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }
