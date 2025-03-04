package certificate

import (
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	pb "github.com/ttpreport/ligolo-mp/v2/protobuf"
)

type Certificate struct {
	Name        string
	caName      string
	Certificate []byte
	Key         []byte
	Thumbprint  [sha1.Size]byte
}

func (cert *Certificate) KeyPair() (tls.Certificate, error) {
	return tls.X509KeyPair(cert.Certificate, cert.Key)
}

func (cert *Certificate) CertPool() (*x509.CertPool, error) {
	certpool := x509.NewCertPool()
	if ok := certpool.AppendCertsFromPEM(cert.Certificate); !ok {
		return nil, errors.New("error parsing certificate")
	}

	return certpool, nil
}

func (cert *Certificate) ExpiryDate() time.Time {
	keypair, err := cert.KeyPair()
	if err != nil {
		return time.Time{}
	}

	return keypair.Leaf.NotAfter
}

func (cert *Certificate) String() string {
	return fmt.Sprintf("Name=%s CA=%s", cert.Name, cert.caName)
}

func (cert *Certificate) Proto() *pb.Cert {
	return &pb.Cert{
		Name:        cert.Name,
		ExpiryDate:  cert.ExpiryDate().String(),
		Certificate: cert.Certificate,
		Key:         cert.Key,
	}
}

func ProtoToCertificate(p *pb.Cert) *Certificate {
	return &Certificate{
		Name:        p.Name,
		Certificate: p.Certificate,
		Key:         p.Key,
	}
}
