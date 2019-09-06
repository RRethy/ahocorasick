package main

import (
	"fmt"
	"sort"
)

type matcher struct {
	base   []int
	check  []int
	fail   []int
	output map[int]string
}

func compileMatcher(words []string) (*matcher, error) {
	m := new(matcher)
	m.base = append(m.base, 0)
	m.check = append(m.check, 0)
	m.fail = append(m.fail, 0)

	type tnode struct {
		state    int
		suffixes []string
	}
	queue := []tnode{{0, words}}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		sort.Strings(node.suffixes)
		var edges []byte
		for _, suffix := range node.suffixes {
			if len(suffix) == 0 {
				return nil, nil // TODO return an error
			}

			edge := suffix[0]
			if len(edges) == 0 || edges[len(edges)-1] != edge {
				edges = append(edges, edge)
			}
		}

		base := m.findBase(edges)
		m.base[node.state] = base
		i := 0
		parentFailState := m.fail[node.state]
		for _, edge := range edges {
			edgeVal := int(edge)

			m.check[base+edgeVal] = node.state + 1
			if node.state != 0 {
				parentFailContd := m.base[parentFailState] + edgeVal
				rootContd := m.base[0] + edgeVal
				if parentFailContd < len(m.check) && m.check[parentFailContd] == parentFailState {
					m.fail[base+edgeVal] = m.base[parentFailState] + edgeVal
				} else if rootContd < len(m.check) && m.check[rootContd] == 0 {
					m.fail[base+edgeVal] = m.base[0] + edgeVal
				}
			}

			newnode := tnode{base + edgeVal, []string{}}
			for i < len(node.suffixes) && node.suffixes[i][0] == edge {
				if len(node.suffixes[i]) > 1 {
					newnode.suffixes = append(newnode.suffixes, node.suffixes[i][1:])
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
	m, _ := compileMatcher([]string{"he", "hers", "his", "she", "be"})

	fmt.Println(hasPath([]byte("hers"), m))
	fmt.Println(hasPath([]byte("she"), m))
	fmt.Println(hasPath([]byte("his"), m))
	fmt.Println(hasPath([]byte("he"), m))
	fmt.Println(hasPath([]byte("be"), m))
	fmt.Println(hasPath([]byte("herss"), m))

	fmt.Printf("%v\n", m.base)
	fmt.Printf("%v\n", m.check)
	fmt.Printf("%v\n", m.fail)
}
