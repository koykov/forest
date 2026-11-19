package forest

import "cmp"

const (
	rbBlack uint8 = iota
	rbRed
)

type rb[K cmp.Ordered, V any] struct {
	options
	buf []nodeRB[K, V]
	off int

	root  *nodeRB[K, V]
	size  int
	nullK K
	nullV V
}

type nodeRB[K cmp.Ordered, V any] struct {
	key    K
	value  V
	color  uint8
	parent *nodeRB[K, V]
	left   *nodeRB[K, V]
	right  *nodeRB[K, V]
}

func NewRB[K cmp.Ordered, V any](options ...Option) Binary[K, V] {
	t := &rb[K, V]{}
	t.apply(options...)
	t.buf = make([]nodeRB[K, V], t.options.size)
	return t
}

func (t *rb[K, V]) Insert(key K, value V) {
	if t.root == nil {
		t.root = t.alloc(key, value)
		t.size++
		return
	}
	t.root = t.insert(t.root, key, value)
}

func (t *rb[K, V]) insert(node *nodeRB[K, V], key K, value V) *nodeRB[K, V] {
	// todo implement me
	return nil
}

func (t *rb[K, V]) Delete(key K) (ok bool) {
	if t.root == nil {
		return
	}
	if t.root, ok = t.delete(t.root, key); ok {
		t.size--
	}
	return
}

func (t *rb[K, V]) delete(node *nodeRB[K, V], key K) (*nodeRB[K, V], bool) {
	// todo implement me
	return nil, false
}

func (t *rb[K, V]) Search(key K) (V, bool) {
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

func (t *rb[K, V]) Min() (K, V, bool) {
	if t.root == nil {
		return t.nullK, t.nullV, false
	}
	node := t.root
	for node.left != nil {
		node = node.left
	}
	return node.key, node.value, true
}

func (t *rb[K, V]) Max() (K, V, bool) {
	if t.root == nil {
		return t.nullK, t.nullV, false
	}
	node := t.root
	for node.right != nil {
		node = node.right
	}
	return node.key, node.value, true
}

func (t *rb[K, V]) Size() int {
	return t.size
}

func (t *rb[K, V]) Clear() {
	t.root = nil
	t.size = 0
	t.off = 0
}

func (t *rb[K, V]) alloc(key K, value V) (n *nodeRB[K, V]) {
	if t.off < len(t.buf) {
		n = &t.buf[t.off]
	} else {
		t.buf = append(t.buf, nodeRB[K, V]{})
		n = &t.buf[len(t.buf)-1]
	}
	t.off++
	n.key = key
	n.value = value
	n.parent = nil
	n.left = nil
	n.right = nil
	return
}
