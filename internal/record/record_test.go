package record

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// --- shared test helpers -----------------------------------------------

func testKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating P-256 key: %v", err)
	}
	return key
}

func testSigner(t *testing.T) (Signer, TrustMaterial) {
	t.Helper()
	key := testKey(t)
	signer := Signer{Key: key, KeyID: "authority-key-1"}
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")
	tm, err := NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{
		signer.KeyID: &key.PublicKey,
	})
	if err != nil {
		t.Fatalf("building trust material: %v", err)
	}
	return signer, tm
}

func testAssertions(t *testing.T) Assertions {
	t.Helper()
	instance, err := InstanceIDFromString("inst-0001")
	if err != nil {
		t.Fatalf("instance id: %v", err)
	}
	return Assertions{
		Principal:  spiffeid.RequireFromString("spiffe://domain-a.test/principal"),
		Delegate:   spiffeid.RequireFromString("spiffe://domain-a.test/delegate"),
		Scope:      []string{"read:orders", "write:audit"},
		Expiration: time.Unix(1_800_000_600, 0).UTC(),
		IssuedAt:   time.Unix(1_800_000_000, 0).UTC(),
		Instance:   instance,
	}
}

func mustSeal(t *testing.T, a Assertions, s Signer) *Record {
	t.Helper()
	rec, err := Seal(a, s)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	return rec
}

// --- round-trip and canonicalization ------------------------------------

func TestSealValidateReadRoundTrip(t *testing.T) {
	signer, tm := testSigner(t)
	want := testAssertions(t)
	sealed := mustSeal(t, want, signer)

	validated, outcome := ValidateIntegrity(sealed.Presented(), tm)
	if outcome != Intact {
		t.Fatalf("outcome = %v, want Intact", outcome)
	}
	got := validated.Read()

	if got.Principal != want.Principal || got.Delegate != want.Delegate {
		t.Errorf("identities: got (%s, %s), want (%s, %s)",
			got.Principal, got.Delegate, want.Principal, want.Delegate)
	}
	if len(got.Scope) != 2 || got.Scope[0] != "read:orders" || got.Scope[1] != "write:audit" {
		t.Errorf("scope = %q", got.Scope)
	}
	if !got.Expiration.Equal(want.Expiration) || !got.IssuedAt.Equal(want.IssuedAt) {
		t.Errorf("times: got (%v, %v), want (%v, %v)",
			got.Expiration, got.IssuedAt, want.Expiration, want.IssuedAt)
	}
	if got.Instance != want.Instance {
		t.Errorf("instance = %v, want %v", got.Instance, want.Instance)
	}
	if !got.RevocationBinding.IsAbsent() {
		t.Errorf("revocation binding should be absent, got %d bytes", len(got.RevocationBinding))
	}
}

func TestSealCanonicalizesScope(t *testing.T) {
	signer, _ := testSigner(t)
	a := testAssertions(t)
	a.Scope = []string{"write:audit", "read:orders", "write:audit"}
	sealed := mustSeal(t, a, signer)
	got := sealed.Read().Scope
	if len(got) != 2 || got[0] != "read:orders" || got[1] != "write:audit" {
		t.Errorf("canonical scope = %q, want [read:orders write:audit]", got)
	}
}

func TestSealNormalizesTimePrecisionAndZone(t *testing.T) {
	signer, tm := testSigner(t)
	a := testAssertions(t)
	loc := time.FixedZone("plus5", 5*3600)
	a.Expiration = time.Unix(1_800_000_600, 999_999_999).In(loc)
	a.IssuedAt = time.Unix(1_800_000_000, 123).In(loc)
	sealed := mustSeal(t, a, signer)

	validated, outcome := ValidateIntegrity(sealed.Presented(), tm)
	if outcome != Intact {
		t.Fatalf("outcome = %v, want Intact", outcome)
	}
	got := validated.Read()
	if got.Expiration != time.Unix(1_800_000_600, 0).UTC() {
		t.Errorf("expiration = %v, want normalized UTC second", got.Expiration)
	}
	if got.IssuedAt != time.Unix(1_800_000_000, 0).UTC() {
		t.Errorf("issuedAt = %v, want normalized UTC second", got.IssuedAt)
	}
	// In-memory assertions must equal post-decode assertions exactly.
	if sealed.Read().Expiration != got.Expiration || sealed.Read().IssuedAt != got.IssuedAt {
		t.Error("sealed-side and decoded-side times differ")
	}
}

