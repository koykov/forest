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
}

func (t *bst[K, V]) delete(parent, node *nodeBST[K, V], key K) bool {
	switch {
	case node.key == key:
		switch {
		case node.left == nil && node.right == nil:
			// Delete node without children.
			switch {
			case parent == nil:
				t.root = nil
			case parent.key < key:
				parent.right = nil
			case parent.key > key:
				parent.left = nil
			}
		case node.left != nil && node.right == nil:
			// Move left branch upside (overwrite current node).
			switch {
			case parent == nil:
				t.root = node.left
			case parent.key < key:
				parent.right = node.left
			case parent.key > key:
				parent.left = node.left
			}
		case node.left == nil && node.right != nil:
			// Move right branch upside (overwrite current node).
			switch {
			case parent == nil:
				t.root = node.right
			case parent.key < key:
				parent.right = node.right
			case parent.key > key:
				parent.left = node.right
			}
		case node.left != nil && node.right != nil:
			// Replace current node with successor.
			z := node
			s := node.right
			for {
				if left := s.left; left != nil {
					z = s
					s = left
					continue
				}
				break
			}
			if s.right != nil {
				z.left = s.right
			}
			node.key, node.value = s.key, s.value
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
