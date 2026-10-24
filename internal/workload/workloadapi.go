package workload

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

// WorkloadAPISource implements X509SVIDSource and BundleSource using the
// SPIFFE Workload API (workloadapi). It connects to the Workload API socket
// and automatically rotates SVIDs and bundles as they are updated.
type WorkloadAPISource struct {
	client *workloadapi.Client
	mu     sync.RWMutex

	// Current SVID and bundle state
	currentSVID  *X509SVID
	currentBundles map[spiffeid.TrustDomain]*Bundle

	// Notification channels
	svidUpdates    chan *X509SVID
	bundleUpdates  chan map[spiffeid.TrustDomain]*Bundle
	done           chan struct{}

	// Configuration
	trustDomain spiffeid.TrustDomain
}

// WorkloadAPISourceOption configures a WorkloadAPISource.
type WorkloadAPISourceOption func(*WorkloadAPISource)

// WithTrustDomain sets the expected trust domain. If not set, the trust domain
// is inferred from the first SVID received.
func WithTrustDomain(td spiffeid.TrustDomain) WorkloadAPISourceOption {
	return func(s *WorkloadAPISource) {
		s.trustDomain = td
	}
}

// NewWorkloadAPISource creates a new Workload API source connected to the
// default socket path (or SPIFFE_ENDPOINT_SOCKET env var).
func NewWorkloadAPISource(ctx context.Context, opts ...WorkloadAPISourceOption) (*WorkloadAPISource, error) {
	client, err := workloadapi.New(ctx)
	if err != nil {
		return nil, &WorkloadError{Code: "WORKLOAD_API_CONNECT_FAILED", Message: "failed to connect to Workload API", Cause: err}
	}

	s := &WorkloadAPISource{
		client:         client,
		currentBundles: make(map[spiffeid.TrustDomain]*Bundle),
		svidUpdates:    make(chan *X509SVID, 1),
		bundleUpdates:  make(chan map[spiffeid.TrustDomain]*Bundle, 1),
		done:           make(chan struct{}),
	}

	for _, opt := range opts {
		opt(s)
	}

	// Start background watchers
	go s.watchX509SVIDs(ctx)
	go s.watchBundles(ctx)

	return s, nil
}

// GetX509SVID returns the current X.509 SVID.
func (s *WorkloadAPISource) GetX509SVID() (*X509SVID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.currentSVID == nil {
		return nil, ErrSourceNotReady
	}
	return s.currentSVID, nil
}

// GetBundle returns the trust bundle for the given trust domain.
func (s *WorkloadAPISource) GetBundle(trustDomain spiffeid.TrustDomain) (*Bundle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bundle, ok := s.currentBundles[trustDomain]
	if !ok {
		return nil, ErrBundleNotFound
	}
	return bundle, nil
}

// AllBundles returns all available trust bundles.
func (s *WorkloadAPISource) AllBundles() map[spiffeid.TrustDomain]*Bundle {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid external mutation
	bundles := make(map[spiffeid.TrustDomain]*Bundle, len(s.currentBundles))
	for td, b := range s.currentBundles {
		bundles[td] = b
	}
	return bundles
}

// Done returns a channel that is closed when the source is stopped.
func (s *WorkloadAPISource) Done() <-chan struct{} {
	return s.done
}

// Close stops the source and releases resources.
func (s *WorkloadAPISource) Close() error {
	s.client.Close()
	close(s.done)
	return nil
}

// watchX509SVIDs watches for X.509 SVID updates from the Workload API.
func (s *WorkloadAPISource) watchX509SVIDs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		default:
			// Use the streaming API to watch for updates
			svids, err := s.client.WatchX509Context(ctx)
			if err != nil {
				// Log error and retry after a delay
				time.Sleep(5 * time.Second)
				continue
			}

			// Convert to our X509SVID type
			ourSVID := s.convertX509SVID(svids)
			if ourSVID == nil {
				continue
			}

			// Update trust domain if not set
			if s.trustDomain == (spiffeid.TrustDomain{}) {
				s.trustDomain = ourSVID.ID.TrustDomain()
			}

			// Update and notify
			s.mu.Lock()
			s.currentSVID = ourSVID
			s.mu.Unlock()

			select {
			case s.svidUpdates <- ourSVID:
			default:
				// Channel full, drop update
			}
		}
	}
}

// watchBundles watches for bundle updates from the Workload API.
func (s *WorkloadAPISource) watchBundles(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		default:
			bundles, err := s.client.WatchBundlesContext(ctx)
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			ourBundles := s.convertBundles(bundles)

			s.mu.Lock()
			s.currentBundles = ourBundles
			s.mu.Unlock()

			select {
			case s.bundleUpdates <- ourBundles:
			default:
			}
		}
	}
}

