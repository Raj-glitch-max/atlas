package workload

import (
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/svid/x509svid"
)

// X509SVIDSource is the interface for obtaining X.509 SVIDs from the
// SPIFFE Workload API. Implementations must be safe for concurrent use.
// The source must rotate SVIDs automatically as they approach expiry.
type X509SVIDSource interface {
	// GetX509SVID returns the current X.509 SVID for the workload.
	// It returns an error if no SVID is available or the source is not ready.
	GetX509SVID() (*X509SVID, error)
	// Done returns a channel that is closed when the source is stopped
	// and no longer provides SVIDs.
	Done() <-chan struct{}
}

// X509SVID represents an X.509 SVID (SPIFFE Verifiable Identity Document)
// with its associated private key and trust bundle.
type X509SVID struct {
	// ID is the SPIFFE ID of the workload.
	ID spiffeid.ID
	// Certificates is the chain of X.509 certificates. The first element
	// is the leaf SVID; subsequent elements form the chain up to the trust
	// domain root.
	Certificates []*x509.Certificate
	// PrivateKey is the private key corresponding to the leaf certificate.
	PrivateKey interface{} // crypto.Signer
	// Hint is an optional hint for the server to identify this SVID.
	Hint string
}

// TLSCertificate returns a tls.Certificate suitable for use in tls.Config.
func (s *X509SVID) TLSCertificate() tls.Certificate {
	cert := tls.Certificate{
		Certificate: make([][]byte, len(s.Certificates)),
		PrivateKey:  s.PrivateKey,
	}
	for i, c := range s.Certificates {
		cert.Certificate[i] = c.Raw
	}
	return cert
}

// ExpiresAt returns the expiration time of the leaf certificate.
func (s *X509SVID) ExpiresAt() time.Time {
	if len(s.Certificates) > 0 {
		return s.Certificates[0].NotAfter
	}
	return time.Time{}
}

// BundleSource is the interface for obtaining trust bundles (root CAs)
// for one or more trust domains.
type BundleSource interface {
	// GetBundle returns the trust bundle for the given trust domain.
	// It returns an error if the bundle is not available.
	GetBundle(trustDomain spiffeid.TrustDomain) (*Bundle, error)
	// AllBundles returns all available trust bundles.
	AllBundles() map[spiffeid.TrustDomain]*Bundle
	// Done returns a channel that is closed when the source is stopped.
	Done() <-chan struct{}
}

// Bundle represents a trust bundle for a SPIFFE trust domain.
type Bundle struct {
	// TrustDomain is the SPIFFE trust domain this bundle belongs to.
	TrustDomain spiffeid.TrustDomain
	// X509Authorities is the set of root CA certificates for validating
	// X.509 SVIDs from this trust domain.
	X509Authorities []*x509.Certificate
}

// VerifyOptions contains options for verifying an X.509 SVID.
type VerifyOptions struct {
	// TrustDomain is the expected trust domain of the SVID.
	TrustDomain spiffeid.TrustDomain
	// DNSNames is the expected DNS names in the SVID (for server verification).
	DNSNames []string
	// Time is the time to verify against (defaults to time.Now()).
	Time time.Time
}

// VerifyX509SVID verifies an X.509 SVID against a trust bundle.
// It returns the verified SPIFFE ID on success.
func VerifyX509SVID(svid *X509SVID, bundle *Bundle, opts VerifyOptions) (spiffeid.ID, error) {
	// Use the first certificate as the leaf
	if len(svid.Certificates) == 0 {
		return spiffeid.ID{}, ErrNoCertificates
	}
	leaf := svid.Certificates[0]

	// Verify the SPIFFE ID in the certificate
	id, err := x509svid.IDFromCert(leaf)
	if err != nil {
		return spiffeid.ID{}, ErrInvalidSPIFFEID
	}

	// Check trust domain matches
	if id.TrustDomain() != opts.TrustDomain {
		return spiffeid.ID{}, ErrTrustDomainMismatch
	}

	// Verify certificate chain against bundle
	roots := x509.NewCertPool()
	for _, ca := range bundle.X509Authorities {
		roots.AddCert(ca)
	}

	verifyOpts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
		CurrentTime:   opts.Time,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	}

	// Add intermediate certificates
	for i := 1; i < len(svid.Certificates); i++ {
		verifyOpts.Intermediates.AddCert(svid.Certificates[i])
	}

	_, err = leaf.Verify(verifyOpts)
	if err != nil {
		return spiffeid.ID{}, ErrVerificationFailed
	}

	// Verify DNS names if provided (for server verification)
	for _, dns := range opts.DNSNames {
		if err := leaf.VerifyHostname(dns); err != nil {
			return spiffeid.ID{}, ErrDNSMismatch
		}
	}

	return id, nil
}

// Error types for workload operations.
var (
	ErrNoCertificates      = &WorkloadError{Code: "NO_CERTIFICATES", Message: "X509 SVID has no certificates"}
	ErrInvalidSPIFFEID     = &WorkloadError{Code: "INVALID_SPIFFE_ID", Message: "certificate does not contain a valid SPIFFE ID"}
	ErrTrustDomainMismatch = &WorkloadError{Code: "TRUST_DOMAIN_MISMATCH", Message: "SVID trust domain does not match expected"}
	ErrVerificationFailed  = &WorkloadError{Code: "VERIFICATION_FAILED", Message: "certificate chain verification failed"}
	ErrDNSMismatch         = &WorkloadError{Code: "DNS_MISMATCH", Message: "certificate does not match expected DNS name"}
	ErrSourceNotReady      = &WorkloadError{Code: "SOURCE_NOT_READY", Message: "workload API source not ready"}
	ErrSourceClosed        = &WorkloadError{Code: "SOURCE_CLOSED", Message: "workload API source has been closed"}
	ErrBundleNotFound      = &WorkloadError{Code: "BUNDLE_NOT_FOUND", Message: "trust bundle not found for trust domain"}
	ErrNoPrivateKey        = &WorkloadError{Code: "NO_PRIVATE_KEY", Message: "X509 SVID has no private key"}
)

// WorkloadError is a structured error for workload operations.
type WorkloadError struct {
	Code    string
	Message string
	Cause   error
}

func (e *WorkloadError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *WorkloadError) Unwrap() error {
	return e.Cause
}

func (e *WorkloadError) Is(target error) bool {
	if te, ok := target.(*WorkloadError); ok {
		return e.Code == te.Code
	}
	return false
}
