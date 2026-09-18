package forest

import "cmp"

type bst[K cmp.Ordered, V any] struct {
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

func NewBST[K cmp.Ordered, V any]() Binary[K, V] {
	return &bst[K, V]{}
}

func (t *bst[K, V]) Insert(key K, value V) {
	if t.root == nil {
		t.root = &nodeBST[K, V]{
			key:   key,
			value: value,
		}
		t.size++
		return
	}
	t.insert(t.root, key, value)
}

func (t *bst[K, V]) insert(node *nodeBST[K, V], key K, value V) {
	switch {
	case node.key == key:
		node.value = value
	case node.key > key:
		if node.left == nil {
			node.left = &nodeBST[K, V]{
				key:   key,
				value: value,
			}
			t.size++
			return
		}
		t.insert(node.left, key, value)
	case node.key < key:
		if node.right == nil {
			node.right = &nodeBST[K, V]{
				key:   key,
				value: value,
			}
			t.size++
			return
		}
		t.insert(node.right, key, value)
	}
}

func (t *bst[K, V]) Delete(key K) bool {
	if t.root == nil {
		return false
	}
	return t.delete(nil, t.root, key)
	// switch {
	// case t.root.key == key:
	// 	t.root = nil
	// 	t.size--
	// 	return true
	// case t.root.key < key && t.root.right != nil:
	// 	return t.delete(t.root, t.root.right, key)
	// case t.root.key > key && t.root.left != nil:
	// 	return t.delete(t.root, t.root.left, key)
	// }
	// return false
}

func (t *bst[K, V]) delete(parent, node *nodeBST[K, V], key K) bool {
	switch {
	case node.key == key:
		if parent == nil {
			t.root = nil
			t.size--
			return true
		}
		ptr := parent.left
		if parent.key < key {
			ptr = parent.right
		}
		_ = ptr
		switch {
		case node.left == nil && node.right == nil:
			ptr = nil
		case node.left != nil && node.right == nil:
			ptr = node.left
		case node.left == nil && node.right != nil:
			ptr = node.right
		case node.left != nil && node.right != nil:
			successor := node.right
			for left := successor.left; left != nil; {
				successor = left
			}
			successor.left = node.left
			ptr = node.right
			// todo finish me
		}
		t.size--
		return true
	case node.key < key && node.right != nil:
		return t.delete(node, node.right, key)
	case node.key > key && node.left != nil:
		return t.delete(node, node.left, key)
	}
	return false
}

func (t *bst[K, V]) Search(key K) (V, bool) {
	if t.root == nil {
		return t.nullV, false
	}
	node, ok := t.search(t.root, key)
	if !ok {
		return t.nullV, false
	}
	return node.value, true
}

func (t *bst[K, V]) search(node *nodeBST[K, V], key K) (*nodeBST[K, V], bool) {
	switch {
	case node.key == key:
		return node, true
	case node.key < key && node.right != nil:
		return t.search(node.right, key)
	case node.key > key && node.left != nil:
		return t.search(node.left, key)
	}
	return nil, false
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
	return t.size
}

func (t *bst[K, V]) Clear() {
	t.root = nil
	t.size = 0
}
