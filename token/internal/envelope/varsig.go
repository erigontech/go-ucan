package envelope

import (
	"sync"

	"github.com/ucan-wg/go-varsig"
)

var (
	varsigsMu sync.RWMutex
	varsigs   = varsig.DefaultRegistry()
)

// RegisterVarsig adds a signature algorithm that envelopes may be sealed with.
func RegisterVarsig(alg varsig.Algorithm, decode varsig.DecodeFunc) {
	varsigsMu.Lock()
	defer varsigsMu.Unlock()
	varsigs.Register(alg, decode)
}

func decodeVarsig(data []byte) (varsig.Varsig, error) {
	varsigsMu.RLock()
	defer varsigsMu.RUnlock()
	return varsigs.Decode(data)
}
