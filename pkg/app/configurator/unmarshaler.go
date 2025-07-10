package configurator

import (
	"encoding/json"
	"sync"

	"gopkg.in/yaml.v3"
	"github.com/BurntSushi/toml"
)

var registry = &unmarshalerRegistry{
	unmarshalers: map[string]LoaderFn{
		"json": json.Unmarshal,
		"toml": toml.Unmarshal,
		"yaml": yaml.Unmarshal,
	},
}

type (
	// LoaderFn is the function type for loading configuration.
	LoaderFn func([]byte, any) error

	// unmarshalerRegistry is the registry for unmarshalers.
	unmarshalerRegistry struct {
		unmarshalers map[string]LoaderFn
		mu           sync.RWMutex
	}
)

// RegisterUnmarshaler registers an unmarshaler.
func RegisterUnmarshaler(name string, fn LoaderFn) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.unmarshalers[name] = fn
}

// Unmarshaler returns the unmarshaler by name.
func Unmarshaler(name string) (LoaderFn, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	fn, ok := registry.unmarshalers[name]
	return fn, ok
}