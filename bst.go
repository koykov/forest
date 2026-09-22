package forest

import "cmp"

type bst[K cmp.Ordered, V any] struct {
	options
	buf []nodeBST[K, V]
	off int

	root  *nodeBST[K, V]
	size  int
	nullK K
	nullV V
}

type nodeBST[K cmp.Ordered, V any] struct {
	key   K
	value V
	left  *nodeBST[K, V]
	right *nodeBST[K, V]
}

func NewBST[K cmp.Ordered, V any](options ...Option) Binary[K, V] {
	t := &bst[K, V]{}
	t.apply(options...)
	t.buf = make([]nodeBST[K, V], t.options.size)
	return t
}

func (t *bst[K, V]) Insert(key K, value V) {
	if t.root == nil {
		t.root = t.alloc(key, value)
		t.size++
		return
	}

	node := t.root
	for {
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
			node = node.left
		default:
			if node.right == nil {
				node.right = t.alloc(key, value)
				t.size++
				return
			}
			node = node.right
		}
	}
}

func (t *bst[K, V]) Delete(key K) bool {
	var parent *nodeBST[K, V]
	node := t.root

	for node != nil && node.key != key {
		parent = node
		if key < node.key {
			node = node.left
		} else {
			node = node.right
		}
	}
	if node == nil {
		return false
	}

	switch {
	case node.left == nil && node.right == nil:
		t.replace(parent, node, nil)
	case node.left == nil && node.right != nil:
		t.replace(parent, node, node.right)
	case node.left != nil && node.right == nil:
		t.replace(parent, node, node.left)
	default:
		t.delete2C(node)
	}
	t.size--
	return true
}

// replace swaps child into parent's slot (or into root when parent is nil).
func (t *bst[K, V]) replace(parent, child, with *nodeBST[K, V]) {
	if parent == nil {
		t.root = with
		return
	}
	if parent.left == child {
		parent.left = with
	} else {
		parent.right = with
	}
}

// delete2C replaces node's key/value with its in-order successor's key/value, then removes the successor node,
// that has at most one child (its right subtree).
func (t *bst[K, V]) delete2C(node *nodeBST[K, V]) {
	var parent *nodeBST[K, V]
	succ := node.right
	for succ.left != nil {
		parent = succ
		succ = succ.left
	}

	node.key, node.value = succ.key, succ.value

	if parent == nil {
		// Successor is node.right itself.
		node.right = succ.right
	} else {
		parent.left = succ.right
	}
}

func (t *bst[K, V]) Search(key K) (V, bool) {
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

func (t *bst[K, V]) Min() (K, V, bool) {
	if t.root == nil {
		return t.nullK, t.nullV, false
	}
	node := t.root
	for node.left != nil {
		node = node.left
	}
	return node.key, node.value, true
}

func (t *bst[K, V]) Max() (K, V, bool) {
	if t.root == nil {
		return t.nullK, t.nullV, false
	}
	node := t.root
	for node.right != nil {
		node = node.right
	}
	return node.key, node.value, true
}

func (t *bst[K, V]) Size() int {
	return t.size
}

func (t *bst[K, V]) Clear() {
	t.root = nil
	t.size = 0
	t.off = 0
}

func (t *bst[K, V]) alloc(key K, value V) (n *nodeBST[K, V]) {
	if t.off < len(t.buf) {
		n = &t.buf[t.off]
	} else {
		t.buf = append(t.buf, nodeBST[K, V]{})
		n = &t.buf[len(t.buf)-1]
	}
	t.off++
	n.key = key
	n.value = value
	n.left = nil
	n.right = nil
	return
}
