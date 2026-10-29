package vectors

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"testing"
	"time"

	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/Raj-glitch-max/atlas/internal/revstatus"
	"github.com/Raj-glitch-max/atlas/internal/truststore"
	"github.com/Raj-glitch-max/atlas/internal/verify"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/stretchr/testify/require"
)

// PresenterVector represents a test vector for presenter binding verification.
type PresenterVector struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Policy      PolicyConfig      `json:"policy"`
	Record      RecordConfig      `json:"record"`
	Presenter   string            `json:"presenter"`
	Expect      string            `json:"expect"`
	Causes      []string          `json:"causes"`
	Trace       TraceConfig       `json:"trace"`
}

type PolicyConfig struct {
	R              string `json:"R"`
	SkewTolerance  string `json:"skewTolerance"`
	RequirePresenter bool   `json:"requirePresenter"`
}

type RecordConfig struct {
	Principal    string   `json:"principal"`
	Delegate     string   `json:"delegate"`
	Scope        []string `json:"scope"`
	IssuedAt     string   `json:"issuedAt"`
	Expiration   string   `json:"expiration"`
	IssuerKeyId  string   `json:"issuerKeyId"`
}

type TraceConfig struct {
	Checks []CheckConfig `json:"checks"`
}

type CheckConfig struct {
	Name     string `json:"name"`
	Outcome  string `json:"outcome"`
	Cause    string `json:"cause,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

type PresenterVectorsFile struct {
	Vectors []PresenterVector `json:"vectors"`
}

func loadPresenterVectors(t *testing.T) PresenterVectorsFile {
	data := readTestFile(t, "presenter-vectors.json")
	var vectors PresenterVectorsFile
	err := json.Unmarshal(data, &vectors)
	require.NoError(t, err)
	return vectors
}

func TestPresenterVectors(t *testing.T) {
	vectors := loadPresenterVectors(t)

	// Create a test key for signing
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	trustDomain := spiffeid.RequireTrustDomainFromString("example.org")
	tm, err := record.NewTrustMaterial(trustDomain, map[string]*ecdsa.PublicKey{"authority-key-1": &key.PublicKey})
	require.NoError(t, err)

	trustStore := truststore.New()
	require.NoError(t, trustStore.Provision(tm, time.Now()))

	// Create a revocation publisher and provider
	listID := "atlas-revlist:example.org"
	publisher, err := revstatus.NewPublisher(key, listID)
	require.NoError(t, err)
	provider := revstatus.NewSignedSetProvider(&key.PublicKey, listID)

	// Publish an empty revocation set
	snap, err := publisher.Publish([]record.InstanceID{}, time.Now())
	require.NoError(t, err)
	_, err = provider.Ingest(snap)
	require.NoError(t, err)

	for _, v := range vectors.Vectors {
		t.Run(v.Name, func(t *testing.T) {
			// Parse policy
			R, err := time.ParseDuration(v.Policy.R)
			require.NoError(t, err)
			skew, err := time.ParseDuration(v.Policy.SkewTolerance)
			require.NoError(t, err)
			policy, err := verify.NewPolicy(R, skew, v.Policy.RequirePresenter)
			require.NoError(t, err)

			// Parse record config
			principal, err := spiffeid.FromString(v.Record.Principal)
			require.NoError(t, err)
			delegate, err := spiffeid.FromString(v.Record.Delegate)
			require.NoError(t, err)
			issuedAt, err := time.Parse(time.RFC3339, v.Record.IssuedAt)
			require.NoError(t, err)
			expiration, err := time.Parse(time.RFC3339, v.Record.Expiration)
			require.NoError(t, err)

			// Create the record
			rec, err := record.NewRecord(record.Input{
				Principal:    principal,
				Delegate:     delegate,
				Scope:        v.Record.Scope,
				IssuedAt:     issuedAt,
				Expiration:   expiration,
				IssuerKeyID:  v.Record.IssuerKeyId,
				Signer: record.Signer{
					Key:    key,
					KeyID:  v.Record.IssuerKeyId,
				},
			})
			require.NoError(t, err)

			// Parse presenter if provided
			var presenter spiffeid.ID
			if v.Presenter != "" {
				presenter, err = spiffeid.FromString(v.Presenter)
				require.NoError(t, err)
			}

			// Verify
			verifier, err := verify.NewVerifier(policy, trustStore, provider, realClock{})
			require.NoError(t, err)

			verdict, trace := verifier.Verify(rec.Presented(), presenter)

			// Check decision
			require.Equal(t, v.Expect, verdict.Decision.String(), "decision mismatch")

			// Check causes
			if len(v.Causes) > 0 {
				require.ElementsMatch(t, v.Causes, causeStrings(verdict.Causes), "causes mismatch")
			} else {
				require.Empty(t, verdict.Causes, "expected no causes")
			}

			// Check trace checks
			require.Len(t, trace.Entries, len(v.Trace.Checks))
			for i, expectedCheck := range v.Trace.Checks {
				entry := trace.Entries[i]
				require.Equal(t, expectedCheck.Name, entry.Check.String())
				require.Equal(t, expectedCheck.Outcome, entry.Outcome.String())
				if expectedCheck.Cause != "" {
					require.Equal(t, expectedCheck.Cause, entry.Cause.String())
				}
				if expectedCheck.Detail != "" {
					require.Contains(t, entry.Detail, expectedCheck.Detail)
				}
			}
		})
	}
}

// realClock implements verify.Clock for tests.
type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// causeStrings converts a slice of Causes to strings.
func causeStrings(causes []verify.Cause) []string {
	result := make([]string, len(causes))
	for i, c := range causes {
		result[i] = c.String()
	}
	return result
}