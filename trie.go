package search_trie

import (
	"sort"
	"sync"
)

type Trie struct {
	mu   sync.RWMutex
	root *node
}

// NewTrie creates a new Trie with the given topK limit.
func NewTrie(topK int) *Trie {
	return &Trie{root: newnode(topK)}
}

// TopK returns the top K most frequent words for prefix.
func (t *Trie) TopK(key string) []nodeInfo {
	if key == "" {
		return nil
	}

	t.mu.RLock() // Блокируем чтение
	defer t.mu.RUnlock()

	topK := t.root.getTopK(key)
	out := make([]nodeInfo, len(topK))
	for i, item := range topK {
		out[i] = nodeInfo{
			Key:       item.key,
			Frequency: item.freq,
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Frequency > out[j].Frequency
	})

	return out
}

// Has checks trie has the key.
func (t *Trie) Has(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.root.has(key)
}

// Put inserts the given key/frequency pair into the Trie.
func (t *Trie) Put(key string, frequency uint) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.root.put(key, frequency)
}

// Inc increments the frequency of the given key.
func (t *Trie) Inc(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.root.inc(key)
}

// Traverse returns all keys in the Trie.
func (t *Trie) Traverse() <-chan nodeInfo {
	return t.root.traverse()
}
