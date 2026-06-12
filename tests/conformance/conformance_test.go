package conformance_test

import (
	"testing"

	"github.com/Raj-glitch-max/atlas/tests/conformance"
)

// The V1 verifier must pass the full conformance corpus. When a competing
// implementation exists, it runs this same Run with its own factory — the
// oracle is shared, so verifier differentials (the Frankencerts failure)
// surface as a corpus failure rather than an exploitable production gap.
func TestV1Conformance(t *testing.T) {
	conformance.Run(t, conformance.V1Factory)
}
