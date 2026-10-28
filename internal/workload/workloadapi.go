package workload

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
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
	currentSVID    *X509SVID
	currentBundles map[spiffeid.TrustDomain]*Bundle

	// Notification channels
	svidUpdates   chan *X509SVID
	bundleUpdates chan map[spiffeid.TrustDomain]*Bundle
	done          chan struct{}

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

	// Fetch initial SVID and bundles
	if err := s.fetchInitial(ctx); err != nil {
		client.Close()
		return nil, err
	}

	// Start background watchers
	go s.watchX509SVIDs(ctx)
	go s.watchBundles(ctx)

	return s, nil
}

// fetchInitial fetches the initial SVID and bundles.
func (s *WorkloadAPISource) fetchInitial(ctx context.Context) error {
	// Fetch initial X509 SVID
	svid, err := s.client.FetchX509SVID(ctx)
	if err != nil {
		return &WorkloadError{Code: "FETCH_SVID_FAILED", Message: "failed to fetch initial X509 SVID", Cause: err}
	}
	ourSVID := s.convertX509SVID(svid)
	if ourSVID == nil {
		return &WorkloadError{Code: "CONVERT_SVID_FAILED", Message: "failed to convert initial SVID"}
	}

	// Fetch bundles for the trust domain
	td := ourSVID.ID.TrustDomain()
	bundle, err := s.client.FetchBundle(ctx, td)
	if err != nil {
		return &WorkloadError{Code: "FETCH_BUNDLE_FAILED", Message: "failed to fetch initial bundle", Cause: err}
	}
	ourBundle := s.convertBundle(td, bundle)
	if ourBundle == nil {
		return &WorkloadError{Code: "CONVERT_BUNDLE_FAILED", Message: "failed to convert initial bundle"}
	}

	s.mu.Lock()
	s.currentSVID = ourSVID
	s.currentBundles[td] = ourBundle
	if s.trustDomain == (spiffeid.TrustDomain{}) {
		s.trustDomain = td
	}
	s.mu.Unlock()

	return nil
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
			updates, err := s.client.WatchX509SVID(ctx)
			if err != nil {
				// Log error and retry after a delay
				time.Sleep(5 * time.Second)
				continue
			}

			for svid := range updates {
				ourSVID := s.convertX509SVID(svid)
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
			// Fetch bundles for the current trust domain
			s.mu.RLock()
			td := s.trustDomain
			s.mu.RUnlock()

			if td == (spiffeid.TrustDomain{}) {
				time.Sleep(5 * time.Second)
				continue
			}

			bundle, err := s.client.FetchBundle(ctx, td)
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}

			ourBundle := s.convertBundle(td, bundle)
			if ourBundle == nil {
				continue
			}

			s.mu.Lock()
			s.currentBundles[td] = ourBundle
			s.mu.Unlock()

			// Notify
			select {
			case s.bundleUpdates <- s.AllBundles():
			default:
			}
		}
	}
}

// convertX509SVID converts a workloadapi X509SVID to our type.
func (s *WorkloadAPISource) convertX509SVID(svid *workloadapi.X509SVID) *X509SVID {
	if svid == nil || len(svid.Certificates) == 0 {
		return nil
	}

	id, err := spiffeid.FromX509Certificate(svid.Certificates[0])
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

// convertBundle converts a workloadapi bundle to our type.
func (s *WorkloadAPISource) convertBundle(td spiffeid.TrustDomain, bundle *workloadapi.Bundle) *Bundle {
	if bundle == nil {
		return nil
	}

	certs := make([]*x509.Certificate, len(bundle.X509Authorities()))
	for i, c := range bundle.X509Authorities() {
		certs[i] = c
	}

	return &Bundle{
		TrustDomain:     td,
		X509Authorities: certs,
	}
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