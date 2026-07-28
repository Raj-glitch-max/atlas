package conformance

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	jose "github.com/go-jose/go-jose/v3"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Revocation-layer conformance vectors.
//
// WHY THIS FILE EXISTS, separately from the verdict/negative vectors:
//
// The schema-1 vectors supply `revocation.state` as an ALREADY-COMPUTED answer
// (Indeterminate | NotObservedRevoked | ObservablyRevoked). That tests the
// verifier's HANDLING of an answer and completely skips its DERIVATION — the
// snapshot's signature is never checked, and freshness is never computed, by
// anything in the language-neutral suite.
//
// Freshness-bounded, fail-closed revocation is the most novel mechanism in
// Atlas and was consequently its least conformance-tested one: an independent
// implementation could pass all 30 existing vectors while getting revocation
// semantics entirely wrong.
//
// These vectors carry the SIGNED SNAPSHOT itself, so a foreign implementation
// must perform the real work: verify the snapshot signature with the published
// key, reject a snapshot whose list id does not match, enforce monotonic
// adoption against rollback, compute the snapshot's age against the verifier's
// clock, and only then derive a revocation state and a verdict.

// RevocationSchemaVersion is independent of SchemaVersion: this is a distinct
// artifact with a distinct shape, so it carries its own version rather than
// forcing a bump on the verdict/negative files.
const RevocationSchemaVersion = 1

// VectorSnapshot is a signed revoked-set snapshot exactly as it travels in a
// trust bundle. Every field an implementation needs to verify it is present;
// nothing is pre-digested.
type VectorSnapshot struct {
	ListID  string   `json:"list_id"`
	AsOf    string   `json:"as_of"`   // RFC3339 UTC — the currency of this knowledge
	Revoked []string `json:"revoked"` // instance ids in the set
	Sig     string   `json:"sig"`     // base64 (std) signature over the set
}

// RevocationVector is one revocation-layer case.
//
// Ingest is a SEQUENCE, deliberately: rollback is only expressible as an
// ordered pair of ingests, and a verifier that adopts snapshots without
// enforcing monotonicity can only be caught by presenting a newer snapshot
// followed by an older one.
type RevocationVector struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Record      string           `json:"record"`       // the presented record (compact JWS)
	Trust       VectorTrust      `json:"trust"`        // RP-held trust material for the record
	SnapshotKey json.RawMessage  `json:"snapshot_key"` // JWK the RP holds for the revocation origin
	ExpectList  string           `json:"expect_list"`  // the list id the RP expects; a mismatch must be refused
	Ingest      []VectorSnapshot `json:"ingest"`       // applied in order, before verification
	ExpectAdopt []bool           `json:"expect_adopt"` // per-ingest: must this snapshot be adopted?
	Now         string           `json:"now"`          // RFC3339 UTC — the verifier's clock
	Policy      VectorPolicy     `json:"policy"`       // R and skew
	Expect      VectorExpect     `json:"expect"`       // the required verdict
}

// RevocationVectorFile is the committed artifact.
type RevocationVectorFile struct {
	Schema      int                `json:"schema"`
	Description string             `json:"description"`
	Vectors     []RevocationVector `json:"vectors"`
}

func mustID(s string) spiffeid.ID {
	return spiffeid.RequireFromString(s)
}

func jwkOf(pub *ecdsa.PublicKey, kid string) (json.RawMessage, error) {
	k := jose.JSONWebKey{Key: pub, KeyID: kid, Algorithm: string(jose.ES256), Use: "sig"}
	return k.MarshalJSON()
}

func snapshotToVector(s revstatus.SignedRevokedSet) VectorSnapshot {
	ids := make([]string, 0, len(s.Revoked))
	for _, r := range s.Revoked {
		ids = append(ids, r.String())
	}
	return VectorSnapshot{
		ListID:  s.ListID,
		AsOf:    s.AsOf.UTC().Format(time.RFC3339Nano),
		Revoked: ids,
		Sig:     base64.StdEncoding.EncodeToString(s.Sig),
	}
}

// SnapshotFromVector rebuilds a snapshot from its published form — the exact
// operation a foreign implementation performs. Exported so the replay test and
// any external harness use the same path.
func SnapshotFromVector(v VectorSnapshot) (revstatus.SignedRevokedSet, error) {
	asOf, err := time.Parse(time.RFC3339Nano, v.AsOf)
	if err != nil {
		return revstatus.SignedRevokedSet{}, fmt.Errorf("as_of: %w", err)
	}
	sig, err := base64.StdEncoding.DecodeString(v.Sig)
	if err != nil {
		return revstatus.SignedRevokedSet{}, fmt.Errorf("sig: %w", err)
	}
	ids := make([]record.InstanceID, 0, len(v.Revoked))
	for _, s := range v.Revoked {
		id, err := record.InstanceIDFromString(s)
		if err != nil {
			return revstatus.SignedRevokedSet{}, fmt.Errorf("revoked id %q: %w", s, err)
		}
		ids = append(ids, id)
	}
	return revstatus.SignedRevokedSet{
		ListID: v.ListID, AsOf: asOf, Revoked: ids, Sig: sig,
	}, nil
}

