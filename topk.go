package search_trie

type topKHeapItem struct {
	key  string
	freq uint
}

type topKHeap struct {
	items []topKHeapItem
	limit int
}

func (h *topKHeap) Len() int {
	return len(h.items)
}

func (h *topKHeap) Less(i, j int) bool {
	return h.items[i].freq < h.items[j].freq
}

func (h *topKHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *topKHeap) Push(x interface{}) {
	h.items = append(h.items, x.(topKHeapItem))
}

func (h *topKHeap) Pop() interface{} {
	old := h.items
	item := old[0]
	h.items = old[1:]
	return item
}