func TestRevocationBindingRoundTripAndAbsence(t *testing.T) {
	signer, tm := testSigner(t)

	a := testAssertions(t)
	a.RevocationBinding = RevBinding{0x01, 0x02, 0xFF}
	sealed := mustSeal(t, a, signer)
	validated, outcome := ValidateIntegrity(sealed.Presented(), tm)
	if outcome != Intact {
		t.Fatalf("outcome = %v, want Intact", outcome)
	}
	got := validated.Read().RevocationBinding
	if got.IsAbsent() || len(got) != 3 || got[0] != 0x01 || got[2] != 0xFF {
		t.Errorf("binding round-trip = %v", got)
	}

	// Empty non-nil binding normalizes to absent.
	b := testAssertions(t)
	b.RevocationBinding = RevBinding{}
	sealed2 := mustSeal(t, b, signer)
	if !sealed2.Read().RevocationBinding.IsAbsent() {
		t.Error("empty binding should normalize to absent")
	}
}

// --- defensive copies and opacity ----------------------------------------

func TestReadReturnsIndependentCopies(t *testing.T) {
	signer, _ := testSigner(t)
	a := testAssertions(t)
	a.RevocationBinding = RevBinding{9, 9}
	sealed := mustSeal(t, a, signer)

	first := sealed.Read()
	first.Scope[0] = "tampered"
	first.RevocationBinding[0] = 0

	second := sealed.Read()
	if second.Scope[0] != "read:orders" {
		t.Error("caller mutation reached the record's scope")
	}
	if second.RevocationBinding[0] != 9 {
		t.Error("caller mutation reached the record's revocation binding")
	}
}

func TestSealDoesNotRetainCallerSlices(t *testing.T) {
	signer, _ := testSigner(t)
	a := testAssertions(t)
	scope := []string{"read:orders", "write:audit"}
	binding := RevBinding{7}
	a.Scope = scope
	a.RevocationBinding = binding
	sealed := mustSeal(t, a, signer)

	scope[0] = "mutated"
	binding[0] = 0
	got := sealed.Read()
	if got.Scope[0] != "read:orders" || got.RevocationBinding[0] != 7 {
		t.Error("record retained caller-owned slices")
	}
}

func TestInstanceIDOpacityAndEquality(t *testing.T) {
	a, err := InstanceIDFromString("alpha")
	if err != nil {
		t.Fatal(err)
	}
	b, err := InstanceIDFromString("alpha")
	if err != nil {
		t.Fatal(err)
	}
	c, err := InstanceIDFromString("beta")
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Error("equal content must compare equal")
	}
	if a == c {
		t.Error("distinct content must compare unequal")
	}
	if a.IsZero() || !(InstanceID{}).IsZero() {
		t.Error("IsZero misreports")
	}
	if _, err := InstanceIDFromString(""); err == nil {
		t.Error("empty instance id must refuse")
	}
	var m map[InstanceID]bool = map[InstanceID]bool{a: true}
	if !m[b] {
		t.Error("InstanceID must be usable as a map key by equality")
	}
}

func TestOutcomeZeroValueIsAltered(t *testing.T) {
	var o Outcome
	if o != Altered {
		t.Error("the zero Outcome must be Altered (fail-safe)")
	}
	if Altered.String() != "Altered" || Intact.String() != "Intact" {
		t.Error("Outcome.String misrenders")
	}
}

