package issuance_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/issuance"
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/tests/harness"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

const (
	principalID = "spiffe://domain-a.test/principal"
	delegateID  = "spiffe://domain-a.test/delegate"
	keyID       = "authority-key-1"
)

var (
	issuedAt = time.Unix(1_800_000_000, 0).UTC()
	expiry   = time.Unix(1_800_003_600, 0).UTC()
)

type fixture struct {
	authority *issuance.Authority
	principal spiffeid.ID
	delegate  spiffeid.ID
	pubKey    *ecdsa.PublicKey
}

// newFixture builds an authority whose principal holds a 3-permission set, so
// a 2-permission scope is a proper subset and the full 3 is not (strict).
func newFixture(t *testing.T) *fixture {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	principal := spiffeid.RequireFromString(principalID)
	perms := harness.NewPermissionSource().Grant(principal, "read:orders", "write:audit", "admin:all")
	auth, err := issuance.NewAuthority(
		record.Signer{Key: key, KeyID: keyID},
		perms,
		issuance.NoBinding{},
		issuance.RandomMinter{},
		harness.NewClock(issuedAt),
	)
	if err != nil {
		t.Fatalf("authority: %v", err)
	}
	return &fixture{
		authority: auth,
		principal: principal,
		delegate:  spiffeid.RequireFromString(delegateID),
		pubKey:    &key.PublicKey,
	}
}

func (f *fixture) req(scope ...string) issuance.Request {
	return issuance.Request{
		Principal:  f.principal,
		Delegate:   f.delegate,
		Scope:      scope,
		Expiration: expiry,
	}
}

func TestIssueProperSubsetSucceeds(t *testing.T) {
	f := newFixture(t)
	res, err := f.authority.Issue(f.req("read:orders", "write:audit"))
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if res.Outcome != issuance.Issued || res.Record == nil {
		t.Fatalf("outcome = %s (record nil=%v), want Issued", res.Outcome, res.Record == nil)
	}
	// The issued record is authentic under the authority's key and carries a
	// fresh instance identity and the issuance time.
	tm, err := record.NewTrustMaterial(f.principal.TrustDomain(), map[string]*ecdsa.PublicKey{keyID: f.pubKey})
	if err != nil {
		t.Fatal(err)
	}
	validated, outcome := record.ValidateIntegrity(res.Record.Presented(), tm)
	if outcome != record.Intact {
		t.Fatal("issued record must be authentic under the authority key")
	}
	a := validated.Read()
	if a.Instance.IsZero() {
		t.Error("issued record must carry a fresh instance identity")
	}
	if !a.IssuedAt.Equal(issuedAt) {
		t.Errorf("issued-at = %v, want %v", a.IssuedAt, issuedAt)
	}
	if res.Trace.Outcome != issuance.Issued || !res.Trace.SubsetSatisfied || res.Trace.Instance == "" {
		t.Errorf("trace incomplete: %+v", res.Trace)
	}
}

func TestIssueRefusesOverScope(t *testing.T) {
	f := newFixture(t)
	cases := map[string][]string{
		"reaches outside permissions": {"read:orders", "delete:everything"},
		"equals full permission set":  {"read:orders", "write:audit", "admin:all"}, // not strict
	}
	for name, scope := range cases {
		t.Run(name, func(t *testing.T) {
			res, err := f.authority.Issue(f.req(scope...))
			if err != nil {
				t.Fatalf("issue: %v", err)
			}
			if res.Outcome != issuance.Refused || res.Refusal != issuance.OverScope {
				t.Fatalf("outcome=%s refusal=%s, want Refused/OverScope", res.Outcome, res.Refusal)
			}
			if res.Record != nil {
				t.Error("refusal must create nothing (FM6)")
			}
		})
	}
}

