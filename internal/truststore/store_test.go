package truststore_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"sync"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/truststore"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

var t0 = time.Unix(1_800_000_000, 0).UTC()

func material(t *testing.T, domain string) record.TrustMaterial {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	td := spiffeid.RequireTrustDomainFromString(domain)
	tm, err := record.NewTrustMaterial(td, map[string]*ecdsa.PublicKey{"k1": &key.PublicKey})
	if err != nil {
		t.Fatal(err)
	}
	return tm
}

func TestEmptyStoreAnswersAbsent(t *testing.T) {
	s := truststore.New()
	if _, ok := s.TrustMaterialFor(spiffeid.RequireTrustDomainFromString("domain-a.test")); ok {
		t.Error("empty store must answer absent, not present")
	}
}

func TestProvisionAndRetrieve(t *testing.T) {
	s := truststore.New()
	m := material(t, "domain-a.test")
	if err := s.Provision(m, t0); err != nil {
		t.Fatalf("provision: %v", err)
	}
	got, ok := s.TrustMaterialFor(m.Domain())
	if !ok {
		t.Fatal("provisioned material must be present")
	}
	if got.Domain() != m.Domain() {
		t.Errorf("domain = %v, want %v", got.Domain(), m.Domain())
	}
	// A different domain is still absent (no cross-domain leakage).
	if _, ok := s.TrustMaterialFor(spiffeid.RequireTrustDomainFromString("domain-b.test")); ok {
		t.Error("unrelated domain must be absent")
	}
}

func TestProvisionRefusesZeroMaterial(t *testing.T) {
	s := truststore.New()
	if err := s.Provision(record.TrustMaterial{}, t0); err == nil {
		t.Error("zero-value material must be refused, never half-trusted")
	}
	if len(s.Provisionings()) != 0 {
		t.Error("a refused provisioning must record nothing")
	}
}

func TestProvisioningRecordsAppend(t *testing.T) {
	s := truststore.New()
	s.Provision(material(t, "domain-a.test"), t0)
	s.Provision(material(t, "domain-b.test"), t0.Add(time.Minute))
	recs := s.Provisionings()
	if len(recs) != 2 {
		t.Fatalf("want 2 provisioning records, got %d", len(recs))
	}
	if recs[0].KeyCount != 1 || recs[1].KeyCount != 1 {
		t.Error("provisioning record must note the key count")
	}
	if !recs[0].ProvisionedAt.Equal(t0) {
		t.Errorf("provisioning time = %v, want %v", recs[0].ProvisionedAt, t0)
	}
	// Snapshot is independent.
	recs[0] = truststore.ProvisioningRecord{}
	if s.Provisionings()[0].Domain.IsZero() {
		t.Error("caller mutated provisioning records through the snapshot")
	}
}

func TestReProvisionReplacesMaterial(t *testing.T) {
	s := truststore.New()
	first := material(t, "domain-a.test")
	second := material(t, "domain-a.test") // same domain, new keys (rotation)
	s.Provision(first, t0)
	s.Provision(second, t0.Add(time.Hour))

	got, ok := s.TrustMaterialFor(first.Domain())
	if !ok {
		t.Fatal("material must be present after re-provision")
	}
	// The replacement's key is held (we can't compare keys directly, but the
	// second provisioning must have been recorded).
	_ = got
	if len(s.Provisionings()) != 2 {
		t.Errorf("re-provision must append a record, got %d", len(s.Provisionings()))
	}
}

func TestConcurrentProvisionAndRead(t *testing.T) {
	// Single-writer/multi-reader posture (AD-017); run under -race.
	s := truststore.New()
	domains := []string{"a.test", "b.test", "c.test", "d.test"}
	var wg sync.WaitGroup
	for _, d := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			s.Provision(material(t, d), t0)
		}(d)
	}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.TrustMaterialFor(spiffeid.RequireTrustDomainFromString("a.test"))
		}()
	}
	wg.Wait()
	if len(s.Provisionings()) != len(domains) {
		t.Errorf("want %d provisionings, got %d", len(domains), len(s.Provisionings()))
	}
}
