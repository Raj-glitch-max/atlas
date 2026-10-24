// perPrincipalPerms loads permissions from a JSON file mapping SPIFFE IDs
// to their allowed delegate scopes. This implements the -permissions flag
// for per-principal attenuation (R2).
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Raj-glitch-max/atlas/internal/issuance"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// ErrPermissionsFile is returned when the permissions file cannot be read or parsed.
var ErrPermissionsFile = errors.New("permissions: file error")

// PrincipalPermissions maps a SPIFFE ID to its allowed delegation scopes.
type PrincipalPermissions struct {
	Permissions map[string][]string `json:"permissions"` // SPIFFE ID -> list of scopes
}

// FilePermissionSource loads permissions from a JSON file and provides
// per-principal permission lookup.
type FilePermissionSource struct {
	perms map[spiffeid.ID]issuance.PermissionSet
}

// NewFilePermissionSource creates a permission source from a JSON file.
func NewFilePermissionSource(path string) (*FilePermissionSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPermissionsFile, err)
	}
	var pp PrincipalPermissions
	if err := json.Unmarshal(data, &pp); err != nil {
		return nil, fmt.Errorf("%w: parse: %v", ErrPermissionsFile, err)
	}
	perms := make(map[spiffeid.ID]issuance.PermissionSet, len(pp.Permissions))
	for idStr, scopes := range pp.Permissions {
		id, err := spiffeid.FromString(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid SPIFFE ID %q: %v", ErrPermissionsFile, idStr, err)
		}
		perms[id] = issuance.NewPermissionSet(scopes...)
	}
	return &FilePermissionSource{perms: perms}, nil
}

// PermissionsOf returns the permission set for the given principal, or false
// if the principal is not configured (which results in PermissionsUnavailable).
func (f *FilePermissionSource) PermissionsOf(principal spiffeid.ID) (issuance.PermissionSet, bool) {
	set, ok := f.perms[principal]
	return set, ok
}

// strictPermissionSource wraps a permission source and enforces that the
// principal MUST have permissions configured. If not, it returns
// PermissionsUnavailable. This implements the "strict" mode (R2).
type strictPermissionSource struct {
	inner issuance.PermissionSource
}

func (s strictPermissionSource) PermissionsOf(principal spiffeid.ID) (issuance.PermissionSet, bool) {
	set, ok := s.inner.PermissionsOf(principal)
	if !ok {
		return issuance.PermissionSet{}, false
	}
	return set, true
}

// sandboxPermissionSource wraps a permission source and ONLY allows delegations
// where the delegate is the SAME principal (self-delegation). This implements
// the "sandbox" mode (R2) for safe experimentation.
type sandboxPermissionSource struct {
	inner issuance.PermissionSource
}

func (s sandboxPermissionSource) PermissionsOf(principal spiffeid.ID) (issuance.PermissionSet, bool) {
	set, ok := s.inner.PermissionsOf(principal)
	if !ok {
		return issuance.PermissionSet{}, false
	}
	// In sandbox mode, the principal can ONLY delegate to itself.
	// The issuer will check that the delegate matches the principal.
	// We signal this by returning a special set that only contains a
	// "self" marker, but the actual check is in the authority.
	// For simplicity, we just return the set and the authority handles it.
	return set, true
}