package bench_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// Self-contained rig (no test scaffolding): a real key, trust material, a
// sealed record, and the real verifier. Inline ports keep the lab decoupled
// from tests/harness.

type freshRev struct {
	state verify.RevocationState
	now   time.Time
}

func (r freshRev) StatusOf(record.InstanceID) verify.RevocationStatus {
	if r.state == verify.Indeterminate {
		return verify.RevocationStatus{State: verify.Indeterminate}
	}
	return verify.RevocationStatus{State: r.state, AsOf: r.now}
}

type trustAll struct {
	domain spiffeid.TrustDomain
	mat    record.TrustMaterial
}

func (t trustAll) TrustMaterialFor(d spiffeid.TrustDomain) (record.TrustMaterial, bool) {
	if d == t.domain {
		return t.mat, true
	}
	return record.TrustMaterial{}, false
}

type clk struct{ t time.Time }

func (c clk) Now() time.Time { return c.t }

var benchNow = time.Unix(1_800_000_300, 0).UTC()

type rig struct {
	signer    record.Signer
	trust     trustAll
	domain    spiffeid.TrustDomain
	principal spiffeid.ID
	delegate  spiffeid.ID
	record    []byte
	instance  record.InstanceID
}

func newRig(b testing.TB) rig {
	b.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	domain := spiffeid.RequireTrustDomainFromString("domain-a.test")
	mat, err := record.NewTrustMaterial(domain, map[string]*ecdsa.PublicKey{"k1": &key.PublicKey})
	if err != nil {
		b.Fatal(err)
	}
	signer := record.Signer{Key: key, KeyID: "k1"}
	inst, _ := record.InstanceIDFromString("inst-bench")
	rec, err := record.Seal(record.Assertions{
		Principal:  spiffeid.RequireFromString("spiffe://domain-a.test/principal"),
		Delegate:   spiffeid.RequireFromString("spiffe://domain-a.test/delegate"),
		Scope:      []string{"read:orders", "write:audit"},
		Expiration: benchNow.Add(time.Hour),
		IssuedAt:   benchNow.Add(-time.Minute),
		Instance:   inst,
	}, signer)
	if err != nil {
		b.Fatal(err)
	}
	return rig{
		signer: signer, trust: trustAll{domain: domain, mat: mat}, domain: domain,
		principal: spiffeid.RequireFromString("spiffe://domain-a.test/principal"),
		delegate:  spiffeid.RequireFromString("spiffe://domain-a.test/delegate"),
		record:    rec.Presented(), instance: inst,
	}
}

func (r rig) verifier(b testing.TB, state verify.RevocationState) *verify.Verifier {
	b.Helper()
	policy, err := verify.NewPolicy(2*time.Second, 30*time.Second)
	if err != nil {
		b.Fatal(err)
	}
	v, err := verify.NewVerifier(policy, r.trust, freshRev{state: state, now: benchNow}, clk{benchNow})
	if err != nil {
		b.Fatal(err)
	}
	return v
}

func BenchmarkIssue(b *testing.B) {
	r := newRig(b)
	a := record.Assertions{
		Principal: r.principal, Delegate: r.delegate,
		Scope:      []string{"read:orders", "write:audit"},
		Expiration: benchNow.Add(time.Hour), IssuedAt: benchNow, Instance: r.instance,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := record.Seal(a, r.signer); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateIntegrity(b *testing.B) {
	r := newRig(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, out := record.ValidateIntegrity(r.record, r.trust.mat); out != record.Intact {
			b.Fatal("not intact")
		}
	}
}

func BenchmarkVerifyAccept(b *testing.B) {
	r := newRig(b)
	v := r.verifier(b, verify.NotObservedRevoked)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if verdict, _ := v.Verify(r.record); !verdict.IsAccept() {
			b.Fatalf("want Accept, got %s", verdict.Decision)
		}
	}
}

func BenchmarkVerifyRevoked(b *testing.B) {
	r := newRig(b)
	v := r.verifier(b, verify.ObservablyRevoked)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if verdict, _ := v.Verify(r.record); verdict.Decision != verify.Reject {
			b.Fatalf("want Reject, got %s", verdict.Decision)
		}
	}
}

// BenchmarkSignedSetPublish measures the cost of publishing a signed revoked
// set as the number of revoked instances scales — the revocation-artifact
// production cost of the alpha-path realization.
func benchmarkSignedSetPublish(b *testing.B, n int) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := revstatus.NewPublisher(key, "bench-list")
	revoked := make([]record.InstanceID, n)
	for i := range revoked {
		revoked[i], _ = record.InstanceIDFromString("inst-" + itoa(i))
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := pub.Publish(revoked, benchNow); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSignedSetPublish100(b *testing.B)   { benchmarkSignedSetPublish(b, 100) }
func BenchmarkSignedSetPublish1000(b *testing.B)  { benchmarkSignedSetPublish(b, 1000) }
func BenchmarkSignedSetPublish10000(b *testing.B) { benchmarkSignedSetPublish(b, 10000) }

// itoa avoids strconv import churn in a hot helper.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}

// checkpoint: chore(security): document simulated agent node