func TestIssueRefusesPermissionsUnavailable(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	principal := spiffeid.RequireFromString(principalID)
	// Permission source knows nothing about this principal.
	perms := harness.NewPermissionSource()
	auth, err := issuance.NewAuthority(record.Signer{Key: key, KeyID: keyID}, perms, issuance.NoBinding{}, issuance.RandomMinter{}, harness.NewClock(issuedAt))
	if err != nil {
		t.Fatal(err)
	}
	res, err := auth.Issue(issuance.Request{
		Principal: principal, Delegate: spiffeid.RequireFromString(delegateID),
		Scope: []string{"read:orders"}, Expiration: expiry,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Refusal != issuance.PermissionsUnavailable {
		t.Fatalf("refusal = %s, want PermissionsUnavailable", res.Refusal)
	}
	if res.Record != nil {
		t.Error("refusal must create nothing")
	}
}

func TestIssueRefusesMalformedRequest(t *testing.T) {
	f := newFixture(t)
	cases := map[string]issuance.Request{
		"missing principal":  {Delegate: f.delegate, Scope: []string{"read:orders"}, Expiration: expiry},
		"missing delegate":   {Principal: f.principal, Scope: []string{"read:orders"}, Expiration: expiry},
		"empty scope":        {Principal: f.principal, Delegate: f.delegate, Scope: nil, Expiration: expiry},
		"empty scope entry":  {Principal: f.principal, Delegate: f.delegate, Scope: []string{"read:orders", ""}, Expiration: expiry},
		"missing expiration": {Principal: f.principal, Delegate: f.delegate, Scope: []string{"read:orders"}},
	}
	for name, req := range cases {
		t.Run(name, func(t *testing.T) {
			res, err := f.authority.Issue(req)
			if err != nil {
				t.Fatalf("issue: %v", err)
			}
			if res.Refusal != issuance.MalformedRequest {
				t.Fatalf("refusal = %s, want MalformedRequest", res.Refusal)
			}
		})
	}
}

func TestIssueMintsUniqueInstances(t *testing.T) {
	f := newFixture(t)
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		res, err := f.authority.Issue(f.req("read:orders"))
		if err != nil || res.Record == nil {
			t.Fatalf("issue %d: %v", i, err)
		}
		id := res.Record.Read().Instance.String()
		if seen[id] {
			t.Fatalf("instance identity collision at issuance %d: %q", i, id)
		}
		seen[id] = true
	}
}

func TestNewAuthorityRefusals(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	good := record.Signer{Key: key, KeyID: keyID}
	perms := harness.NewPermissionSource()
	clock := harness.NewClock(issuedAt)

	if _, err := issuance.NewAuthority(record.Signer{}, perms, issuance.NoBinding{}, issuance.RandomMinter{}, clock); err == nil {
		t.Error("empty signer must refuse")
	}
	if _, err := issuance.NewAuthority(good, nil, issuance.NoBinding{}, issuance.RandomMinter{}, clock); err == nil {
		t.Error("nil permission source must refuse")
	}
	if _, err := issuance.NewAuthority(good, perms, nil, issuance.RandomMinter{}, clock); err == nil {
		t.Error("nil binding source must refuse")
	}
	if _, err := issuance.NewAuthority(good, perms, issuance.NoBinding{}, nil, clock); err == nil {
		t.Error("nil minter must refuse")
	}
	if _, err := issuance.NewAuthority(good, perms, issuance.NoBinding{}, issuance.RandomMinter{}, nil); err == nil {
		t.Error("nil clock must refuse")
	}
}

func TestPermissionSetProperSubsetSemantics(t *testing.T) {
	s := issuance.NewPermissionSet("a", "b", "c")
	if s.Len() != 3 {
		t.Fatalf("len = %d, want 3", s.Len())
	}
	if !s.Contains("a") || s.Contains("z") {
		t.Error("Contains wrong")
	}
	// Duplicates and empties ignored in construction.
	s2 := issuance.NewPermissionSet("a", "a", "")
	if s2.Len() != 1 {
		t.Errorf("dedup/empty handling: len = %d, want 1", s2.Len())
	}
}
