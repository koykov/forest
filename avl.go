package forest

import "cmp"

type avl[K cmp.Ordered, V any] struct {
	root  *nodeAVL[K, V]
	size  int
	nullK K
	nullV V
}

type nodeAVL[K cmp.Ordered, V any] struct {
	key   K
	value V
	left  *nodeAVL[K, V]
	right *nodeAVL[K, V]
}

func (t *avl[K, V]) Insert(key K, value V) {
	// todo implement me
}

func (t *avl[K, V]) Delete(key K) bool {
	// todo implement me
	return false
}

func (t *avl[K, V]) Search(key K) (V, bool) {
	// todo implement me
	return t.nullV, false
}

func (t *avl[K, V]) Min() (K, V, bool) {
	// todo implement me
	return t.nullK, t.nullV, false
}

func (t *avl[K, V]) Max() (K, V, bool) {
	// todo implement me
	return t.nullK, t.nullV, false
}

func (t *avl[K, V]) Size() int {
	// todo implement me
	return 0
}

func (t *avl[K, V]) Clear() {
	// todo implement me
}
