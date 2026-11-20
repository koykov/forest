package forest

func (t *rb[K, V]) rebalance(node *nodeRB[K, V]) {
	if node.parent == nil {
		node.color = rbBlack
	} else {
		t.rebalance2(node)
	}
}

func (t *rb[K, V]) rebalance2(node *nodeRB[K, V]) {
	if node.parent.color == rbBlack {
		return
	}
	t.rebalance3(node)
}

func (t *rb[K, V]) rebalance3(node *nodeRB[K, V]) {
	u := t.uncle(node)
	if u != nil && u.color == rbRed {
		node.parent.color = rbBlack
		u.color = rbBlack
		if g := t.grandparent(node); g != nil {
			g.color = rbRed
			t.rebalance(node)
		}
	} else {
		t.rebalance4(node)
	}
}

func (t *rb[K, V]) rebalance4(node *nodeRB[K, V]) {
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
	t.rebalance5(node)
}

func (t *rb[K, V]) rebalance5(node *nodeRB[K, V]) {
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
