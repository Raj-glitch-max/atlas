package vectors_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/issuance"
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

const issuancePath = "issuance-vectors.json"

var updateIssuance = flag.Bool("update-issuance", false, "regenerate issuance vectors (not implemented - vectors are hand-written)")

// IssuanceVector is one issuance-layer conformance case.
type IssuanceVector struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Principal   string            `json:"principal"`
	Delegate    string            `json:"delegate"`
	Scope       []string          `json:"scope"`
	Permissions map[string][]string `json:"permissions"` // SPIFFE ID -> scopes
	Mode        string            `json:"mode"`          // open, strict, sandbox
	Expect      IssuanceExpect    `json:"expect"`
}

type IssuanceExpect struct {
	Outcome  string `json:"outcome"`  // Issued or Refused
	Refusal  string `json:"refusal"`  // NoRefusal, OverScope, PermissionsUnavailable, MalformedRequest
}

type IssuanceVectorFile struct {
	Schema      int                `json:"schema"`
	Description string             `json:"description"`
	Vectors     []IssuanceVector   `json:"vectors"`
}

func TestIssuanceVectorsReplay(t *testing.T) {
	data, err := os.ReadFile(issuancePath)
	if err != nil {
		t.Fatalf("read issuance vectors: %v", err)
	}
	var vf IssuanceVectorFile
	if err := json.Unmarshal(data, &vf); err != nil {
		t.Fatalf("parse issuance vectors: %v", err)
	}
	if vf.Schema != 1 {
		t.Fatalf("issuance vector schema %d, code expects 1", vf.Schema)
	}
	if len(vf.Vectors) == 0 {
		t.Fatal("no issuance vectors")
	}

	for _, vec := range vf.Vectors {
		vec := vec
		t.Run(vec.Name, func(t *testing.T) {
			// Build permission source from vector
			permSrc := buildPermissionSource(t, vec.Permissions)

			// Build authority
			key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				t.Fatal(err)
			}
			const keyID = "authority-key-1"
			now := time.Unix(1_800_000_000, 0).UTC()
			expiry := now.Add(time.Hour)

			auth, err := issuance.NewAuthority(
				record.Signer{Key: key, KeyID: keyID},
				permSrc,
				issuance.NoBinding{},
				issuance.RandomMinter{},
				fixedClock{now},
			)
			if err != nil {
				t.Fatalf("authority: %v", err)
			}

			principal := spiffeid.RequireFromString(vec.Principal)
			delegate := spiffeid.RequireFromString(vec.Delegate)

			res, err := auth.Issue(issuance.Request{
				Principal:  principal,
				Delegate:   delegate,
				Scope:      vec.Scope,
				Expiration: expiry,
			})
			if err != nil {
				t.Fatalf("issue: %v", err)
			}

			if string(res.Outcome) != vec.Expect.Outcome {
				t.Fatalf("outcome = %s, expect %s", res.Outcome, vec.Expect.Outcome)
			}
			if string(res.Refusal) != vec.Expect.Refusal {
				t.Fatalf("refusal = %s, expect %s", res.Refusal, vec.Expect.Refusal)
			}

			// If Issued, verify the record is valid
			if res.Outcome == issuance.Issued {
				if res.Record == nil {
					t.Fatal("Issued but Record is nil")
				}
				// Verify the record can be validated
				pubKey := &key.PublicKey
				tm, err := record.NewTrustMaterial(principal.TrustDomain(), map[string]*ecdsa.PublicKey{keyID: pubKey})
				if err != nil {
					t.Fatalf("trust material: %v", err)
				}
				validated, outcome := record.ValidateIntegrity(res.Record.Presented(), tm)
				if outcome != record.Intact {
					t.Fatalf("issued record not authentic: %s", outcome)
				}
				assertions := validated.Read()
				if assertions.Instance.IsZero() {
					t.Error("issued record must carry instance identity")
				}
				if assertions.Principal != principal {
					t.Errorf("principal mismatch: %v != %v", assertions.Principal, principal)
				}
				if assertions.Delegate != delegate {
					t.Errorf("delegate mismatch: %v != %v", assertions.Delegate, delegate)
				}
			} else {
				if res.Record != nil {
					t.Error("refusal must create nothing (FM6)")
				}
			}
		})
	}
}

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

func buildPermissionSource(t *testing.T, perms map[string][]string) issuance.PermissionSource {
	t.Helper()
	if perms == nil {
		return nil
	}
	m := make(map[spiffeid.ID]issuance.PermissionSet, len(perms))
	for idStr, scopes := range perms {
		id, err := spiffeid.FromString(idStr)
		if err != nil {
			t.Fatalf("invalid SPIFFE ID %q: %v", idStr, err)
		}
		m[id] = issuance.NewPermissionSet(scopes...)
	}
	return testPermissionSource{perms: m}
}

type testPermissionSource struct {
	perms map[spiffeid.ID]issuance.PermissionSet
}

func (t testPermissionSource) PermissionsOf(principal spiffeid.ID) (issuance.PermissionSet, bool) {
	set, ok := t.perms[principal]
	return set, ok
}