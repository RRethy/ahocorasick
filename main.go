package main

import (
	"fmt"
	"sort"
)

type indexedStringSlice struct {
	strs  [][]byte
	depth int
}

func (sslice *indexedStringSlice) Len() int {
	return len(sslice.strs)
}
func (sslice *indexedStringSlice) Less(i, j int) bool {
	return sslice.strs[i][sslice.depth] < sslice.strs[j][sslice.depth]
}
func (sslice *indexedStringSlice) Swap(i, j int) {
	sslice.strs[i], sslice.strs[j] = sslice.strs[j], sslice.strs[i]
}

type matcher struct {
	base   []int
	check  []int
	fail   []int
	output map[int][][]byte
}

func compileMatcher(words [][]byte) (*matcher, error) {
	m := new(matcher)
	m.base = append(m.base, 0)
	m.check = append(m.check, 0)
	m.fail = append(m.fail, 0)
	m.output = map[int][][]byte{}

	type tnode struct {
		state    int
		suffixes indexedStringSlice
	}
	queue := []tnode{{0, indexedStringSlice{words, 0}}}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		depth := node.suffixes.depth

		// Get all the edges
		sort.Sort(&node.suffixes)
		var edges []byte
		for _, suffix := range node.suffixes.strs {
			edge := suffix[depth]
			if len(edges) == 0 || edges[len(edges)-1] != edge {
				edges = append(edges, edge)
			}
		}

		base := m.findBase(edges)
		m.base[node.state] = base

		i := 0
		for _, edge := range edges {
			offset := int(edge)

			m.check[base+offset] = node.state + 1
			if node.state != 0 {
				if m.hasEdge(m.fail[node.state], offset) {
					m.fail[base+offset] = m.base[m.fail[node.state]] + offset
				} else if m.hasEdge(0, offset) {
					m.fail[base+offset] = m.base[0] + offset
				}
			}

			newnode := tnode{base + offset, indexedStringSlice{[][]byte{}, depth + 1}}
			for i < len(node.suffixes.strs) && node.suffixes.strs[i][depth] == edge {
				if len(node.suffixes.strs[i]) > depth+1 {
					newnode.suffixes.strs = append(newnode.suffixes.strs, node.suffixes.strs[i])
				} else {
					m.output[newnode.state] = append(m.output[newnode.state], node.suffixes.strs[i])
				}
				i++
			}
			queue = append(queue, newnode)
		}
	}

	return m, nil
}

func (m *matcher) findBase(edges []byte) int {
	if len(edges) == 0 {
		return -300
	}
	min := int(edges[0])
	max := int(edges[len(edges)-1])
	width := max - min

	for i := range m.check[1:] {
		i++ // fix i since we are using range [1:]
		if i+width >= len(m.check) {
			break
		}

		fits := true
		for _, e := range edges {
			if m.check[i+int(e)-min] != 0 {
				fits = false
				break
			}
		}
		if fits {
			return i - min
		}
	}

	m.increaseSize(width + 1)
	return len(m.base) - 1 - max
}

func (m *matcher) increaseSize(dsize int) {
	m.base = append(m.base, make([]int, dsize)...)
	m.check = append(m.check, make([]int, dsize)...)
	m.fail = append(m.fail, make([]int, dsize)...)
}

func (m *matcher) hasEdge(fromState, offset int) bool {
	toState := m.base[fromState] + offset
	return toState < len(m.check) && m.check[toState] == fromState+1
}

func hasPath(word []byte, m *matcher) bool {
	state := 0
	for _, b := range word {
		base := m.base[state]
		if base == -300 {
			return false
		}
		if base+int(b) >= len(m.check) || m.check[base+int(b)]-1 != state {
			return false
		}
		state = base + int(b)
	}
	return true
}

func main() {
	// m, _ := compileMatcher([]string{"hers", "she"})
	m, _ := compileMatcher([][]byte{
		[]byte("he"),
		[]byte("hers"),
		[]byte("his"),
		[]byte("she"),
		[]byte("be"),
	})

	fmt.Println(hasPath([]byte("hers"), m))
	fmt.Println(hasPath([]byte("she"), m))
	fmt.Println(hasPath([]byte("his"), m))
	fmt.Println(hasPath([]byte("he"), m))
	fmt.Println(hasPath([]byte("be"), m))
	fmt.Println(hasPath([]byte("herss"), m))

	fmt.Printf("%v\n", m.base)
	fmt.Printf("%v\n", m.check)
	fmt.Printf("%v\n", m.fail)
	for s, words := range m.output {
		fmt.Printf("%d =>\n", s)
		for _, word := range words {
			fmt.Println(string(word))
		}
	}
}
