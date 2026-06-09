package conformance

import (
	"encoding/json"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/verify"
	jose "github.com/go-jose/go-jose/v3"
)

// Language-neutral conformance vectors: the executable corpus serialized to
// static JSON any implementation — Go, Rust, Python, anything — can replay
// without importing this module. This is the artifact that makes "independent
// implementations" real: a foreign verifier reads the vectors, reproduces the
// verdict for each, and is conformant iff it matches. It is the interop
// backbone every serious trust/crypto primitive ships (Wycheproof, the JOSE
// cookbook, WebAuthn/FIDO conformance, the Ed25519 RFC vectors).
//
// Single source of truth: EmitVectors derives every vector from BuildCorpus
// and records the *actual* verdict the reference verifier (verify) produces,
// so the vectors are the reference behavior by construction. TestV1Conformance
// separately asserts that behavior matches the spec, so the chain is
// spec -> reference -> vectors, with no drift.

// Vector is one language-neutral conformance case.
type Vector struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Record      string           `json:"record"`     // compact JWS (the presented delegation record)
	Trust       VectorTrust      `json:"trust"`      // trust material the relying party holds
	Revocation  VectorRevocation `json:"revocation"` // the revocation-status answer supplied to the verifier
	Now         string           `json:"now"`        // RFC3339 UTC — the verifier's clock reading
	Policy      VectorPolicy     `json:"policy"`     // R and skew tolerance
	Expect      VectorExpect     `json:"expect"`     // the verdict every conformant verifier must produce
}

// VectorTrust is the RP-held trust material: a set of JWK public keys for a
// trust domain. An empty Keys set models "no trust material for the domain."
type VectorTrust struct {
	Domain string            `json:"domain"`
	Keys   []json.RawMessage `json:"keys"` // each an RFC 7517 JWK (EC P-256, alg ES256)
}

// VectorRevocation is the revocation-observation answer supplied to the
// verifier (the M5 output the RP would have obtained locally).
type VectorRevocation struct {
	State string  `json:"state"`           // Indeterminate | NotObservedRevoked | ObservablyRevoked
	AsOf  *string `json:"as_of,omitempty"` // RFC3339 UTC; null for Indeterminate
}

// VectorPolicy is the verifier's freshness/skew policy.
type VectorPolicy struct {
	RSeconds    float64 `json:"r_seconds"`
	SkewSeconds float64 `json:"skew_seconds"`
}

// VectorExpect is the required verdict.
type VectorExpect struct {
	Decision string   `json:"decision"` // Accept | Reject | InconclusiveRejected
	Causes   []string `json:"causes"`   // cause names (see internal/verify Cause)
}

// SchemaVersion is bumped when the vector JSON shape changes (never
// repurpose a field; add and bump — mirrors the interface-spec discipline).
const SchemaVersion = 1

// VectorFile is the top-level committed artifact.
type VectorFile struct {
	Schema      int      `json:"schema"`
	Description string   `json:"description"`
	Vectors     []Vector `json:"vectors"`
}

func decisionToken(d verify.Decision) string {
	switch d {
	case verify.Accept:
		return "Accept"
	case verify.Reject:
		return "Reject"
	default:
		return "InconclusiveRejected"
	}
}

func (k *Kit) publicJWK() (json.RawMessage, error) {
	jwk := jose.JSONWebKey{Key: &k.signer.Key.PublicKey, KeyID: k.signer.KeyID, Algorithm: "ES256", Use: "sig"}
	return jwk.MarshalJSON()
}

// EmitVectors derives the full vector set from the conformance corpus,
// recording the reference verifier's actual verdict for each. Records carry
// real signatures (verification is deterministic; only signing is
// randomized), so the emitted set is stable to replay though a fresh emit
// produces fresh signatures — regenerate deliberately, commit the result.
func EmitVectors() (VectorFile, error) {
	k, err := NewKit()
	if err != nil {
		return VectorFile{}, err
	}
	jwk, err := k.publicJWK()
	if err != nil {
		return VectorFile{}, err
	}

	var vectors []Vector
	for _, s := range BuildCorpus(k) {
		assertions, ok := k.Decode(s.Presented)

		trust := VectorTrust{Domain: k.domain.String(), Keys: []json.RawMessage{}}
		if _, present := s.Trust.TrustMaterialFor(k.domain); present {
			trust.Keys = []json.RawMessage{jwk}
		}

		var rev VectorRevocation
		if ok {
			st := s.Revocation.StatusOf(assertions.Instance)
			rev.State = st.State.String()
			if st.State != verify.Indeterminate {
				a := st.AsOf.UTC().Format(time.RFC3339)
				rev.AsOf = &a
			}
		} else {
			// Record could not be decoded even with the authority key: it is
			// intrinsically malformed; revocation is not consulted. Report the
			// answer the port would give for the zero instance, honestly.
			rev.State = verify.Indeterminate.String()
		}

		// Authoritative expected verdict = what the reference verifier does.
		v, err := verify.NewVerifier(s.Policy, s.Trust, s.Revocation, s.Clock)
		if err != nil {
			return VectorFile{}, err
		}
		verdict, _ := v.Verify(s.Presented)
		causes := make([]string, 0, len(verdict.Causes))
		for _, c := range verdict.Causes {
			causes = append(causes, c.String())
		}

		vectors = append(vectors, Vector{
			Name:        s.Name,
			Description: "reference verdict for the corpus scenario of the same name",
			Record:      string(s.Presented),
			Trust:       trust,
			Revocation:  rev,
			Now:         s.Clock.Now().UTC().Format(time.RFC3339),
			Policy:      VectorPolicy{RSeconds: s.Policy.R().Seconds(), SkewSeconds: s.Policy.SkewTolerance().Seconds()},
			Expect:      VectorExpect{Decision: decisionToken(verdict.Decision), Causes: causes},
		})
	}
	return VectorFile{
		Schema:      SchemaVersion,
		Description: "Atlas Verification Core (M3) conformance vectors. Any conformant verifier must, given each vector's record + trust + revocation + now + policy, produce the expected decision and causes. See tests/vectors/VECTORS.md.",
		Vectors:     vectors,
	}, nil
}
