package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// loadOrCreateKey returns the authority signing key. With an empty path it
// generates an ephemeral key (fine for demos/tests, but records won't survive
// a restart). With a path it loads the PEM key if present, otherwise generates
// one and writes it (0600) — so issued records stay verifiable across restarts.
func loadOrCreateKey(path string) (*ecdsa.PrivateKey, error) {
	if path == "" {
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	raw, err := os.ReadFile(path)
	switch {
	case err == nil:
		blk, _ := pem.Decode(raw)
		if blk == nil {
			return nil, fmt.Errorf("key file %s: not PEM-encoded", path)
		}
		k, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
		if err != nil {
			return nil, fmt.Errorf("key file %s: %w", path, err)
		}
		ek, ok := k.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key file %s: not an ECDSA key", path)
		}
		if ek.Curve != elliptic.P256() {
			return nil, fmt.Errorf("key file %s: curve must be P-256", path)
		}
		return ek, nil
	case os.IsNotExist(err):
		k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		der, err := x509.MarshalPKCS8PrivateKey(k)
		if err != nil {
			return nil, err
		}
		pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
		if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
			// Refusing to start is deliberate: -key asked for durability, and
			// silently falling back to an ephemeral key would mean every record
			// issued before a restart stops verifying, with nothing in the logs
			// to explain why. Fail loudly instead of pretending.
			//
			// The overwhelmingly common cause is a container that correctly runs
			// as a non-root user meeting a volume the platform mounted as root,
			// so the message says so rather than leaving an operator to guess.
			if os.IsPermission(err) {
				return nil, fmt.Errorf(
					"write key file %s: %w\n"+
						"  This image runs as uid 65532 (non-root), and the directory holding that\n"+
						"  path is not writable by it — typically a platform-mounted volume owned by\n"+
						"  root. Either:\n"+
						"    • give the volume to uid 65532 (Docker: `chown -R 65532:65532` the host\n"+
						"      dir; Kubernetes: fsGroup: 65532), or\n"+
						"    • run without -key/-store for an ephemeral in-memory instance, or\n"+
						"    • run the container as root (Railway: RAILWAY_RUN_UID=0), which trades\n"+
						"      away the non-root hardening this image is built for.",
					path, err)
			}
			return nil, fmt.Errorf("write key file %s: %w", path, err)
		}
		return k, nil
	default:
		return nil, err
	}
}
