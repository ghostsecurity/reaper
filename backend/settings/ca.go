package settings

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	certFilename = "ca.crt"
	pkFilename   = "ca.key"
)

func LoadCA() (cert []byte, pk []byte, err error) {

	configDir, err := getDir()
	if err != nil {
		return nil, nil, err
	}

	certPath := filepath.Join(configDir, certFilename)
	pkPath := filepath.Join(configDir, pkFilename)

	// return ca/key stored on disk if it already exists
	cert, err = os.ReadFile(certPath)
	if err == nil {
		pk, err = os.ReadFile(pkPath)
		if err == nil {
			return cert, pk, nil
		}
	}

	// otherwise generate a new ca/key pair and save it to disk for next time
	cert, pk, err = RegenerateCA()
	if err != nil {
		return nil, nil, err
	}

	if err := saveCA(cert, pk); err != nil {
		return nil, nil, err
	}

	return cert, pk, nil
}

func saveCA(cert []byte, pk []byte) error {

	configDir, err := getDir()
	if err != nil {
		return err
	}

	certPath := filepath.Join(configDir, certFilename)
	pkPath := filepath.Join(configDir, pkFilename)

	if err := os.WriteFile(certPath, cert, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(pkPath, pk, 0600); err != nil {
		return err
	}

	return nil
}

func RegenerateCA() (cert []byte, pk []byte, err error) {

	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization:  []string{"Ghost Security"},
			Country:       []string{"Ghost Security"},
			Province:      []string{"Ghost Security"},
			Locality:      []string{"Ghost Security"},
			StreetAddress: []string{"Ghost Security"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}); err != nil {
		return nil, nil, err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	}); err != nil {
		return nil, nil, err
	}

	caPemBytes := caPEM.Bytes()
	caPrivKeyPemBytes := caPrivKeyPEM.Bytes()

	return caPemBytes, caPrivKeyPemBytes, nil
}
