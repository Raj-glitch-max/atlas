package issuance

import "sort"

// PermissionSet is a principal's own permissions, against which a delegation's
// requested scope is checked for the strict-subset property at issuance
// (FR2, ER2, INV2, SO6). It is a set: order and duplicates do not matter.
type PermissionSet struct {
	perms map[string]struct{}
}

// NewPermissionSet builds a permission set from permission strings, ignoring
// duplicates and empty entries.
func NewPermissionSet(perms ...string) PermissionSet {
	m := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		if p != "" {
			m[p] = struct{}{}
		}
	}
	return PermissionSet{perms: m}
}

// Contains reports whether the set holds a permission.
func (s PermissionSet) Contains(p string) bool {
	_, ok := s.perms[p]
	return ok
}

// Len returns the number of distinct permissions.
func (s PermissionSet) Len() int { return len(s.perms) }

// list returns the permissions sorted, for trace summaries.
func (s PermissionSet) list() []string {
	out := make([]string, 0, len(s.perms))
	for p := range s.perms {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// isProperSupersetOf reports whether this set strictly contains the scope:
// every scope permission is held AND the scope is strictly smaller than the
// set (ER2 "strict subset" — an equal scope is refused). scope may contain
// duplicates; distinctness is computed here.
func (s PermissionSet) isProperSupersetOf(scope []string) bool {
	distinct := make(map[string]struct{}, len(scope))
	for _, p := range scope {
		if !s.Contains(p) {
			return false // scope reaches outside the principal's permissions
		}
		distinct[p] = struct{}{}
	}
	// Strict: the scope must be smaller than the full permission set.
	return len(distinct) < s.Len()
}

// checkpoint: feat(issuance): add boundary check (#60)
