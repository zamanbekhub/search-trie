package search_trie

import (
	"container/heap"
)

type nodeInfo struct {
	Key       string
	Frequency uint
}

type node struct {
	frequency uint
	isEnd     bool
	children  map[string]*node
	topK      *topKHeap
}

func newnode(topK int) *node {
	heapInstance := &topKHeap{limit: topK}
	heap.Init(heapInstance) // Инициализация кучи
	return &node{
		children: map[string]*node{},
		topK:     heapInstance,
	}
}

func (root *node) getTopK(key string) []topKHeapItem {
	curr := root
	prefix := ""
	for _, r := range key {
		prefix += string(r)
		curr = curr.children[prefix]
		if curr == nil {
			return nil
		}
	}
	return curr.topK.items
}

func (root *node) put(key string, frequency uint) {
	curr, path := root.getOrCreate(key)
	curr.isEnd = true
	curr.frequency = frequency

	for _, n := range path {
		n.updateTopK(key, frequency)
	}
}

func (root *node) inc(key string) {
	curr, path := root.getOrCreate(key)
	curr.frequency += 1

	for _, n := range path {
		n.updateTopK(key, curr.frequency)
	}
}

func (root *node) getOrCreate(key string) (*node, []*node) {
	curr, path := root, make([]*node, 0, len(key))
	prefix := ""

	path = append(path, curr)
	for _, r := range key {
		prefix += string(r)
		child := curr.children[prefix]
		if child == nil {
			if curr.children == nil {
				curr.children = map[string]*node{}
			}
			child = newnode(curr.topK.limit)
			curr.children[prefix] = child
		}

		curr = curr.children[prefix]
		path = append(path, curr)
	}

	return curr, path
}

func (root *node) updateTopK(key string, freq uint) {
	for i, item := range root.topK.items {
		if item.key == key {
			// Update existing key
			root.topK.items[i].freq = freq
			heap.Fix(root.topK, i) // Reorder the heap
			return
		}
	}

	// Add new key
	heap.Push(root.topK, topKHeapItem{key: key, freq: freq})
	if root.topK.Len() > root.topK.limit {
		root.topK.Pop()
	}
}

func (root *node) traverse() <-chan nodeInfo {
	out := make(chan nodeInfo, 100)
	go func() {
		defer close(out)
		root.traverseHelper("", out)
	}()

	return out
}

func (root *node) traverseHelper(prefix string, out chan<- nodeInfo) {
	if prefix != "" && root.isEnd {
		out <- nodeInfo{
			Key:       prefix,
			Frequency: root.frequency,
		}
	}

	for r, child := range root.children {
		child.traverseHelper(r, out)
	}
}
