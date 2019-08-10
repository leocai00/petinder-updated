package indexes

import (
	"testing"
)

//TODO: implement automated tests for your trie data structure
func TestLen(t *testing.T) {
	cases := []struct {
		name  string
		key   string
		value int64
		trie  *Trie
	}{
		{
			"Match Sizes",
			"ab12",
			34,
			NewTrie(),
		},
	}

	for _, c := range cases {
		c.trie.Add(c.key, c.value)
		if c.trie.Len() != len(c.key) {
			t.Errorf("sizes do not match")
		}
	}
}

func TestAdd(t *testing.T) {
	cases := []struct {
		name  string
		key   string
		value []int64
		trie  *Trie
	}{
		{
			"Single",
			"ab12",
			[]int64{34},
			NewTrie(),
		},
		{
			"Multiple",
			"abcd",
			[]int64{34, 45, 56},
			NewTrie(),
		},
	}

	for _, c := range cases {
		for _, v := range c.value {
			c.trie.Add(c.key, v)
		}
		curr := c.trie.root
		values := make([]int64, len(c.value))
		for k, m := range c.key {
			r := rune(m)
			val, ok := curr.children[r]
			if !ok {
				t.Errorf("keys do not match")
			}
			if k < len(c.key)-1 {
				curr = curr.children[r]
			} else {
				copy(values, val.values)
			}
		}

		for k, v := range values {
			if c.value[k] != v {
				t.Errorf("values do not match")
			}
		}
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		name          string
		key           []string
		removedKey    string
		value         []int64
		removedValue  int64
		expectedValue int64
		deleteChild   bool
		expectedSize  int
	}{
		{
			"Single Input",
			[]string{"go"},
			"go",
			[]int64{1},
			1,
			0,
			true,
			0,
		},
		{
			"Empty Value",
			[]string{"gogoo", "go"},
			"go",
			[]int64{123, 12},
			12,
			0,
			false,
			5,
		},
		{
			"Empty Trie",
			[]string{},
			"go",
			[]int64{1},
			1,
			0,
			true,
			0,
		},
		{
			"Does Not Contain Value",
			[]string{"go"},
			"go",
			[]int64{1},
			12,
			0,
			false,
			2,
		},
	}

	for _, c := range cases {
		root := NewTrie()
		for _, k := range c.key {
			for _, v := range c.value {
				root.Add(k, v)
			}
		}

		root.Remove(c.removedKey, c.removedValue)
		current := root
		if c.deleteChild {
			r := rune(c.removedKey[0])
			if next := current.root.children[r]; next != nil {
				t.Errorf("case %s: should remove child node", c.name)
			}
		} else {
			for _, k := range c.removedKey {
				next, ok := current.root.children[rune(k)]
				if !ok {
					t.Errorf("case %s: did not add to child", c.name)
				}
				current.root = next
			}
		}
		if root.Len() != c.expectedSize {
			t.Errorf("case %s: expected size %d but got %d", c.name, c.expectedSize, root.size)
		}
	}
}
func TestFind(t *testing.T) {
	cases := []struct {
		name   string
		key1   string
		key2   string
		value1 []int64
		value2 []int64
		n      int
		trie   *Trie
	}{
		{
			"Single Value",
			"abcd",
			"",
			[]int64{10, 11},
			[]int64{11, 10},
			1,
			NewTrie(),
		},
		{
			"More Value",
			"aabbcs",
			"",
			[]int64{26, 40},
			[]int64{59, 34},
			3,
			NewTrie(),
		},
		{
			"Less Value",
			"aabbcs",
			"ab",
			[]int64{26, 40, 60},
			[]int64{34, 59},
			4,
			NewTrie(),
		},
		{
			"Same Value",
			"aabbcs",
			"ab",
			[]int64{26, 40, 60},
			[]int64{34, 59},
			3,
			NewTrie(),
		},
		{
			"No Given Rune",
			"ab",
			"cd",
			[]int64{10, 11, 12},
			[]int64{13, 14},
			3,
			NewTrie(),
		},
	}

	for _, c := range cases {
		for _, v := range c.value1 {
			c.trie.Add(c.key1, v)
		}
		c.trie.Add("kf", c.value2[0])
		result := c.trie.Find(c.n, "ab")
		temp := append(c.value1, c.value2...)
		for i := 0; i < len(result); i++ {
			if result[i] != temp[i] {
				t.Errorf("results do not match")
			}
		}
	}
}