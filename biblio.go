package biblio

import (
	"fmt"
	"sort"
	"time"
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
	Base   []int
	Check  []int
	Fail   []int
	Output map[int][][]byte
}

// Compile TODO
func Compile(words [][]byte) *Matcher {
	m := new(Matcher)
	m.Base = make([]int, 2048)[:1]
	m.Check = make([]int, 2048)[:1]
	m.Fail = make([]int, 2048)[:1]
	m.Output = map[int][][]byte{}

	// Represents a node in the implicit trie representing words
	type trienode struct {
		state    int
		suffixes indexedStringSlice
	}
	queue := make([]trienode, 256)[:1]
	queue[0] = trienode{0, indexedStringSlice{words, 0}}

	sorttime := 0
	basecalctime := 0
	bfstime := 0
	totalchild := 0
	failtime := 0

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		depth := node.suffixes.depth

		startSort := time.Now().Nanosecond()
		// Get all the edges in lexicographical order
		var edges []byte
		sort.Sort(&node.suffixes)
		for _, suffix := range node.suffixes.slices {
			edge := suffix[depth]
			if len(edges) == 0 || edges[len(edges)-1] != edge {
				edges = append(edges, edge)
			}
		}
		sorttime += time.Now().Nanosecond() - startSort

		startBaseCalc := time.Now().Nanosecond()
		// Calculate a suitable Base value where each edge will fit into the
		// double array trie
		Base := m.findBase(edges)
		m.Base[node.state] = Base
		basecalctime += time.Now().Nanosecond() - startBaseCalc

		starttotalchild := time.Now().Nanosecond()
		i := 0
		for _, edge := range edges {
			offset := int(edge)
			newState := Base + offset

			// Setup the state=Check[Base[state]+offset] identity so we know
			// this edge exists in the trie
			// We always increase the state held in Check by 1 to avoid zero
			// values
			m.Check[newState] = node.state + 1

			startfailtime := time.Now().Nanosecond()
			// Setup the Fail function for the child nodes. This Check will
			// ensure nodes at level 0 and level 1 have a Fail state of 0
			if node.state != 0 {
				if m.hasEdge(m.Fail[node.state], offset) {
					// We can continue from the Fail state of the parent
					m.Fail[newState] = m.Base[m.Fail[node.state]] + offset
				} else if m.hasEdge(0, offset) {
					// We can continue from the Fail state of root
					m.Fail[newState] = m.Base[0] + offset
				}

				// Setup the Output function
				failState := m.Fail[newState]
				for _, word := range m.Output[failState] {
					m.Output[newState] = append(m.Output[newState], word)
				}
			}
			failtime += time.Now().Nanosecond() - startfailtime

			startBfs := time.Now().Nanosecond()
			// Add the child nodes to the queue to continue down the BFS
			newnode := trienode{newState, indexedStringSlice{[][]byte{}, depth + 1}}
			for i < len(node.suffixes.slices) && node.suffixes.slices[i][depth] == edge {
				if len(node.suffixes.slices[i]) > depth+1 {
					newnode.suffixes.slices = append(newnode.suffixes.slices, node.suffixes.slices[i])
				} else {
					m.Output[newState] = append(m.Output[newState], node.suffixes.slices[i])
				}
				i++
			}
			queue = append(queue, newnode)
			bfstime += time.Now().Nanosecond() - startBfs
		}
		totalchild += time.Now().Nanosecond() - starttotalchild
	}
	fmt.Printf("sorttime: %d ns\n", sorttime)
	fmt.Printf("basecalctime: %d ns\n", basecalctime)
	fmt.Printf("bfstime: %d ns\n", bfstime)
	fmt.Printf("totalchild: %d ns\n", totalchild)
	fmt.Printf("failtime: %d ns\n", failtime)

	return m
}

func (m *Matcher) findBase(edges []byte) int {
	if len(edges) == 0 {
		return LEAF
	}
	min := int(edges[0])
	max := int(edges[len(edges)-1])
	width := max - min

	if len(edges) < 3 {
		for i := range m.Check[1:] {
			i++ // fix i since we are using range [1:], simplifies calculations
			if i+width >= len(m.Check) {
				break
			}

			fits := true
			for _, e := range edges {
				if m.Check[i+int(e)-min] != 0 {
					fits = false
					break
				}
			}
			if fits {
				return i - min
			}
		}
	}

	m.increaseSize(width + 1)
	return len(m.Base) - 1 - max
}

func (m *Matcher) increaseSize(dsize int) {
	m.Base = append(m.Base, make([]int, dsize)...)
	m.Check = append(m.Check, make([]int, dsize)...)
	m.Fail = append(m.Fail, make([]int, dsize)...)
}

func (m *Matcher) hasEdge(fromState, offset int) bool {
	toState := m.Base[fromState] + offset
	return toState >= 0 && toState < len(m.Check) && m.Check[toState] == fromState+1
}

func hasPath(word []byte, m *Matcher) bool {
	state := 0
	for _, b := range word {
		Base := m.Base[state]
		if Base == LEAF {
			return false
		}
		if Base+int(b) >= len(m.Check) || m.Check[Base+int(b)]-1 != state {
			return false
		}
		state = Base + int(b)
	}
	return true
}

// Match TODO
type Match struct {
	Word  []byte
	Index int
}

// FindAll TODO
func (m *Matcher) FindAll(text []byte) (matches []Match) {
	state := 0
	for i, b := range text {
		offset := int(b)
		for {
			if m.hasEdge(state, offset) {
				state = m.Base[state] + offset
				break
			} else if state == 0 {
				break
			} else if m.Base[state] == LEAF {
				state = m.Fail[state]
			} else {
				state = m.Fail[state]
			}
		}
		for _, word := range m.Output[state] {
			matches = append(matches, Match{word, i})
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

// 	fmt.Printf("%v\n", m.Base)
// 	fmt.Printf("%v\n", m.Check)
// 	fmt.Printf("%v\n", m.Fail)
// 	for s, words := range m.Output {
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
