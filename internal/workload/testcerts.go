package workload

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/url"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
)

// TestCertConfig holds configuration for generating test certificates.
type TestCertConfig struct {
	// TrustDomain is the SPIFFE trust domain.
	TrustDomain string
	// SPIFFEID is the SPIFFE ID for the certificate (e.g., "spiffe://example.org/workload").
	SPIFFEID string
	// DNSNames are the DNS names to include in the certificate.
	DNSNames []string
	// IPAddresses are the IP addresses to include in the certificate.
	IPAddresses []net.IP
	// NotBefore is the certificate start time.
	NotBefore time.Time
	// NotAfter is the certificate expiration time.
	NotAfter time.Time
	// IsCA indicates if this certificate should be a CA.
	IsCA bool
	// Parent is the parent certificate for signing (nil for self-signed root).
	Parent *TestCertConfig
	// ParentKey is the parent's private key for signing.
	ParentKey *ecdsa.PrivateKey
}

// GenerateTestCert generates a test X.509 certificate with SPIFFE ID.
// It returns the certificate, private key, and PEM-encoded versions.
func GenerateTestCert(config TestCertConfig) (*x509.Certificate, *ecdsa.PrivateKey, []byte, []byte, error) {
	// Generate private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Parse SPIFFE ID
	id, err := spiffeid.FromString(config.SPIFFEID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Build template
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1).SetBytes(randomBytes(16)),
		Subject: pkix.Name{
			CommonName: id.String(),
		},
		NotBefore:             config.NotBefore,
		NotAfter:              config.NotAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  config.IsCA,
		DNSNames:              config.DNSNames,
		IPAddresses:           config.IPAddresses,
	}

	// Add SPIFFE ID as URN SAN
	template.URIs = []*url.URL{id.URL()}

	// Determine signing key and parent
	var parent *x509.Certificate
	signKey := privateKey
	if config.Parent != nil && config.ParentKey != nil {
		parent = &x509.Certificate{
			Subject: pkix.Name{CommonName: config.Parent.SPIFFEID},
			NotBefore: config.Parent.NotBefore,
			NotAfter:  config.Parent.NotAfter,
			IsCA:      config.Parent.IsCA,
		}
		signKey = config.ParentKey
	}

	// Self-sign if no parent
	if parent == nil {
		parent = template
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, &privateKey.PublicKey, signKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Parse the created certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Encode to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})

	return cert, privateKey, certPEM, keyPEM, nil
}

// GenerateTestCA generates a test root CA certificate.
func GenerateTestCA(trustDomain string, notBefore, notAfter time.Time) (*x509.Certificate, *ecdsa.PrivateKey, []byte, []byte, error) {
	config := TestCertConfig{
		TrustDomain: trustDomain,
		SPIFFEID:    "spiffe://" + trustDomain + "/spiffe/ca",
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		IsCA:        true,
	}
	return GenerateTestCert(config)
}

// GenerateTestWorkloadSVID generates a test workload SVID signed by a CA.
func GenerateTestWorkloadSVID(trustDomain, workloadName string, caCert *x509.Certificate, caKey *ecdsa.PrivateKey, notBefore, notAfter time.Time) (*x509.Certificate, *ecdsa.PrivateKey, []byte, []byte, error) {
	config := TestCertConfig{
		TrustDomain: trustDomain,
		SPIFFEID:    "spiffe://" + trustDomain + "/workload/" + workloadName,
		DNSNames:    []string{"localhost"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		IsCA:        false,
		Parent: &TestCertConfig{
			SPIFFEID: "spiffe://" + trustDomain + "/spiffe/ca",
			NotBefore: notBefore,
			NotAfter:  notAfter,
			IsCA:      true,
		},
		ParentKey: caKey,
	}
	return GenerateTestCert(config)
}

// GenerateTestBundle creates a test trust bundle from a CA certificate.
func GenerateTestBundle(caCert *x509.Certificate) *Bundle {
	trustDomain := spiffeid.RequireTrustDomainFromString(caCert.Subject.CommonName)
	return &Bundle{
		TrustDomain:     trustDomain,
		X509Authorities: []*x509.Certificate{caCert},
	}
}

// GenerateTestX509SVID generates a complete X509SVID for testing.
func GenerateTestX509SVID(trustDomain, workloadName string, caCert *x509.Certificate, caKey *ecdsa.PrivateKey, notBefore, notAfter time.Time) (*X509SVID, error) {
	cert, key, _, _, err := GenerateTestWorkloadSVID(trustDomain, workloadName, caCert, caKey, notBefore, notAfter)
	if err != nil {
		return nil, err
	}

	id, err := spiffeid.FromString("spiffe://" + trustDomain + "/workload/" + workloadName)
	if err != nil {
		return nil, err
	}

	return &X509SVID{
		ID:           id,
		Certificates: []*x509.Certificate{cert, caCert},
		PrivateKey:   key,
		Hint:         "test-svid-" + workloadName,
	}, nil
}

// randomBytes generates random bytes.
func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

// TestX509SVIDSource creates a StaticSource with a test SVID and bundle.
func TestX509SVIDSource(trustDomain, workloadName string, notBefore, notAfter time.Time) (*StaticSource, *x509.Certificate, *ecdsa.PrivateKey, error) {
	caCert, caKey, _, _, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	if err != nil {
		return nil, nil, nil, err
	}

	svid, err := GenerateTestX509SVID(trustDomain, workloadName, caCert, caKey, notBefore, notAfter)
	if err != nil {
		return nil, nil, nil, err
	}

	bundle := &Bundle{
		TrustDomain:     spiffeid.RequireTrustDomainFromString(trustDomain),
		X509Authorities: []*x509.Certificate{caCert},
	}

	source := NewStaticSource(svid, map[spiffeid.TrustDomain]*Bundle{
		bundle.TrustDomain: bundle,
	})

	return source, caCert, caKey, nil
}