// --- Seal refusals (negative) ---------------------------------------------

func TestSealRefusesIncompleteAssertions(t *testing.T) {
	signer, _ := testSigner(t)
	base := testAssertions(t)

	cases := map[string]func(*Assertions){
		"missing principal":  func(a *Assertions) { a.Principal = spiffeid.ID{} },
		"missing delegate":   func(a *Assertions) { a.Delegate = spiffeid.ID{} },
		"nil scope":          func(a *Assertions) { a.Scope = nil },
		"empty scope":        func(a *Assertions) { a.Scope = []string{} },
		"empty scope entry":  func(a *Assertions) { a.Scope = []string{"read:orders", ""} },
		"missing expiration": func(a *Assertions) { a.Expiration = time.Time{} },
		"missing issuedAt":   func(a *Assertions) { a.IssuedAt = time.Time{} },
		"missing instance":   func(a *Assertions) { a.Instance = InstanceID{} },
	}
	for name, mutate := range cases {
		a := base
		a.Scope = append([]string(nil), base.Scope...)
		mutate(&a)
		if _, err := Seal(a, signer); err == nil {
			t.Errorf("%s: Seal must refuse", name)
		}
	}
}

func TestSealRefusesInvalidSigner(t *testing.T) {
	a := testAssertions(t)

	if _, err := Seal(a, Signer{Key: nil, KeyID: "k"}); err == nil {
		t.Error("nil key must refuse")
	}
	if _, err := Seal(a, Signer{Key: testKey(t), KeyID: ""}); err == nil {
		t.Error("empty key ID must refuse")
	}
	p384, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Seal(a, Signer{Key: p384, KeyID: "k"}); err == nil {
		t.Error("non-P-256 key must refuse (ES256 pin)")
	}
}

// --- TrustMaterial construction (negative + boundary) ----------------------

func TestNewTrustMaterialRefusals(t *testing.T) {
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")
	key := testKey(t)

	if _, err := NewTrustMaterial(spiffeid.TrustDomain{}, map[string]*ecdsa.PublicKey{"k": &key.PublicKey}); err == nil {
		t.Error("zero domain must refuse")
	}
	if _, err := NewTrustMaterial(domain, nil); err == nil {
		t.Error("no keys must refuse")
	}
	if _, err := NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{"": &key.PublicKey}); err == nil {
		t.Error("empty kid must refuse")
	}
	if _, err := NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{"k": nil}); err == nil {
		t.Error("nil key must refuse")
	}
	p384, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{"k": &p384.PublicKey}); err == nil {
		t.Error("non-P-256 key must refuse (ES256 pin)")
	}
}

func TestTrustMaterialCopiesCallerMap(t *testing.T) {
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")
	key := testKey(t)
	keys := map[string]*ecdsa.PublicKey{"k": &key.PublicKey}
	tm, err := NewTrustMaterial(domain, keys)
	if err != nil {
		t.Fatal(err)
	}
	delete(keys, "k")
	if _, held := tm.keyFor("k"); !held {
		t.Error("trust material must not share the caller's map")
	}
	if tm.Domain() != domain {
		t.Errorf("Domain() = %v", tm.Domain())
	}
}

// --- representation boundary: validity is not integrity --------------------

func TestExpiredAtBirthStillSealsAndValidates(t *testing.T) {
	// exp < iat is representable: the verifier rejects it as expired;
	// the record model neither judges nor blocks it (validity is M3's).
	signer, tm := testSigner(t)
	a := testAssertions(t)
	a.Expiration = a.IssuedAt.Add(-time.Hour)
	sealed := mustSeal(t, a, signer)
	if _, outcome := ValidateIntegrity(sealed.Presented(), tm); outcome != Intact {
		t.Error("integrity must be independent of validity")
	}
}

// checkpoint: refactor(sdk): refactor revstatus snapshot retrieval
