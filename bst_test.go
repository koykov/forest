package forest

import (
	"cmp"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func inorderKeys[K cmp.Ordered, V any](dst []K, n *nodeBST[K, V]) []K {
	if n == nil {
		return dst
	}
	dst = inorderKeys(dst, n.left)
	dst = append(dst, n.key)
	dst = inorderKeys(dst, n.right)
	return dst
}

func TestBST(t *testing.T) {
	type testBST = Binary[int, string]

	newTestBST := func(t *testing.T) testBST {
		t.Helper()
		return NewBST[int, string]()
	}

	fill := func(t *testing.T, tr testBST, keys ...int) {
		t.Helper()
		for _, k := range keys {
			tr.Insert(k, fmt.Sprintf("v%d", k))
		}
	}

	t.Run("search", func(t *testing.T) {
		tests := []struct {
			name    string
			keys    []int
			query   int
			wantVal string
			wantOK  bool
		}{
			{"empty tree", nil, 1, "", false},
			{"single hit", []int{5}, 5, "v5", true},
			{"single miss smaller", []int{5}, 1, "", false},
			{"single miss bigger", []int{5}, 9, "", false},
			{"left child hit", []int{5, 3}, 3, "v3", true},
			{"right child hit", []int{5, 7}, 7, "v7", true},
			{"miss between", []int{5, 3, 7}, 4, "", false},
			{"deep left", []int{5, 3, 1}, 1, "v1", true},
			{"deep right", []int{5, 7, 9}, 9, "v9", true},
			{"miss below min", []int{5, 3, 7}, 0, "", false},
			{"miss above max", []int{5, 3, 7}, 100, "", false},
			{"root after many", []int{10, 5, 15, 3, 7, 12, 20}, 10, "v10", true},
			{"leaf", []int{10, 5, 15, 3, 7, 12, 20}, 12, "v12", true},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tree := newTestBST(t)
				fill(t, tree, tt.keys...)

				got, ok := tree.Search(tt.query)
				require.Equal(t, tt.wantOK, ok)
				assert.Equal(t, tt.wantVal, got)
			})
		}
	})

	t.Run("insert", func(t *testing.T) {
		tests := []struct {
			name     string
			ops      []int
			wantSize int
		}{
			{"empty", nil, 0},
			{"one", []int{5}, 1},
			{"two distinct", []int{5, 3}, 2},
			{"increasing", []int{1, 2, 3, 4, 5}, 5},
			{"decreasing", []int{5, 4, 3, 2, 1}, 5},
			{"random-ish", []int{5, 3, 7, 1, 4, 6, 8}, 7},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tr := newTestBST(t)
				fill(t, tr, tt.ops...)
				assert.Equal(t, tt.wantSize, tr.Size())
			})
		}
	})

	t.Run("upsert", func(t *testing.T) {
		tr := newTestBST(t)
		tr.Insert(5, "first")
		tr.Insert(5, "second")

		require.Equal(t, 1, tr.Size())
		v, ok := tr.Search(5)
		require.True(t, ok)
		assert.Equal(t, "second", v)
	})

	t.Run("delete", func(t *testing.T) {
		tests := []struct {
			name      string
			keys      []int
			del       int
			wantOK    bool
			wantSize  int
			remaining []int
			absent    []int
		}{
			{
				name: "empty tree", keys: nil, del: 1,
				wantOK: false, wantSize: 0,
			},
			{
				name: "delete missing", keys: []int{5, 3, 7}, del: 4,
				wantOK: false, wantSize: 3, remaining: []int{5, 3, 7},
			},
			{
				name: "delete leaf", keys: []int{5, 3, 7}, del: 3,
				wantOK: true, wantSize: 2, remaining: []int{5, 7}, absent: []int{3},
			},
			{
				name: "delete node with one left child", keys: []int{5, 3, 1}, del: 3,
				wantOK: true, wantSize: 2, remaining: []int{5, 1}, absent: []int{3},
			},
			{
				name: "delete node with one right child", keys: []int{5, 7, 9}, del: 7,
				wantOK: true, wantSize: 2, remaining: []int{5, 9}, absent: []int{7},
			},
			{
				name: "delete node with two children", keys: []int{10, 5, 15, 3, 7, 12, 20}, del: 10,
				wantOK: true, wantSize: 6,
				remaining: []int{15, 5, 3, 7, 12, 20}, absent: []int{10},
			},
			{
				name: "delete node with two children (non-root)", keys: []int{10, 5, 15, 3, 7, 12, 20}, del: 15,
				wantOK: true, wantSize: 6,
				remaining: []int{10, 5, 3, 7, 12, 20}, absent: []int{15},
			},
			{
				name: "delete root leaf", keys: []int{1}, del: 1,
				wantOK: true, wantSize: 0, absent: []int{1},
			},
			{
				name: "delete root with one child", keys: []int{5, 3}, del: 5,
				wantOK: true, wantSize: 1, remaining: []int{3}, absent: []int{5},
			},
			{
				name: "delete root with two children", keys: []int{5, 3, 7}, del: 5,
				wantOK: true, wantSize: 2, remaining: []int{3, 7}, absent: []int{5},
			},
			{
				name: "delete min", keys: []int{5, 3, 7, 1}, del: 1,
				wantOK: true, wantSize: 3, remaining: []int{5, 3, 7}, absent: []int{1},
			},
			{
				name: "delete max", keys: []int{5, 3, 7, 9}, del: 9,
				wantOK: true, wantSize: 3, remaining: []int{5, 3, 7}, absent: []int{9},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tr := newTestBST(t)
				fill(t, tr, tt.keys...)

				ok := tr.Delete(tt.del)
				require.Equal(t, tt.wantOK, ok, "Delete ok")
				require.Equal(t, tt.wantSize, tr.Size(), "Size after delete")

				for _, k := range tt.remaining {
					_, found := tr.Search(k)
					assert.Truef(t, found, "key %d should remain", k)
				}
				for _, k := range tt.absent {
					_, found := tr.Search(k)
					assert.Falsef(t, found, "key %d should be absent", k)
				}
			})
		}
	})

	t.Run("delete twice", func(t *testing.T) {
		tr := newTestBST(t)
		fill(t, tr, 5, 3, 7)

		require.True(t, tr.Delete(3), "first delete should succeed")
		require.False(t, tr.Delete(3), "second delete should fail")
		assert.Equal(t, 2, tr.Size())
	})

	t.Run("min max", func(t *testing.T) {
		tests := []struct {
			name    string
			keys    []int
			wantMin int
			wantMax int
			empty   bool
		}{
			{"empty", nil, 0, 0, true},
			{"single", []int{5}, 5, 5, false},
			{"left chain", []int{5, 4, 3, 2, 1}, 1, 5, false},
			{"right chain", []int{1, 2, 3, 4, 5}, 1, 5, false},
			{"balanced", []int{10, 5, 15, 3, 7, 12, 20}, 3, 20, false},
			{"after deletes", []int{10, 5, 15, 3, 7}, 3, 15, false},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tr := newTestBST(t)
				fill(t, tr, tt.keys...)

				k, v, ok := tr.Min()
				require.Equal(t, !tt.empty, ok, "Min ok")
				if !tt.empty {
					assert.Equal(t, tt.wantMin, k, "Min key")
					assert.Equal(t, fmt.Sprintf("v%d", tt.wantMin), v, "Min value")
				}

				k, _, ok = tr.Max()
				require.Equal(t, !tt.empty, ok, "Max ok")
				if !tt.empty {
					assert.Equal(t, tt.wantMax, k, "Max key")
				}
			})
		}
	})

	t.Run("clear", func(t *testing.T) {
		tr := newTestBST(t)
		fill(t, tr, 5, 3, 7)
		tr.Clear()

		require.Equal(t, 0, tr.Size())
		_, ok := tr.Search(5)
		assert.False(t, ok, "Search should miss after Clear")
		_, _, ok = tr.Min()
		assert.False(t, ok, "Min should be empty after Clear")

		tr.Insert(1, "v1")
		v, ok := tr.Search(1)
		require.True(t, ok, "tree not usable after Clear")
		assert.Equal(t, "v1", v)
	})

	t.Run("in-order", func(t *testing.T) {
		tests := []struct {
			name string
			ops  []int
		}{
			{"increasing", []int{1, 2, 3, 4, 5}},
			{"decreasing", []int{5, 4, 3, 2, 1}},
			{"zigzag", []int{5, 1, 4, 2, 3}},
			{"with deletes", []int{10, 5, 15, 3, 7, 12, 20, -5, -15}},
			{"delete all", []int{5, 3, 7, -5, -3, -7}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tr := newTestBST(t)
				for _, op := range tt.ops {
					if op < 0 {
						tr.Delete(-op)
					} else {
						tr.Insert(op, fmt.Sprintf("v%d", op))
					}
				}
				tree := tr.(*bst[int, string])
				keys := inorderKeys(nil, tree.root)
				for i := 1; i < len(keys); i++ {
					require.Lessf(t, keys[i-1], keys[i], "order broken at %d: %v", i, keys)
				}
			})
		}
	})

	t.Run("order", func(t *testing.T) {
		permutations := [][]int{
			{5, 3, 7, 1, 4, 6, 8},
			{1, 4, 6, 8, 3, 7, 5},
			{8, 6, 4, 1, 7, 3, 5},
		}
		var results []string
		for _, keys := range permutations {
			tr := newTestBST(t)
			fill(t, tr, keys...)
			var got []string
			for _, k := range []int{1, 3, 4, 5, 6, 7, 8} {
				v, _ := tr.Search(k)
				got = append(got, v)
			}
			results = append(results, strings.Join(got, ","))
		}
		for i := 1; i < len(results); i++ {
			assert.Equalf(t, results[0], results[i], "permutation %d", i)
		}
	})
}