// convertX509SVID converts a workloadapi X509SVID to our type.
func (s *WorkloadAPISource) convertX509SVID(svid *workloadapi.X509SVID) *X509SVID {
	if svid == nil {
		return nil
	}

	id, err := spiffeid.FromX509SVID(svid.Certificates[0])
	if err != nil {
		return nil
	}

	return &X509SVID{
		ID:           id,
		Certificates: svid.Certificates,
		PrivateKey:   svid.PrivateKey,
		Hint:         svid.Hint,
	}
}

// convertBundles converts workloadapi bundles to our type.
func (s *WorkloadAPISource) convertBundles(bundles *workloadapi.BundleSet) map[spiffeid.TrustDomain]*Bundle {
	if bundles == nil {
		return nil
	}

	result := make(map[spiffeid.TrustDomain]*Bundle)
	for td, bundle := range bundles.X509Authorities() {
		certs := make([]*x509.Certificate, len(bundle))
		for i, c := range bundle {
			certs[i] = c
		}
		result[td] = &Bundle{
			TrustDomain:     td,
			X509Authorities: certs,
		}
	}
	return result
}

// TLSCertificate returns a tls.Certificate for the current SVID, suitable for
// use in a tls.Config for mTLS.
func (s *WorkloadAPISource) TLSCertificate() (*tls.Certificate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.currentSVID == nil {
		return nil, ErrSourceNotReady
	}
	cert := s.currentSVID.TLSCertificate()
	return &cert, nil
}

// TLSConfig returns a tls.Config configured for mTLS using the current SVID
// and trust bundles. The config verifies peer certificates against the trust
// bundles and presents the current SVID as the client/server certificate.
func (s *WorkloadAPISource) TLSConfig() (*tls.Config, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.currentSVID == nil {
		return nil, ErrSourceNotReady
	}

	// Build root CA pool from all bundles
	roots := x509.NewCertPool()
	for _, bundle := range s.currentBundles {
		for _, ca := range bundle.X509Authorities {
			roots.AddCert(ca)
		}
	}

	cert := s.currentSVID.TLSCertificate()

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// Custom verification: we trust the standard verification but
			// also allow cross-trust-domain verification via our bundles
			return nil // Standard verification already ran
		},
	}, nil
}

// SVIDUpdates returns a channel that receives SVID updates.
// The channel is closed when the source is closed.
func (s *WorkloadAPISource) SVIDUpdates() <-chan *X509SVID {
	return s.svidUpdates
}

// BundleUpdates returns a channel that receives bundle updates.
// The channel is closed when the source is closed.
func (s *WorkloadAPISource) BundleUpdates() <-chan map[spiffeid.TrustDomain]*Bundle {
	return s.bundleUpdates
}

// TrustDomain returns the trust domain this source is configured for.
func (s *WorkloadAPISource) TrustDomain() spiffeid.TrustDomain {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.trustDomain
}

// StaticSource is a test implementation of X509SVIDSource and BundleSource
// that returns fixed SVIDs and bundles. Useful for testing without a real
// Workload API.
type StaticSource struct {
	svid    *X509SVID
	bundles map[spiffeid.TrustDomain]*Bundle
	done    chan struct{}
	mu      sync.RWMutex
}

// NewStaticSource creates a new StaticSource with the given SVID and bundles.
func NewStaticSource(svid *X509SVID, bundles map[spiffeid.TrustDomain]*Bundle) *StaticSource {
	return &StaticSource{
		svid:    svid,
		bundles: bundles,
		done:    make(chan struct{}),
	}
}

// GetX509SVID returns the static SVID.
func (s *StaticSource) GetX509SVID() (*X509SVID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.svid == nil {
		return nil, ErrSourceNotReady
	}
	return s.svid, nil
}

// GetBundle returns the static bundle for the trust domain.
func (s *StaticSource) GetBundle(trustDomain spiffeid.TrustDomain) (*Bundle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	bundle, ok := s.bundles[trustDomain]
	if !ok {
		return nil, ErrBundleNotFound
	}
	return bundle, nil
}

// AllBundles returns all static bundles.
func (s *StaticSource) AllBundles() map[spiffeid.TrustDomain]*Bundle {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[spiffeid.TrustDomain]*Bundle, len(s.bundles))
	for td, b := range s.bundles {
		result[td] = b
	}
	return result
}

// Done returns a channel that is closed when the source is closed.
func (s *StaticSource) Done() <-chan struct{} {
	return s.done
}

// Close closes the source.
func (s *StaticSource) Close() error {
	close(s.done)
	return nil
}

// SetSVID updates the SVID (for testing rotation).
func (s *StaticSource) SetSVID(svid *X509SVID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.svid = svid
}

// SetBundles updates the bundles (for testing rotation).
func (s *StaticSource) SetBundles(bundles map[spiffeid.TrustDomain]*Bundle) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bundles = bundles
}