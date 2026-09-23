package forest

import "cmp"

type avl[K cmp.Ordered, V any] struct {
	options
	buf []nodeAVL[K, V]
	off int

	root  *nodeAVL[K, V]
	size  int
	nullK K
	nullV V
}

type nodeAVL[K cmp.Ordered, V any] struct {
	height int
	key    K
	value  V
	left   *nodeAVL[K, V]
	right  *nodeAVL[K, V]
}

func NewAVL[K cmp.Ordered, V any](options ...Option) Binary[K, V] {
	t := &avl[K, V]{}
	t.apply(options...)
	t.buf = make([]nodeAVL[K, V], t.options.size)
	return t
}

func (t *avl[K, V]) Insert(key K, value V) {
	// todo implement me
}

func (t *avl[K, V]) Delete(key K) bool {
	// todo implement me
	return false
}

func (t *avl[K, V]) Search(key K) (V, bool) {
	node := t.root
	for node != nil {
		switch {
		case key == node.key:
			return node.value, true
		case key < node.key:
			node = node.left
		default:
			node = node.right
		}
	}
	return t.nullV, false
}

func (t *avl[K, V]) Min() (K, V, bool) {
	if t.root == nil {
		return t.nullK, t.nullV, false
	}
	node := t.root
	for node.left != nil {
		node = node.left
	}
	return node.key, node.value, true
}

func (t *avl[K, V]) Max() (K, V, bool) {
	if t.root == nil {
		return t.nullK, t.nullV, false
	}
	node := t.root
	for node.right != nil {
		node = node.right
	}
	return node.key, node.value, true
}

func (t *avl[K, V]) Size() int {
	return t.size
}

func (t *avl[K, V]) Clear() {
	t.root = nil
	t.size = 0
	t.off = 0
}
