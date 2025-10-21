package proxy

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/elazarl/goproxy"
)

// Ghost Labs Reaper CA certificate (exp. 2030-Oct-20)
var caCert = []byte(`-----BEGIN CERTIFICATE-----
MIIENTCCAx2gAwIBAgIUDwdggwuOcbMyK6ftH/tM8Y3gsQgwDQYJKoZIhvcNAQEL
BQAwgakxCzAJBgNVBAYTAlVTMQ4wDAYDVQQIDAVUZXhhczEPMA0GA1UEBwwGQXVz
dGluMR0wGwYDVQQKDBRHaG9zdCBTZWN1cml0eSwgSW5jLjETMBEGA1UECwwKR2hv
c3QgTGFiczEhMB8GA1UEAwwYcmVhcGVyLmdob3N0c2VjdXJpdHkuY29tMSIwIAYJ
KoZIhvcNAQkBFhNsYWJzQGdob3N0LnNlY3VyaXR5MB4XDTI1MTAyMTExMzYwMVoX
DTMwMTAyMDExMzYwMVowgakxCzAJBgNVBAYTAlVTMQ4wDAYDVQQIDAVUZXhhczEP
MA0GA1UEBwwGQXVzdGluMR0wGwYDVQQKDBRHaG9zdCBTZWN1cml0eSwgSW5jLjET
MBEGA1UECwwKR2hvc3QgTGFiczEhMB8GA1UEAwwYcmVhcGVyLmdob3N0c2VjdXJp
dHkuY29tMSIwIAYJKoZIhvcNAQkBFhNsYWJzQGdob3N0LnNlY3VyaXR5MIIBIjAN
BgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvXgOUfacFCgqBfEqyaBef75t5wQf
9Czur5z82iTIlbxw8+YgspyQE9+9v/NCnJ36lYirzT31OJ5fFbllSM9f3SfOYuZN
D4uBUfG9YHt7xU6KIyu2ClxvFJGQZtU+1BA4AmF/cNurj41SBQrDIjxb5gQVq9ZF
PgbviKXureVn8dpGzsICnqQPj1rng1p0FfibjzexLB0ye6UJLSnk3JGjXPgYQUgN
9lgaZ+fASdX3d/veWlh/mcufrKX+WN/X7rcCQlkSVhlO5GP0a2uIC4Rj//j1jUgT
k/mQbvFgN6/ixCV9LhCrtpmjLQZPOaZwc4aJjus6b/GUfrQgeDJjSIWb0wIDAQAB
o1MwUTAdBgNVHQ4EFgQUL48lnRsl0Zfeg+fwHd8cg6RXPoIwHwYDVR0jBBgwFoAU
L48lnRsl0Zfeg+fwHd8cg6RXPoIwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0B
AQsFAAOCAQEAn5zAPqldj148HoqYLwDDLDpkyt2I0473boLSj41iuc6nIokuzrQI
GIyq/lS3isvV7xisMMwLb+ON3/iqGhH7/QVJ69nM5NSvrtLLfy+tqsQjPRx/Gzo4
91WlqwBd90Cj7SoYgGFdr9Tmx/YKGr+EgQfILYwepARRWqbSo50pYJ4Z039v8TDK
EiUQT9pu5aFmQXW8cCizMd4/HH8S/XSuSMabzdoYf5AkOg4tWzUf+MWZtb/FpW4h
8HYkgMSyYIW3/aT8Pv9Tlj5hsDbS753Y818tLH6VqsMni+YiHyS5ZqXnvT+wRGy6
NEMKaF7lmehYmo7bRNOFNtFD3qj5ecN/8g==
-----END CERTIFICATE-----`)

// Ghost Labs Reaper CA key
var caKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC9eA5R9pwUKCoF
8SrJoF5/vm3nBB/0LO6vnPzaJMiVvHDz5iCynJAT372/80KcnfqViKvNPfU4nl8V
uWVIz1/dJ85i5k0Pi4FR8b1ge3vFToojK7YKXG8UkZBm1T7UEDgCYX9w26uPjVIF
CsMiPFvmBBWr1kU+Bu+Ipe6t5Wfx2kbOwgKepA+PWueDWnQV+JuPN7EsHTJ7pQkt
KeTckaNc+BhBSA32WBpn58BJ1fd3+95aWH+Zy5+spf5Y39futwJCWRJWGU7kY/Rr
a4gLhGP/+PWNSBOT+ZBu8WA3r+LEJX0uEKu2maMtBk85pnBzhomO6zpv8ZR+tCB4
MmNIhZvTAgMBAAECggEAClDOiso1U5L06Bo/fPclHf8QEcYI4qz4Q98gFH/F3K0u
rdYnjUH2muERfjFwPNcMbSF3T8yMnOFhMO1v+nJb+E5cr3LK2vVP/FZVPBBiZpWG
R4WaOt8vJxIJKCglFxEDqPww9EPph9PRwDeGn0cProAmGT3TIEjquBMSDIK90Y47
qvlO2Vn7VrFTYOqU+eM9VrLOd8Wy2FPUUfO10GVRvN0ZErkndVjevenqyWdmH3qM
DaPk+wIqGqcsk/1aP/ty0FM0/u9ctIJ4GSkyEhn/i6ZPpGcKQ+edcSgKa1OotLMr
Q7KoPgvG51Wpcbr6VWaobe8fSYegQtorY6hl9kB8iQKBgQDcrCf5zuKM9fQvTPSO
GqCGDDyrkUdHP+CC2wr/3Ipb3k/WySb6H9UUysO1uJ+xaRVQA8HzqaieUO3nEJMk
M5tmoJo2ioTagYXMyq+7DkfUmVyzka4RHqOqDiI2kMMoO5XuUwMz+Y57Xz6VxVwP
uJ429jnqxNLVCqaL2yDLUPaLawKBgQDbzRUvtcjBOHYUeMoWNX2E2mp9lOSI20T2
Zynyu1WgtYrbR9QCf9TjnzjSNkXvqVUmQgZXZ0KBZrlalHZkmHDFxk9zXAOK4zth
X6nJtyXwFErs2ndL8RG9noL5uTz2MzFMhkdoEkafyku6aBlxXB7NNwTzW16FTh3Z
dC1Tnl7zOQKBgQCPS6el2wdYW7qWIJXJ1VaZ1UZsbqlnhf5HWu/4mACsiV807WhH
EfavSr/tuBbTAJbbX6VJkckyDQF/g07ZOj3WVcHuWuLMdUEqbA/TGwHf9zqwTJBJ
A6lpm0XyQuzHqnHA0d0JmitAx+d/ICqY9tyeeiO/5NG3j/P4a3IPNOL0QwKBgQCB
e2PCslTsNmWhE7MAuEwUClL3XdHvKTSL7yQQAPmlbay6Fqs3ObTgzng7pYs3bsph
ej2gGY1dC3WffZvtELxGVdeR/p97nvbpGuC7mq+3qUymEOB8FSw9RvajQ9M8udWN
3gCMt09xbEuGKTLry7e9bm71KVsaLnV5F25oNwB6SQKBgAgRQKRhNciBMKMOhVBN
NDGZVdWcsi4Fwk+wpEGYQwd9BJofaJBzAwyPyRYCfFVw1P+baAJBKMNKGrRA+bmT
Ot6KuVlU/CmVgwuo4Y0StrO2Gk1BRaK7nr6t22xwYnPJ0/w/UxpADhrLEdx1Y+Te
vmBKExBNiMt1S065MqQF/LMb
-----END PRIVATE KEY-----`)

func initializeCA() error {
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return err
	}
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}

	return nil
}
