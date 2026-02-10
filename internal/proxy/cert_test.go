package proxy

import "testing"

func TestGenerateCA(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatalf("generating CA: %v", err)
	}

	if ca.Cert == nil {
		t.Fatal("expected non-nil cert")
	}
	if ca.Key == nil {
		t.Fatal("expected non-nil key")
	}
	if !ca.Cert.IsCA {
		t.Error("cert should be a CA")
	}
	if ca.Cert.Subject.CommonName != "Ghost Security Reaper CA - https://ghostsecurity.ai" {
		t.Errorf("CN = %q, want Ghost Security Reaper CA - https://ghostsecurity.ai", ca.Cert.Subject.CommonName)
	}
}
