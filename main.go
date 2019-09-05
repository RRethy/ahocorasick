package main

import (
	"fmt"
	"sort"
)

type matcher struct {
	base   []int
	check  []int
	output map[int]string
}

func compileMatcher(words []string) (*matcher, error) {
	m := new(matcher)
	m.base = append(m.base, 0)
	m.check = append(m.check, 0)

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
		for _, edge := range edges {
			m.check[base+int(edge)] = node.state + 1
		}

		i := 0
		for _, edge := range edges {
			newnode := tnode{base + int(edge), []string{}}
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
	} else if len(edges) < 3 {
		e0 := int(edges[0])
		e1 := int(edges[len(edges)-1])
		dx := e1 - e0

		for i, slot := range m.check[1:] {
			if i+dx+1 >= len(m.check) {
				break
			}
			if slot == 0 && m.check[i+dx+1] == 0 {
				return i - int(edges[0]) + 1
			}
		}
		m.increaseSize(dx + 1)
		return len(m.base) - 1 - e1
	}
	i := len(m.base) - 1
	m.increaseSize(256)
	return i
}

func (m *matcher) increaseSize(dsize int) {
	m.base = append(m.base, make([]int, dsize)...)
	m.check = append(m.check, make([]int, dsize)...)
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
}
