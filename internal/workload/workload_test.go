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
	"testing"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/stretchr/testify/require"
)

func TestGenerateTestCA(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	caCert, caKey, caCertPEM, caKeyPEM, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	require.NoError(t, err)
	require.NotNil(t, caCert)
	require.NotNil(t, caKey)
	require.NotEmpty(t, caCertPEM)
	require.NotEmpty(t, caKeyPEM)

	// Verify it's a CA
	require.True(t, caCert.IsCA)
	require.True(t, caCert.BasicConstraintsValid)

	// Verify SPIFFE ID in certificate
	id, err := spiffeid.FromX509SVID(caCert)
	require.NoError(t, err)
	require.Equal(t, "spiffe://"+trustDomain+"/spiffe/ca", id.String())

	// Verify PEM can be parsed
	block, _ := pem.Decode(caCertPEM)
	require.NotNil(t, block)
	require.Equal(t, "CERTIFICATE", block.Type)

	block, _ = pem.Decode(caKeyPEM)
	require.NotNil(t, block)
	require.Equal(t, "PRIVATE KEY", block.Type)
}

func TestGenerateTestWorkloadSVID(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	caCert, caKey, _, _, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	require.NoError(t, err)

	workloadCert, workloadKey, workloadCertPEM, workloadKeyPEM, err := GenerateTestWorkloadSVID(trustDomain, "test-workload", caCert, caKey, notBefore, notAfter)
	require.NoError(t, err)
	require.NotNil(t, workloadCert)
	require.NotNil(t, workloadKey)
	require.NotEmpty(t, workloadCertPEM)
	require.NotEmpty(t, workloadKeyPEM)

	// Verify it's not a CA
	require.False(t, workloadCert.IsCA)

	// Verify SPIFFE ID in certificate
	id, err := spiffeid.FromX509SVID(workloadCert)
	require.NoError(t, err)
	require.Equal(t, "spiffe://"+trustDomain+"/workload/test-workload", id.String())

	// Verify DNS names
	require.Contains(t, workloadCert.DNSNames, "localhost")

	// Verify IP addresses
	require.Contains(t, workloadCert.IPAddresses, net.ParseIP("127.0.0.1"))
}

func TestGenerateTestX509SVID(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	caCert, caKey, _, _, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	require.NoError(t, err)

	svid, err := GenerateTestX509SVID(trustDomain, "test-workload", caCert, caKey, notBefore, notAfter)
	require.NoError(t, err)
	require.NotNil(t, svid)
	require.Equal(t, "spiffe://"+trustDomain+"/workload/test-workload", svid.ID.String())
	require.Len(t, svid.Certificates, 2) // leaf + CA
	require.NotNil(t, svid.PrivateKey)
	require.Equal(t, "test-svid-test-workload", svid.Hint)
}

func TestVerifyX509SVID(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	caCert, caKey, _, _, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	require.NoError(t, err)

	svid, err := GenerateTestX509SVID(trustDomain, "test-workload", caCert, caKey, notBefore, notAfter)
	require.NoError(t, err)

	bundle := &Bundle{
		TrustDomain:     spiffeid.RequireTrustDomainFromString(trustDomain),
		X509Authorities: []*x509.Certificate{caCert},
	}

	verifiedID, err := VerifyX509SVID(svid, bundle, VerifyOptions{
		TrustDomain: spiffeid.RequireTrustDomainFromString(trustDomain),
		Time:        now,
	})
	require.NoError(t, err)
	require.Equal(t, svid.ID, verifiedID)
}

func TestVerifyX509SVID_WrongTrustDomain(t *testing.T) {
	trustDomain := "example.org"
	otherDomain := "other.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	caCert, caKey, _, _, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	require.NoError(t, err)

	svid, err := GenerateTestX509SVID(trustDomain, "test-workload", caCert, caKey, notBefore, notAfter)
	require.NoError(t, err)

	bundle := &Bundle{
		TrustDomain:     spiffeid.RequireTrustDomainFromString(otherDomain),
		X509Authorities: []*x509.Certificate{caCert},
	}

	_, err = VerifyX509SVID(svid, bundle, VerifyOptions{
		TrustDomain: spiffeid.RequireTrustDomainFromString(otherDomain),
		Time:        now,
	})
	require.Error(t, err)
	require.True(t, IsErrTrustDomainMismatch(err))
}

func TestVerifyX509SVID_Expired(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-2 * time.Hour)
	notAfter := now.Add(-time.Hour) // Expired

	caCert, caKey, _, _, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	require.NoError(t, err)

	svid, err := GenerateTestX509SVID(trustDomain, "test-workload", caCert, caKey, notBefore, notAfter)
	require.NoError(t, err)

	bundle := &Bundle{
		TrustDomain:     spiffeid.RequireTrustDomainFromString(trustDomain),
		X509Authorities: []*x509.Certificate{caCert},
	}

	_, err = VerifyX509SVID(svid, bundle, VerifyOptions{
		TrustDomain: spiffeid.RequireTrustDomainFromString(trustDomain),
		Time:        now,
	})
	require.Error(t, err)
	require.True(t, IsErrVerificationFailed(err))
}

