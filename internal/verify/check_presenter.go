package verify

import (
	"github.com/Raj-glitch-max/atlas/internal/record"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// checkPresenter is the presenter-binding stage (INV1 extension): when the
// policy requires a presenter, the verifier must be supplied with a presenter
// SPIFFE ID (from the authenticated peer, e.g., mTLS), and it must match the
// record's delegate. If the presenter is absent or mismatched, verification
// fails definitively.
func checkPresenter(a record.Assertions, presenter spiffeid.ID, requirePresenter bool) CheckEntry {
	dig := digest("presenter", a.Delegate.String(), presenter.String())

	if !requirePresenter {
		return pass(CheckPresenterBinding, "presenter binding not required", dig)
	}
	if presenter.IsZero() {
		return failDefinitive(CheckPresenterBinding, PresenterAbsent,
			"presenter required but not supplied", dig)
	}
	if a.Delegate != presenter {
		return failDefinitive(CheckPresenterBinding, PresenterMismatch,
			"presenter does not match record delegate", dig)
	}
	return pass(CheckPresenterBinding, "presenter matches delegate", dig)
}
