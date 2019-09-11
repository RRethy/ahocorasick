package biblio

import (
	"sort"
)

const (
	// LEAF represents a leaf on the trie
	// This must be <255 since the offsets used are in [0,255]
	LEAF = -1867
)

type indexedStringSlice struct {
	slices [][]byte
	depth  int
}

func (sslice *indexedStringSlice) Len() int {
	return len(sslice.slices)
}
func (sslice *indexedStringSlice) Less(i, j int) bool {
	return sslice.slices[i][sslice.depth] < sslice.slices[j][sslice.depth]
}
func (sslice *indexedStringSlice) Swap(i, j int) {
	sslice.slices[i], sslice.slices[j] = sslice.slices[j], sslice.slices[i]
}

// Matcher TODO
type Matcher struct {
	base   []int
	check  []int
	fail   []int
	output map[int][][]byte
}

// Compile TODO
func Compile(words [][]byte) *Matcher {
	m := new(Matcher)
	m.base = append(m.base, 0)
	m.check = append(m.check, 0)
	m.fail = append(m.fail, 0)
	m.output = map[int][][]byte{}

	// Represents a node in the implicit trie representing words
	type trienode struct {
		state    int
		suffixes indexedStringSlice
	}
	queue := []trienode{{0, indexedStringSlice{words, 0}}}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		depth := node.suffixes.depth

		// Get all the edges in lexicographical order
		var edges []byte
		sort.Sort(&node.suffixes)
		for _, suffix := range node.suffixes.slices {
			edge := suffix[depth]
			if len(edges) == 0 || edges[len(edges)-1] != edge {
				edges = append(edges, edge)
			}
		}

		// Calculate a suitable base value where each edge will fit into the
		// double array trie
		base := m.findBase(edges)
		m.base[node.state] = base

		i := 0
		for _, edge := range edges {
			offset := int(edge)
			newState := base + offset

			// Setup the state=check[base[state]+offset] identity so we know
			// this edge exists in the trie
			// We always increase the state held in check by 1 to avoid zero
			// values
			m.check[newState] = node.state + 1

			// Setup the fail function for the child nodes. This check will
			// ensure nodes at level 0 and level 1 have a fail state of 0
			if node.state != 0 {
				if m.hasEdge(m.fail[node.state], offset) {
					// We can continue from the fail state of the parent
					m.fail[newState] = m.base[m.fail[node.state]] + offset
				} else if m.hasEdge(0, offset) {
					// We can continue from the fail state of root
					m.fail[newState] = m.base[0] + offset
				}

				// Setup the output function
				failState := m.fail[newState]
				for _, word := range m.output[failState] {
					m.output[newState] = append(m.output[newState], word)
				}
			}

			// Add the child nodes to the queue to continue down the BFS
			newnode := trienode{newState, indexedStringSlice{[][]byte{}, depth + 1}}
			for i < len(node.suffixes.slices) && node.suffixes.slices[i][depth] == edge {
				if len(node.suffixes.slices[i]) > depth+1 {
					newnode.suffixes.slices = append(newnode.suffixes.slices, node.suffixes.slices[i])
				} else {
					m.output[newState] = append(m.output[newState], node.suffixes.slices[i])
				}
				i++
			}
			queue = append(queue, newnode)
		}
	}

	return m
}

func (m *Matcher) findBase(edges []byte) int {
	if len(edges) == 0 {
		return LEAF
	}
	min := int(edges[0])
	max := int(edges[len(edges)-1])
	width := max - min

	for i := range m.check[1:] {
		i++ // fix i since we are using range [1:], simplifies calculations
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

func (m *Matcher) increaseSize(dsize int) {
	m.base = append(m.base, make([]int, dsize)...)
	m.check = append(m.check, make([]int, dsize)...)
	m.fail = append(m.fail, make([]int, dsize)...)
}

func (m *Matcher) hasEdge(fromState, offset int) bool {
	toState := m.base[fromState] + offset
	return toState >= 0 && toState < len(m.check) && m.check[toState] == fromState+1
}

func hasPath(word []byte, m *Matcher) bool {
	state := 0
	for _, b := range word {
		base := m.base[state]
		if base == LEAF {
			return false
		}
		if base+int(b) >= len(m.check) || m.check[base+int(b)]-1 != state {
			return false
		}
		state = base + int(b)
	}
	return true
}

// Match TODO
type Match struct {
	word  []byte
	index int
}

// FindAll TODO
func (m *Matcher) FindAll(text []byte) (matches []Match) {
	state := 0
	for i, b := range text {
		offset := int(b)
		for {
			if m.base[state] == LEAF {
				state = m.fail[state]
			} else if m.hasEdge(state, offset) {
				state = m.base[state] + offset
				break
			} else if state == 0 {
				break
			} else {
				state = m.fail[state]
			}
		}
		if m.hasEdge(state, offset) {
			state = m.base[state] + offset
		}
		for _, word := range m.output[state] {
			matches = append(matches, match{word, i})
		}
	}
	return
}

// func main() {
// 	// m, _ := Compile([]string{"hers", "she"})
// 	m := Compile([][]byte{
// 		[]byte("he"),
// 		[]byte("she"),
// 		[]byte("they"),
// 		[]byte("their"),
// 		[]byte("where"),
// 		[]byte("bear"),
// 		[]byte("taratula"),
// 		[]byte("adam"),
// 		[]byte("regard-rethy"),
// 		[]byte("panda"),
// 		[]byte("bear"),
// 		[]byte("golang"),
// 		[]byte("his"),
// 		[]byte("hers"),
// 		[]byte("her"),
// 	})

// 	fmt.Printf("%v\n", m.base)
// 	fmt.Printf("%v\n", m.check)
// 	fmt.Printf("%v\n", m.fail)
// 	for s, words := range m.output {
// 		fmt.Printf("%d =>\n", s)
// 		for _, word := range words {
// 			fmt.Println(string(word))
// 		}
// 	}

// 	fmt.Println("beshe hers ")
// 	matches := m.findAll([]byte("beshe hers "))
// 	for _, match := range matches {
// 		fmt.Printf("%d - %s\n", match.index, string(match.word))
// 	}
// }
