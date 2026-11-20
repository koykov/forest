package forest

func (t *rb[K, V]) insert1(node *nodeRB[K, V]) {
	if node.parent == nil {
		node.color = rbBlack
	} else {
		t.insert2(node)
	}
}

func (t *rb[K, V]) insert2(node *nodeRB[K, V]) {
	if node.parent.color == rbBlack {
		return
	}
	t.insert3(node)
}

func (t *rb[K, V]) insert3(node *nodeRB[K, V]) {
	u := t.uncle(node)
	if u != nil && u.color == rbRed {
		node.parent.color = rbBlack
		u.color = rbBlack
		if g := t.grandparent(node); g != nil {
			g.color = rbRed
			t.insert1(node)
		}
	} else {
		t.insert4(node)
	}
}

func (t *rb[K, V]) insert4(node *nodeRB[K, V]) {
	g := t.grandparent(node)
	if g == nil {
		return
	}
	switch {
	case node.parent.right == node && node.parent == g.left:
		node = node.left
	case node.parent.left == node && node.parent == g.right:
		node = node.right
	}
	t.insert5(node)
}

func (t *rb[K, V]) insert5(node *nodeRB[K, V]) {
	g := t.grandparent(node)
	if g == nil {
		return
	}
	node.parent.color = rbBlack
	g.color = rbRed
	if node.parent.left == node && node.parent == g.left {
		t.rotr(g)
	} else {
		t.rotl(g)
	}
}
