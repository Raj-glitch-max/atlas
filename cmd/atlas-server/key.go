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
			return nil, fmt.Errorf("write key file %s: %w", path, err)
		}
		return k, nil
	default:
		return nil, err
	}
}
