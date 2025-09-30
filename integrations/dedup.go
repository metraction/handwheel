package integrations

import (
	"sync"

	"github.com/metraction/handwheel/metrics"
)

// Dedup holds state for filtering duplicate images in-memory.
type Dedup struct {
	seen map[string]struct{}
	mu   sync.Mutex
}

// NewDedup constructs a Dedup instance.
func NewDedup() *Dedup {
	return &Dedup{
		seen: make(map[string]struct{}),
	}
}

// FilterDublicates is a predicate for flow.NewFilter to filter out already seen images.
func (d *Dedup) FilterDublicates(item any) bool {
	image, ok := item.(string)
	if !ok {
		// If not a string, let it pass (or change to false to drop)
		return true
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, exists := d.seen[image]; exists {
		return false
	}
	d.seen[image] = struct{}{}
	
	// Update unique images gauge
	metrics.UniqueImagesGauge.Set(float64(len(d.seen)))
	
	return true
}