func BenchmarkBST(b *testing.B) {
	sizes := []int{100, 1_000, 10_000, 100_000}

	b.Run("insert", func(b *testing.B) {
		b.Run("random", func(b *testing.B) {
			for _, n := range sizes {
				b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
					keys := rand.Perm(n)
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						tr := NewBST[int, string]()
						for _, k := range keys {
							tr.Insert(k, "v")
						}
					}
				})
			}
		})
		b.Run("sort", func(b *testing.B) {
			for _, n := range sizes {
				b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						tr := NewBST[int, string]()
						for k := 0; k < n; k++ {
							tr.Insert(k, "v")
						}
					}
				})
			}
		})
	})

	b.Run("search", func(b *testing.B) {
		b.Run("hit", func(b *testing.B) {
			for _, n := range sizes {
				b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
					tr := NewBST[int, string]()
					keys := rand.Perm(n)
					for _, k := range keys {
						tr.Insert(k, "v")
					}
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						tr.Search(keys[i%n])
					}
				})
			}
		})
		b.Run("miss", func(b *testing.B) {
			for _, n := range sizes {
				b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
					tr := NewBST[int, string]()
					for k := 0; k < n; k++ {
						tr.Insert(k*2, "v")
					}
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						tr.Search(i%n*2 + 1)
					}
				})
			}
		})
	})

	b.Run("delete", func(b *testing.B) {
		for _, n := range sizes {
			b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					tr := NewBST[int, string]()
					keys := rand.Perm(n)
					for _, k := range keys {
						tr.Insert(k, "v")
					}
					b.StartTimer()
					for _, k := range keys {
						tr.Delete(k)
					}
				}
			})
		}
	})

	b.Run("min max", func(b *testing.B) {
		const n = 10_000
		tr := NewBST[int, string]()
		for _, k := range rand.Perm(n) {
			tr.Insert(k, "v")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tr.Min()
			tr.Max()
		}
	})
}