func TestStaticSource(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	source, _, _, err := TestX509SVIDSource(trustDomain, "test-workload", notBefore, notAfter)
	require.NoError(t, err)
	require.NotNil(t, source)

	// Get SVID
	svid, err := source.GetX509SVID()
	require.NoError(t, err)
	require.NotNil(t, svid)
	require.Equal(t, "spiffe://"+trustDomain+"/workload/test-workload", svid.ID.String())

	// Get bundle
	bundle, err := source.GetBundle(spiffeid.RequireTrustDomainFromString(trustDomain))
	require.NoError(t, err)
	require.NotNil(t, bundle)
	require.Len(t, bundle.X509Authorities, 1)

	// All bundles
	bundles := source.AllBundles()
	require.Len(t, bundles, 1)

	// Close
	err = source.Close()
	require.NoError(t, err)

	// After close, Done should be closed
	select {
	case <-source.Done():
		// Expected
	default:
		t.Fatal("Done channel not closed after Close")
	}
}

func TestStaticSource_SetSVID(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	source, caCert, caKey, err := TestX509SVIDSource(trustDomain, "workload1", notBefore, notAfter)
	require.NoError(t, err)

	// Create new SVID for workload2
	newSVID, err := GenerateTestX509SVID(trustDomain, "workload2", caCert, caKey, notBefore, notAfter)
	require.NoError(t, err)

	// Update source
	source.SetSVID(newSVID)

	// Verify new SVID is returned
	svid, err := source.GetX509SVID()
	require.NoError(t, err)
	require.Equal(t, "spiffe://"+trustDomain+"/workload/workload2", svid.ID.String())
}

func TestStaticSource_SetBundles(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	source, caCert, caKey, err := TestX509SVIDSource(trustDomain, "test-workload", notBefore, notAfter)
	require.NoError(t, err)

	// Create new CA and bundle
	newCaCert, newCaKey, _, _, err := GenerateTestCA("other.org", notBefore, notAfter)
	require.NoError(t, err)

	newBundle := &Bundle{
		TrustDomain:     spiffeid.RequireTrustDomainFromString("other.org"),
		X509Authorities: []*x509.Certificate{newCaCert},
	}

	source.SetBundles(map[spiffeid.TrustDomain]*Bundle{
		newBundle.TrustDomain: newBundle,
	})

	bundle, err := source.GetBundle(spiffeid.RequireTrustDomainFromString("other.org"))
	require.NoError(t, err)
	require.Equal(t, newCaCert, bundle.X509Authorities[0])
}

func TestX509SVID_TLSCertificate(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	caCert, caKey, _, _, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	require.NoError(t, err)

	svid, err := GenerateTestX509SVID(trustDomain, "test-workload", caCert, caKey, notBefore, notAfter)
	require.NoError(t, err)

	tlsCert := svid.TLSCertificate()
	require.Len(t, tlsCert.Certificate, 2)
	require.NotNil(t, tlsCert.PrivateKey)
}

func TestX509SVID_ExpiresAt(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	caCert, caKey, _, _, err := GenerateTestCA(trustDomain, notBefore, notAfter)
	require.NoError(t, err)

	svid, err := GenerateTestX509SVID(trustDomain, "test-workload", caCert, caKey, notBefore, notAfter)
	require.NoError(t, err)

	expiresAt := svid.ExpiresAt()
	require.Equal(t, notAfter, expiresAt)
}

func TestBundleSource_AllBundles(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()
	notBefore := now.Add(-time.Hour)
	notAfter := now.Add(time.Hour)

	source, caCert, caKey, err := TestX509SVIDSource(trustDomain, "test-workload", notBefore, notAfter)
	require.NoError(t, err)

	bundles := source.AllBundles()
	require.Len(t, bundles, 1)
	require.Contains(t, bundles, spiffeid.RequireTrustDomainFromString(trustDomain))
	require.Equal(t, caCert, bundles[spiffeid.RequireTrustDomainFromString(trustDomain)].X509Authorities[0])
}

// Helper to create a test certificate without the test helpers
func createTestCert(t *testing.T, isCA bool, spiffeID string, parent *x509.Certificate, parentKey *ecdsa.PrivateKey) (*x509.Certificate, *ecdsa.PrivateKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	id, err := spiffeid.FromString(spiffeID)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1).SetBytes(randomBytes(16)),
		Subject: pkix.Name{
			CommonName: id.String(),
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  isCA,
		URIs:                  []*url.URL{id.URL()},
	}

	signKey := privateKey
	if parent != nil && parentKey != nil {
		signKey = parentKey
	} else {
		parent = template
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, &privateKey.PublicKey, signKey)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)

	return cert, privateKey
}

func TestVerifyX509SVID_DNSNames(t *testing.T) {
	trustDomain := "example.org"
	now := time.Now()

	caCert, caKey := createTestCert(t, true, "spiffe://"+trustDomain+"/spiffe/ca", nil, nil)
	workloadCert, workloadKey := createTestCert(t, false, "spiffe://"+trustDomain+"/workload/test", caCert, caKey)

	svid := &X509SVID{
		ID:           spiffeid.RequireIDFromString("spiffe://"+trustDomain+"/workload/test"),
		Certificates: []*x509.Certificate{workloadCert, caCert},
		PrivateKey:   workloadKey,
	}

	bundle := &Bundle{
		TrustDomain:     spiffeid.RequireTrustDomainFromString(trustDomain),
		X509Authorities: []*x509.Certificate{caCert},
	}

	// Verify with matching DNS name
	_, err := VerifyX509SVID(svid, bundle, VerifyOptions{
		TrustDomain: spiffeid.RequireTrustDomainFromString(trustDomain),
		DNSNames:    []string{"test.example.org"},
		Time:        now,
	})
	// Should fail because cert doesn't have this DNS name
	require.Error(t, err)
}