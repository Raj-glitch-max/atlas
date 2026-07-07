package record

// AT19/AT21-class reconstruction tests (ER4, FR6, INV9, SO8): a third party
// holding only the presented record and trust material — no access to the
// issuance context or any verifier's runtime state — recovers who delegated
// to whom, with what scope, at what time; and reconstruction remains
// possible after the delegation's validity has ended (the record persists
// past Expired/Revoked; RFC-002 §9.1).

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// thirdPartyReconstruct is the reviewer's entire toolkit: bytes + trust
// material in, assertions out. It is a free function over the public
// surface only — it deliberately takes nothing else, mirroring AT19's
// "record alone" stimulus.
func thirdPartyReconstruct(presented []byte, tm TrustMaterial) (Assertions, bool) {
	rec, outcome := ValidateIntegrity(presented, tm)
	if outcome != Intact {
		return Assertions{}, false
	}
	return rec.Read(), true
}

func TestThirdPartyReconstructionFromRecordAlone(t *testing.T) {
	// Issuance side: seal a record, keep only its bytes.
	signer, tm := testSigner(t)
	sealed := mustSeal(t, testAssertions(t), signer)
	presented := append([]byte(nil), sealed.Presented()...)

	// Third-party side: expectations are hardcoded constants, not values
	// shared with the issuance-side variables — the reviewer knows what
	// the record should say, not what the issuer's process state was.
	got, ok := thirdPartyReconstruct(presented, tm)
	if !ok {
		t.Fatal("reconstruction failed on an unaltered record")
	}
	if got.Principal.String() != "spiffe://domain-a.test/principal" {
		t.Errorf("delegator = %s", got.Principal)
	}
	if got.Delegate.String() != "spiffe://domain-a.test/delegate" {
		t.Errorf("delegate = %s", got.Delegate)
	}
	if len(got.Scope) != 2 || got.Scope[0] != "read:orders" || got.Scope[1] != "write:audit" {
		t.Errorf("scope = %q", got.Scope)
	}
	if got.IssuedAt != time.Unix(1_800_000_000, 0).UTC() {
		t.Errorf("issuance time = %v", got.IssuedAt)
	}
	if got.Expiration != time.Unix(1_800_000_600, 0).UTC() {
		t.Errorf("expiration = %v", got.Expiration)
	}
}

func TestReconstructionAfterValidityEnded(t *testing.T) {
	// A record whose validity window is long past still reconstructs:
	// terminal lifecycle states end validity, not the record (INV9).
	signer, tm := testSigner(t)
	a := testAssertions(t)
	a.IssuedAt = time.Unix(1_000_000_000, 0).UTC()   // 2001
	a.Expiration = time.Unix(1_000_003_600, 0).UTC() // expired decades ago
	sealed := mustSeal(t, a, signer)

	got, ok := thirdPartyReconstruct(sealed.Presented(), tm)
	if !ok {
		t.Fatal("an expired record must still reconstruct")
	}
	if got.Expiration != a.Expiration || got.IssuedAt != a.IssuedAt {
		t.Error("reconstructed times differ from issuance")
	}
}

func TestReconstructionHasNoPrivilegedFallback(t *testing.T) {
	// With the wrong trust material, reconstruction fails closed: there
	// is no side channel to the issuer and no privileged path (AP8).
	signer, _ := testSigner(t)
	sealed := mustSeal(t, testAssertions(t), signer)

	otherKey := testKey(t)
	otherDomain := spiffeid.RequireTrustDomainFromString("domain-b.test")
	wrongTM, err := NewTrustMaterial(otherDomain, map[string]*ecdsa.PublicKey{
		"authority-key-1": &otherKey.PublicKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := thirdPartyReconstruct(sealed.Presented(), wrongTM); ok {
		t.Fatal("reconstruction must fail with material that cannot attribute the record")
	}
}

func TestReconstructionIsIndependentOfOriginalVerification(t *testing.T) {
	// AT21: two independent parties validate the same bytes with their
	// own copies of the material; both reconstruct identically, sharing
	// no state — a *Record is never passed between them.
	signer, tm := testSigner(t)
	sealed := mustSeal(t, testAssertions(t), signer)
	presented := sealed.Presented()

	first, ok1 := thirdPartyReconstruct(append([]byte(nil), presented...), tm)
	second, ok2 := thirdPartyReconstruct(append([]byte(nil), presented...), tm)
	if !ok1 || !ok2 {
		t.Fatal("independent reconstructions failed")
	}
	if first.Principal != second.Principal || first.Delegate != second.Delegate ||
		first.Instance != second.Instance ||
		first.Expiration != second.Expiration || first.IssuedAt != second.IssuedAt {
		t.Error("independent reconstructions disagree")
	}
}
