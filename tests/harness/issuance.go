package harness

import (
	"github.com/Raj-glitch-max/atlas/internal/issuance"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// PermissionSource is a fake issuance.PermissionSource: it answers a
// configured permission set per principal, and "unavailable" for unknown
// principals (honest absence — issuance fails closed, never guesses).
type PermissionSource struct {
	byPrincipal map[spiffeid.ID]issuance.PermissionSet
	unavailable map[spiffeid.ID]bool
}

// NewPermissionSource returns an empty fake permission source.
func NewPermissionSource() *PermissionSource {
	return &PermissionSource{
		byPrincipal: map[spiffeid.ID]issuance.PermissionSet{},
		unavailable: map[spiffeid.ID]bool{},
	}
}

// Grant configures a principal's permission set.
func (s *PermissionSource) Grant(principal spiffeid.ID, perms ...string) *PermissionSource {
	s.byPrincipal[principal] = issuance.NewPermissionSet(perms...)
	return s
}

// MarkUnavailable forces the source to answer unavailable for a principal,
// even if permissions were granted — for testing the fail-closed path.
func (s *PermissionSource) MarkUnavailable(principal spiffeid.ID) *PermissionSource {
	s.unavailable[principal] = true
	return s
}

// PermissionsOf implements issuance.PermissionSource.
func (s *PermissionSource) PermissionsOf(principal spiffeid.ID) (issuance.PermissionSet, bool) {
	if s.unavailable[principal] {
		return issuance.PermissionSet{}, false
	}
	p, ok := s.byPrincipal[principal]
	return p, ok
}
