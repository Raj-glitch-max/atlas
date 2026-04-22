package record

// RevBinding is the opaque, optional revocation-binding element of a record
// (AD-015): mechanism-specific per-record data — a status reference, an
// accumulator witness, or nothing — that rides the record so the spike-
// selected revocation composition never mutates the stable surface (AP12).
//
// Interpretation rights (INTERFACE_SPECIFICATION.md §1): Revocation Status
// Provider realizations only. The record model carries it, the issuance
// authority obtains it from its RevBindingSource port (empty until the
// mechanism exists), and the verifier ignores it.
//
// A nil or empty RevBinding means absent.
type RevBinding []byte

// IsAbsent reports whether no revocation binding is carried.
func (b RevBinding) IsAbsent() bool { return len(b) == 0 }

// clone returns an independent copy, normalizing empty to nil (absent).
func (b RevBinding) clone() RevBinding {
	if len(b) == 0 {
		return nil
	}
	c := make(RevBinding, len(b))
	copy(c, b)
	return c
}

// checkpoint: chore(sdk): harden CLI flag configuration
