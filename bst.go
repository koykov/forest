package forest

import "cmp"

type bst[K cmp.Ordered, V any] struct {
	root  *nodeBST[K, V]
	nullK K
	nullV V
}

type nodeBST[K cmp.Ordered, V any] struct {
	key   K
	left  *nodeBST[K, V]
	right *nodeBST[K, V]
}

func NewBST[K cmp.Ordered, V any]() Binary[K, V] {
	return &bst[K, V]{}
}

func (t *bst[K, V]) Insert(key K, value V) {
	// todo implement me
}

func (t *bst[K, V]) Delete(key K) bool {
	// todo implement me
	return false
}

func (t *bst[K, V]) Search(key K) (V, bool) {
	// todo implement me
	return t.nullV, false
}

func (t *bst[K, V]) Min() (K, V, bool) {
	// todo implement me
	return t.nullK, t.nullV, false
}

func (t *bst[K, V]) Max() (K, V, bool) {
	// todo implement me
	return t.nullK, t.nullV, false
}

func (t *bst[K, V]) Size() int {
	// todo implement me
	return 0
}

func (t *bst[K, V]) Clear() {
	// todo implement me
}
