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
	if t.root == nil {
		t.root = t.alloc(key, value)
		t.size++
		return
	}
	t.insert(t.root, key, value)
}

func (t *avl[K, V]) insert(node *nodeAVL[K, V], key K, value V) {
	switch {
	case key == node.key:
		node.value = value
		return
	case key < node.key:
		if node.left == nil {
			node.left = t.alloc(key, value)
			t.size++
			return
		}
		t.insert(node.left, key, value)
	default:
		if node.right == nil {
			node.right = t.alloc(key, value)
			t.size++
			return
		}
		t.insert(node.right, key, value)
	}

	// Update height.
	node.height = max(t.hOf(node.left), t.hOf(node.right)) + 1

	// Check balance factor.
	bf := t.bfOf(node)
	switch {
	case bf > -2 && bf < 2:
		// Subtree is balanced, do nothing.
		return
	case bf > 1 && key < node.left.key:
		// todo implement LL
	case bf > 1 && key > node.left.key:
		// todo implement LR
	case bf < -1 && key > node.right.key:
		// todo implement RR
	case bf < -1 && key < node.right.key:
		// todo implement RL
	}
}

func (t *avl[K, V]) bfOf(node *nodeAVL[K, V]) int {
	if node == nil {
		return 0
	}
	var hl, hr int
	if node.left != nil {
		hl = node.left.height
	}
	if node.right != nil {
		hr = node.right.height
	}
	return hl - hr
}

func (t *avl[K, V]) hOf(node *nodeAVL[K, V]) int {
	if node == nil {
		return 0
	}
	return node.height
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

func (t *avl[K, V]) alloc(key K, value V) (n *nodeAVL[K, V]) {
	if t.off < len(t.buf) {
		n = &t.buf[t.off]
	} else {
		t.buf = append(t.buf, nodeAVL[K, V]{})
		n = &t.buf[len(t.buf)-1]
	}
	t.off++
	n.height = 0
	n.key = key
	n.value = value
	n.left = nil
	n.right = nil
	return
}