// EmitRevocationVectors builds the revocation-layer corpus.
//
// Each case is constructed from real signed snapshots, so the committed file
// contains genuine cryptographic material rather than assertions about it.
func EmitRevocationVectors() (RevocationVectorFile, error) {
	// Authority key (signs records) and revocation-origin key (signs snapshots)
	// are deliberately DISTINCT here, even though the reference server uses one
	// key for both. Separating them in the vectors prevents an implementation
	// from accidentally passing by verifying a snapshot with the record key.
	authKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return RevocationVectorFile{}, err
	}
	revKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return RevocationVectorFile{}, err
	}
	otherKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return RevocationVectorFile{}, err
	}

	const (
		domain = "domain-a.test"
		listID = "atlas-revlist:domain-a.test"
		kid    = "authority-key-1"
		revKid = "revocation-origin-1"
	)

	base := time.Unix(1_800_000_000, 0).UTC()

	// A record valid across the whole window, so every verdict below is
	// determined by the revocation layer alone and never by expiry.
	inst, err := record.InstanceIDFromString("inst-revocation-vectors")
	if err != nil {
		return RevocationVectorFile{}, err
	}
	sealed, err := record.Seal(record.Assertions{
		Principal:  mustID("spiffe://" + domain + "/workload/payments-api"),
		Delegate:   mustID("spiffe://domain-b.test/agent/booking-worker"),
		Scope:      []string{"read:orders"},
		Expiration: base.Add(24 * time.Hour),
		IssuedAt:   base.Add(-time.Minute),
		Instance:   inst,
	}, record.Signer{Key: authKey, KeyID: kid})
	if err != nil {
		return RevocationVectorFile{}, err
	}

	authJWK, err := jwkOf(&authKey.PublicKey, kid)
	if err != nil {
		return RevocationVectorFile{}, err
	}
	revJWK, err := jwkOf(&revKey.PublicKey, revKid)
	if err != nil {
		return RevocationVectorFile{}, err
	}

	pub, err := revstatus.NewPublisher(revKey, listID)
	if err != nil {
		return RevocationVectorFile{}, err
	}
	otherPub, err := revstatus.NewPublisher(otherKey, listID)
	if err != nil {
		return RevocationVectorFile{}, err
	}
	wrongListPub, err := revstatus.NewPublisher(revKey, "atlas-revlist:someone-else.test")
	if err != nil {
		return RevocationVectorFile{}, err
	}

	empty, err := pub.Publish(nil, base)
	if err != nil {
		return RevocationVectorFile{}, err
	}
	revoked, err := pub.Publish([]record.InstanceID{inst}, base)
	if err != nil {
		return RevocationVectorFile{}, err
	}
	revokedLater, err := pub.Publish([]record.InstanceID{inst}, base.Add(10*time.Second))
	if err != nil {
		return RevocationVectorFile{}, err
	}
	emptyEarlier, err := pub.Publish(nil, base.Add(-10*time.Second))
	if err != nil {
		return RevocationVectorFile{}, err
	}
	forged, err := otherPub.Publish(nil, base)
	if err != nil {
		return RevocationVectorFile{}, err
	}
	wrongList, err := wrongListPub.Publish(nil, base)
	if err != nil {
		return RevocationVectorFile{}, err
	}

	// A snapshot whose signature is corrupted after signing.
	tampered := snapshotToVector(empty)
	sigRaw, _ := base64.StdEncoding.DecodeString(tampered.Sig)
	if len(sigRaw) > 0 {
		sigRaw[len(sigRaw)-1] ^= 0x01
	}
	tampered.Sig = base64.StdEncoding.EncodeToString(sigRaw)

	// A snapshot with an entry stripped after signing — the classic
	// "hide the revocation" tamper.
	stripped := snapshotToVector(revoked)
	stripped.Revoked = nil

	trust := VectorTrust{Domain: domain, Keys: []json.RawMessage{authJWK}}
	policy := VectorPolicy{RSeconds: 30, SkewSeconds: 30}
	rec := string(sealed.Presented())

	v := func(name, desc string, ingest []VectorSnapshot, adopt []bool, now time.Time, decision string, causes ...string) RevocationVector {
		if causes == nil {
			causes = []string{}
		}
		return RevocationVector{
			Name: name, Description: desc, Record: rec, Trust: trust,
			SnapshotKey: revJWK, ExpectList: listID,
			Ingest: ingest, ExpectAdopt: adopt,
			Now:    now.UTC().Format(time.RFC3339Nano),
			Policy: policy,
			Expect: VectorExpect{Decision: decision, Causes: causes},
		}
	}

	vectors := []RevocationVector{
		v("fresh-not-revoked",
			"A validly-signed snapshot within R that does not list the instance. The only case that may accept.",
			[]VectorSnapshot{snapshotToVector(empty)}, []bool{true},
			base.Add(5*time.Second), "Accept"),

		v("fresh-revoked",
			"A validly-signed snapshot within R listing the instance. Definitive rejection.",
			[]VectorSnapshot{snapshotToVector(revoked)}, []bool{true},
			base.Add(5*time.Second), "Reject", "RevokedObservable"),

		v("stale-not-revoked-fails-closed",
			"The instance is NOT revoked, but the snapshot is older than R. The verifier cannot prove freshness and MUST refuse rather than accept. This is the fail-closed property; an implementation that accepts here is unsafe even though its answer happens to be right.",
			[]VectorSnapshot{snapshotToVector(empty)}, []bool{true},
			base.Add(90*time.Second), "InconclusiveRejected", "RevocationKnowledgeStale"),

		v("no-snapshot-is-indeterminate",
			"No snapshot was ever ingested. Absence of a revocation record is NOT evidence of non-revocation.",
			nil, nil,
			base.Add(5*time.Second), "InconclusiveRejected", "RevocationStatusIndeterminate"),

		v("bad-signature-not-adopted",
			"A snapshot whose signature was corrupted after signing MUST NOT be adopted, leaving the verifier with no knowledge — not with the snapshot's contents.",
			[]VectorSnapshot{tampered}, []bool{false},
			base.Add(5*time.Second), "InconclusiveRejected", "RevocationStatusIndeterminate"),

		v("stripped-entry-not-adopted",
			"An attacker removes the instance from the revoked set to hide a revocation. The signature covers the set, so the snapshot MUST NOT be adopted — and the verifier must NOT fall back to accepting.",
			[]VectorSnapshot{stripped}, []bool{false},
			base.Add(5*time.Second), "InconclusiveRejected", "RevocationStatusIndeterminate"),

		v("foreign-key-not-adopted",
			"A well-formed snapshot signed by a key the RP does not hold for the revocation origin. Correct signature, wrong signer.",
			[]VectorSnapshot{snapshotToVector(forged)}, []bool{false},
			base.Add(5*time.Second), "InconclusiveRejected", "RevocationStatusIndeterminate"),

		v("wrong-list-id-not-adopted",
			"A validly-signed snapshot from a DIFFERENT revocation stream. Binding to the expected list id is what stops one domain's revocation state being substituted for another's.",
			[]VectorSnapshot{snapshotToVector(wrongList)}, []bool{false},
			base.Add(5*time.Second), "InconclusiveRejected", "RevocationStatusIndeterminate"),

		v("rollback-rejected",
			"A newer snapshot listing the instance, then an OLDER validly-signed snapshot that does not. Both signatures are genuine. Adopting the second would un-revoke the capability, so monotonic adoption MUST reject it and the verdict must remain Reject.",
			[]VectorSnapshot{snapshotToVector(revokedLater), snapshotToVector(emptyEarlier)},
			[]bool{true, false},
			base.Add(15*time.Second), "Reject", "RevokedObservable"),

		v("replay-same-asof-rejected",
			"Re-ingesting the identical snapshot must not be adopted a second time: adoption is strictly monotonic in as-of, not merely non-decreasing.",
			[]VectorSnapshot{snapshotToVector(revoked), snapshotToVector(revoked)},
			[]bool{true, false},
			base.Add(5*time.Second), "Reject", "RevokedObservable"),

		v("newer-snapshot-supersedes",
			"A revoked snapshot followed by a genuinely newer one that still lists the instance. The newer one is adopted and revocation remains terminal.",
			[]VectorSnapshot{snapshotToVector(revoked), snapshotToVector(revokedLater)},
			[]bool{true, true},
			base.Add(15*time.Second), "Reject", "RevokedObservable"),

		v("boundary-exactly-at-R",
			"A snapshot exactly R old. The boundary is inclusive: age == R is still fresh. Pinned because an off-by-one here silently changes every deployment's revocation window.",
			[]VectorSnapshot{snapshotToVector(empty)}, []bool{true},
			base.Add(30*time.Second), "Accept"),
	}

	return RevocationVectorFile{
		Schema: RevocationSchemaVersion,
		Description: "Atlas revocation-layer conformance vectors. Unlike the verdict/negative " +
			"vectors, these carry the SIGNED SNAPSHOT rather than a pre-computed revocation " +
			"state: a conformant implementation must verify the snapshot signature, bind it to " +
			"the expected list id, enforce monotonic adoption, compute the snapshot's age " +
			"against its own clock and the policy bound R, and only then derive a verdict.",
		Vectors: vectors,
	}, nil
}
