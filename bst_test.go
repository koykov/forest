package forest

import (
	"cmp"
	"fmt"
	"math/rand"
	"strings"
	"testing"
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
		type tc struct {
			name    string
			keys    []int
			query   int
			wantVal string
			wantOK  bool
		}
		tests := []tc{
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
				if ok != tt.wantOK {
					t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
				}
				if got != tt.wantVal {
					t.Fatalf("value = %q, want %q", got, tt.wantVal)
				}
			})
		}
	})

	t.Run("insert", func(t *testing.T) {
		type tc struct {
			name     string
			ops      []int
			wantSize int
		}
		tests := []tc{
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
				if got := tr.Size(); got != tt.wantSize {
					t.Fatalf("Size = %d, want %d", got, tt.wantSize)
				}
			})
		}
	})

	t.Run("upsert", func(t *testing.T) {
		tr := newTestBST(t)
		tr.Insert(5, "first")
		tr.Insert(5, "second")

		if got := tr.Size(); got != 1 {
			t.Fatalf("Size = %d, want 1", got)
		}
		v, ok := tr.Search(5)
		if !ok || v != "second" {
			t.Fatalf("Search = (%q, %v), want (\"second\", true)", v, ok)
		}
	})

	t.Run("delete", func(t *testing.T) {
		type tc struct {
			name      string
			keys      []int
			del       int
			wantOK    bool
			wantSize  int
			remaining []int
			absent    []int
		}
		tests := []tc{
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
				name: "delete node with two children", keys: []int{10, 5, 15, 3, 7, 12, 20}, del: 15,
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

				if ok := tr.Delete(tt.del); ok != tt.wantOK {
					t.Fatalf("Delete ok = %v, want %v", ok, tt.wantOK)
				}
				if got := tr.Size(); got != tt.wantSize {
					t.Fatalf("Size = %d, want %d", got, tt.wantSize)
				}
				for _, k := range tt.remaining {
					if _, ok := tr.Search(k); !ok {
						t.Fatalf("key %d should remain", k)
					}
				}
				for _, k := range tt.absent {
					if _, ok := tr.Search(k); ok {
						t.Fatalf("key %d should be absent", k)
					}
				}
			})
		}
	})

	t.Run("delete twice", func(t *testing.T) {
		tr := newTestBST(t)
		fill(t, tr, 5, 3, 7)
		if !tr.Delete(3) {
			t.Fatal("first delete should succeed")
		}
		if tr.Delete(3) {
			t.Fatal("second delete should fail")
		}
		if tr.Size() != 2 {
			t.Fatalf("Size = %d, want 2", tr.Size())
		}
	})

	t.Run("min max", func(t *testing.T) {
		type tc struct {
			name    string
			keys    []int
			wantMin int
			wantMax int
			empty   bool
		}
		tests := []tc{
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
				if ok == tt.empty {
					t.Fatalf("Min ok = %v, want empty=%v", ok, tt.empty)
				}
				if !tt.empty {
					if k != tt.wantMin {
						t.Fatalf("Min key = %d, want %d", k, tt.wantMin)
					}
					if v != fmt.Sprintf("v%d", tt.wantMin) {
						t.Fatalf("Min val = %q", v)
					}
				}

				k, _, ok = tr.Max()
				if ok == tt.empty {
					t.Fatalf("Max ok = %v, want empty=%v", ok, tt.empty)
				}
				if !tt.empty && k != tt.wantMax {
					t.Fatalf("Max key = %d, want %d", k, tt.wantMax)
				}
			})
		}
	})

	t.Run("clear", func(t *testing.T) {
		tr := newTestBST(t)
		fill(t, tr, 5, 3, 7)
		tr.Clear()

		if tr.Size() != 0 {
			t.Fatalf("Size = %d, want 0", tr.Size())
		}
		if _, ok := tr.Search(5); ok {
			t.Fatal("Search should miss after Clear")
		}
		if _, _, ok := tr.Min(); ok {
			t.Fatal("Min should be empty after Clear")
		}

		tr.Insert(1, "v1")
		if v, ok := tr.Search(1); !ok || v != "v1" {
			t.Fatal("tree not usable after Clear")
		}
	})

	t.Run("in-order", func(t *testing.T) {
		type tc struct {
			name string
			ops  []int
		}
		tests := []tc{
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
					if keys[i-1] >= keys[i] {
						t.Fatalf("order broken at %d: %v", i, keys)
					}
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
			if results[i] != results[0] {
				t.Fatalf("permutation %d gives %q, want %q", i, results[i], results[0])
			}
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
