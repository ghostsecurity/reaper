package settings

import "sync"

type Provider struct {
	mu       sync.RWMutex
	settings Settings
}

func newProvider(settings Settings) *Provider {
	return &Provider{
		settings: settings,
	}
}

func (p *Provider) Get() Settings {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.settings
}

func (p *Provider) Modify(modifier func(s *Settings)) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := make([]byte, len(p.settings.CAKey))
	copy(key, p.settings.CAKey)

	cert := make([]byte, len(p.settings.CACert))
	copy(cert, p.settings.CACert)

	p.settings.CAKey = string(key)
	p.settings.CACert = string(cert)

	modifier(&p.settings)
	return Save(&p.settings)
}
