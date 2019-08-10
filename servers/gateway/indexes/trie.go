package indexes

import (
	"sort"
	"sync"
)

//TODO: implement a trie data structure that stores
//keys of type string and values of type int64
type RuneSlice []rune

func (rs RuneSlice) Len() int {
	return len(rs)
}

func (rs RuneSlice) Less(i, j int) bool {
	return rs[i] < rs[j]
}

func (rs RuneSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

type Trie struct {
	root *Node
	size int
	mx   sync.RWMutex
}

type Node struct {
	values   []int64
	children map[rune]*Node
}

func NewTrie() *Trie {
	return &Trie{
		root: &Node{
			values: make([]int64, 0),
			children: make(map[rune]*Node),
		},
		size: 0,
	}
}

func (t *Trie) Len() int {
	return t.size
}

func (t *Trie) Add(key string, value int64) {
	t.mx.Lock()
	if len(key) == 0 {
		return
	}
	addHelper(0, key, value, t.root, t)
	t.mx.Unlock()
}

func addHelper(index int, key string, value int64, curr *Node, root *Trie) {
	if index == len(key) {
		if !contains(value, curr.values) {
			curr.values = append(curr.values, value)
		}
		sort.Slice(curr.values, func(i, j int) bool {
			return curr.values[i] < curr.values[j]
		})
		return
	}

	r := rune(key[index])
	var next *Node
	if curr.children == nil {
		curr.children = make(map[rune]*Node)
	}
	next, _ = curr.children[r]
	if next == nil {
		curr.children[r] = &Node{}
		next = curr.children[r]
		root.size++
	}

	index++
	addHelper(index, key, value, next, root)
}

func (t *Trie) Find(n int, prefix string) []int64 {
	t.mx.RLock()
	defer t.mx.RUnlock()
	if n == 0 || len(t.root.children) == 0 || prefix == "" {
		return nil
	}

	curr := t.root
	for _, character := range prefix {
		next, _ := curr.children[rune(character)]
		if next == nil {
			return nil
		}
		curr = next
	}

	var arr []int64
	findHelper(curr, n, &arr)
	if len(arr) < n {
		return arr
	}
	return arr[:n]
}

func findHelper(curr *Node, n int, arr *[]int64) {
	if curr.values != nil {
		for _, v := range curr.values {
			*arr = append(*arr, v)
		}
	}
	if len(curr.children) == 0 {
		return
	}

	var keys []rune
	for k := range curr.children {
		keys = append(keys, k)
	}
	sort.Sort(RuneSlice(keys))

	for _, k := range keys {
		findHelper(curr.children[k], n, arr)
	}
	return
}

func (t *Trie) Remove(key string, value int64) {
	t.mx.Lock()
	curr := t.root
	for _, character := range key {
		next, _ := curr.children[rune(character)]
		if next == nil {
			return
		}
		curr = next
	}

	if curr.values == nil || !contains(value, curr.values) {
		return
	}
	if curr.values != nil && len(curr.children) != 0 {
		curr.values = nil
		return
	}

	removeHelper(t.root, t, key)
	t.mx.Unlock()
}

func removeHelper(curr *Node, t *Trie, key string) {
	if len(key) == 0 {
		return
	}

	character := rune(key[0])
	current := curr.children[character]
	curr.children[character] = nil
	t.size--
	removeHelper(current, t, key[1:])
}

func contains(element int64, slice []int64) bool {
	for _, ele := range slice {
		return ele == element
	}
	return false
}