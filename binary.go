package forest

import "cmp"

type Binary[K cmp.Ordered, V any] interface {
	Insert(key K, value V)
	Delete(key K) bool
	Search(key K) (V, bool)
	Min() (K, V, bool)
	Max() (K, V, bool)
	Size() int
	Clear()
}
