package revorigin_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revorigin"
)

func inst(t *testing.T, s string) record.InstanceID {
	t.Helper()
	id, err := record.InstanceIDFromString(s)
	if err != nil {
		t.Fatalf("instance %q: %v", s, err)
	}
	return id
}

var t0 = time.Unix(1_800_000_000, 0).UTC()

func TestRevokeAndIsRevoked(t *testing.T) {
	r := revorigin.New()
	a, b := inst(t, "a"), inst(t, "b")

	if r.IsRevoked(a) {
		t.Fatal("nothing revoked yet")
	}
	r.Revoke(a, t0)
	if !r.IsRevoked(a) {
		t.Error("a must be revoked")
	}
	if r.IsRevoked(b) {
		t.Error("revoking a must not affect b (INV6 sibling independence)")
	}
}

func TestRevokeIsTerminalAndIdempotent(t *testing.T) {
	r := revorigin.New()
	a := inst(t, "a")

	r.Revoke(a, t0)
	r.Revoke(a, t0.Add(time.Hour)) // repeat with a later time
	r.Revoke(a, t0.Add(2*time.Hour))

	if r.Len() != 1 {
		t.Fatalf("terminal revocation must record once, got %d entries", r.Len())
	}
	view := r.View()
	if !view[0].RecordedAt.Equal(t0) {
		t.Errorf("original RecordedAt must stand (INV4 terminality), got %v", view[0].RecordedAt)
	}
}

func TestViewIsOrderedAndImmutableSnapshot(t *testing.T) {
	r := revorigin.New()
	ids := []string{"a", "b", "c", "d"}
	for i, s := range ids {
		r.Revoke(inst(t, s), t0.Add(time.Duration(i)*time.Second))
	}

	view := r.View()
	if len(view) != len(ids) {
		t.Fatalf("view has %d entries, want %d", len(view), len(ids))
	}
	for i, s := range ids {
		if view[i].Instance != inst(t, s) {
			t.Errorf("entry %d = %v, want %s (append order)", i, view[i].Instance, s)
		}
	}
	// Mutating the returned snapshot must not affect the register.
	view[0] = revorigin.Entry{}
	if r.View()[0].Instance != inst(t, "a") {
		t.Error("caller mutated the register through the View snapshot")
	}
}

func TestRevokeZeroInstanceIsInert(t *testing.T) {
	r := revorigin.New()
	r.Revoke(record.InstanceID{}, t0)
	if r.Len() != 0 {
		t.Error("a zero instance identifies nothing; recording it must be a no-op")
	}
}

func TestNeverIssuedInstanceRevocationIsRecordedButInert(t *testing.T) {
	// The register does not validate existence — it records whatever
	// non-zero instance it is given. A revocation of a never-issued
	// instance is harmless: no presented record carries that identity, so
	// verification never matches it (documented behavior, not a bug).
	r := revorigin.New()
	ghost := inst(t, "never-issued")
	r.Revoke(ghost, t0)
	if !r.IsRevoked(ghost) {
		t.Error("register records any non-zero instance it is given")
	}
}

func TestConcurrentRevokeAndView(t *testing.T) {
	// Single-writer/multi-reader posture (AD-017): concurrent Revoke and
	// View must not race or corrupt the register. Run under -race.
	r := revorigin.New()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r.Revoke(inst(t, string(rune('A'+i%26))+string(rune('0'+i/26))), t0)
		}(i)
	}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.View()
		}()
	}
	wg.Wait()

	if r.Len() != 50 {
		t.Errorf("expected 50 distinct revocations, got %d", r.Len())
	}
}

// checkpoint: feat(stores): add attenuation rule engine

// checkpoint: chore(internal): harden test assertions (#108)
