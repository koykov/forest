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
	t.root = t.insert(t.root, key, value)
	t.root.height = max(t.hOf(t.root.left), t.hOf(t.root.right)) + 1
}

func (t *avl[K, V]) insert(node *nodeAVL[K, V], key K, value V) *nodeAVL[K, V] {
	switch {
	case key == node.key:
		node.value = value
		return node
	case key < node.key:
		if node.left == nil {
			node.left = t.alloc(key, value)
			node.height = max(t.hOf(node.left), t.hOf(node.right)) + 1
			t.size++
			return node
		}
		node.left = t.insert(node.left, key, value)
	default:
		if node.right == nil {
			node.right = t.alloc(key, value)
			node.height = max(t.hOf(node.left), t.hOf(node.right)) + 1
			t.size++
			return node
		}
		node.right = t.insert(node.right, key, value)
	}

	// Update height.
	node.height = max(t.hOf(node.left), t.hOf(node.right)) + 1

	// Check balance factor.
	bf := t.bfOf(node)
	switch {
	case bf > -2 && bf < 2:
		// Subtree is balanced, do nothing.
		return node
	case bf > 1 && key < node.left.key: // LL
		return t.rotr(node)
	case bf > 1 && key > node.left.key: // LR
		node.left = t.rotl(node.left)
		return t.rotr(node)
	case bf < -1 && key > node.right.key: // RR
		return t.rotl(node)
	case bf < -1 && key < node.right.key: // RL
		node.right = t.rotr(node.right)
		return t.rotl(node)
	}
	return node
}

func (t *avl[K, V]) rotl(x *nodeAVL[K, V]) *nodeAVL[K, V] {
	y := x.right

	x.right = y.left
	y.left = x

	x.height = max(t.hOf(x.left), t.hOf(x.right)) + 1
	y.height = max(t.hOf(y.left), t.hOf(y.right)) + 1

	return y
}

func (t *avl[K, V]) rotr(x *nodeAVL[K, V]) *nodeAVL[K, V] {
	y := x.left

	x.left = y.right
	y.right = x

	x.height = max(t.hOf(x.left), t.hOf(x.right)) + 1
	y.height = max(t.hOf(y.left), t.hOf(y.right)) + 1

	return y
}

func (t *avl[K, V]) bfOf(node *nodeAVL[K, V]) int {
	if node == nil {
		return 0
	}
	return t.hOf(node.left) - t.hOf(node.right)
}

func (t *avl[K, V]) hOf(node *nodeAVL[K, V]) int {
	if node == nil {
		return -1
	}
	return node.height
}

func (t *avl[K, V]) Delete(key K) (ok bool) {
	if t.root == nil {
		return
	}
	t.root, ok = t.delete(t.root, key)
	return
}

func (t *avl[K, V]) delete(node *nodeAVL[K, V], key K) (*nodeAVL[K, V], bool) {
	if node == nil {
		return nil, false
	}
	var ok bool
	switch {
	case key < node.key:
		node.left, ok = t.delete(node.left, key)
	case key > node.key:
		node.right, ok = t.delete(node.right, key)
	default:
		ok = true
		switch {
		case node.left == nil:
			return node.right, ok
		case node.right == nil:
			return node.left, ok
		default:
			succ := node.right
			for succ.left != nil {
				succ = succ.left
			}
			node.key, node.value = succ.key, succ.value
			node.right, _ = t.delete(node.right, succ.key)
		}
	}
	if !ok {
		return node, false
	}

	// Update height.
	node.height = max(t.hOf(node.left), t.hOf(node.right)) + 1

	// Check balance factor.
	bf := t.bfOf(node)
	switch {
	case bf > -2 && bf < 2:
		// Subtree is balanced, do nothing.
		return node, ok
	case bf > 1 && key < node.left.key: // LL
		return t.rotr(node), ok
	case bf > 1 && key > node.left.key: // LR
		node.left = t.rotl(node.left)
		return t.rotr(node), ok
	case bf < -1 && key > node.right.key: // RR
		return t.rotl(node), ok
	case bf < -1 && key < node.right.key: // RL
		node.right = t.rotr(node.right)
		return t.rotl(node), ok
	}

	return nil, false
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

func (t *avl[K, V]) height() int {
	if t.root != nil {
		return t.root.height
	}
	return 0
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
