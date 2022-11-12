// Package cache package implement the basic block to implement the
// cache persistence
package cache

// NodeInfoCache implement the interface to
// store the node information inside the cache.
type NodeInfoCache struct {
	ID       string
	Alias    string
	Color    string
	Features *string
}
