package record

import (
	"fmt"
	"io"
	"os"
)

// MaxBundleBytes bounds a trust bundle read from disk or the network.
//
// A bundle carries the trust material plus a signed revocation snapshot, so its
// size scales with the number of unexpired revocations — legitimately large, but
// not unbounded. 32 MiB accommodates roughly a million revoked instance ids
// (~16 bytes each, base64 and JSON overhead included), which is already beyond
// the point where snapshot signing saturates a 1 Hz republish cycle.
//
// The bound exists because the bundle is read BEFORE its signature can be
// checked — you cannot verify what you have not yet parsed — so an attacker who
// can place bytes on the read path chooses the allocation size. That is the same
// shape as the unauthenticated request-body exhaustion fixed in the server, one
// layer down, and it applies to every relying party: the CLI's offline verify
// and the reference gate both read bundles from disk on demand.
const MaxBundleBytes int64 = 32 << 20

// ReadBundleFile reads a bundle with an enforced size ceiling.
//
// It reads MaxBundleBytes+1 and fails if the extra byte materialises, so an
// oversized file is refused rather than silently truncated — a truncated bundle
// would fail signature verification anyway, but with a misleading error that
// sends an operator hunting for corruption instead of for size.
func ReadBundleFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	raw, err := io.ReadAll(io.LimitReader(f, MaxBundleBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(raw)) > MaxBundleBytes {
		return nil, fmt.Errorf("trust bundle %s exceeds the %d-byte limit — refusing to parse", path, MaxBundleBytes)
	}
	return raw, nil
}
