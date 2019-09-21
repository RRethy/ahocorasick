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

func (m *Matcher) String() string {
	return fmt.Sprintf(`
Base:   %v
Check:  %v
Fail:   %v
Output: %v
`, m.Base, m.Check, m.Fail, m.Output)
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

		// TODO if no edges then make it a leaf and continue
		// if node.suffixes.Len() == 0 {
		// 	m.Base[node.state] = LEAF
		// 	continue
		// }

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
		Base := m.Foofb(edges)
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
			// TODO this comment is wrong
			m.occupyState(newState, node.state)
			// m.Check[newState] = node.state

			startfailtime := time.Now().Nanosecond()
			if depth > 0 {
				m.setupfail(newState, node.state, offset)
			}
			m.unionFailOutput(newState, m.Fail[newState])
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
	// fmt.Printf("sorttime: %d ns\n", sorttime)
	// fmt.Printf("basecalctime: %d ns\n", basecalctime)
	// fmt.Printf("bfstime: %d ns\n", bfstime)
	// fmt.Printf("totalchild: %d ns\n", totalchild)
	// fmt.Printf("failtime: %d ns\n", failtime)

	return m
}

func (m *Matcher) occupyState(state, parentState int) {
	firstFreeState := m.firstFreeState()
	lastFreeState := m.lastFreeState()
	if firstFreeState == lastFreeState {
		m.Check[0] = 0
	} else {
		switch state {
		case firstFreeState:
			next := -1 * m.Check[state]
			m.Check[0] = -1 * next
			m.Base[next] = m.Base[state]
		case lastFreeState:
			prev := -1 * m.Base[state]
			m.Base[firstFreeState] = -1 * prev
			m.Check[prev] = -1
		default:
			next := -1 * m.Check[state]
			prev := -1 * m.Base[state]
			m.Check[prev] = -1 * next
			m.Base[next] = -1 * prev
		}
	}
	m.Check[state] = parentState
	m.Base[state] = LEAF
}

func (m *Matcher) setupfail(state, parentState, offset int) {
	failState := m.Fail[parentState]
	for {
		if m.hasEdge(failState, offset) {
			m.Fail[state] = m.Base[failState] + offset
			break
		}
		if failState == 0 {
			break
		}
		failState = m.Fail[failState]
	}
}

func (m *Matcher) unionFailOutput(state, failState int) {
	for _, word := range m.Output[failState] {
		m.Output[state] = append(m.Output[state], word)
	}
}

// Foofb TODO
func (m *Matcher) Foofb(edges []byte) int {
	if len(edges) == 0 {
		// TODO this should be removed eventually
		return LEAF
	}

	min := int(edges[0])
	max := int(edges[len(edges)-1])
	width := max - min
	freeState := m.firstFreeState()
	for freeState != -1 {
		valid := true
		for _, e := range edges[1:] {
			state := freeState + int(e) - min
			if state >= len(m.Check) {
				break
			} else if m.Check[state] >= 0 {
				valid = false
				break
			}
		}

		if valid {
			if freeState+width >= len(m.Check) {
				m.Foobar(width - len(m.Check) + freeState + 1)
			}
			return freeState - min
		}

		freeState = m.nextFreeState(freeState)
	}
	freeState = len(m.Check)
	m.Foobar(width + 1)
	return freeState - min
}

// findBase TODO
func (m *Matcher) findBase(edges []byte) int {
	if len(edges) == 0 {
		return LEAF
	}
	min := int(edges[0])
	max := int(edges[len(edges)-1])
	width := max - min

	if len(edges) < 1 {
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

// Foobar TODO
// increaseSize increases the size of base, check, and fail to ensure they
// remain the same size.
// It also sets the default value for these new unoccupied states which form
// bidirectional links to allow fast access to empty states. These
// bidirectional links only pertain to base and check.
//
// Example:
// m:
//   base:  [ 5 0 0 ]
//   check: [ 0 0 0 ]
// increaseSize(3):
//   base:  [ 5  0 0 -5 -3 -4 ]
//   check: [ -3 0 0 -4 -5 -1 ]
// increaseSize(3):
//   base:  [ 5  0 0 -8 -3 -4 -5 -6 -7]
//   check: [ -3 0 0 -4 -5 -6 -7 -8 -1]
//
// m:
//   base:  [ 5 0 0 ]
//   check: [ 0 0 0 ]
// increaseSize(1):
//   base:  [ 5  0 0 -3 ]
//   check: [ -3 0 0 -1 ]
// increaseSize(1):
//   base:  [ 5  0 0 -4 -3 ]
//   check: [ -3 0 0 -4 -1 ]
// increaseSize(1):
//   base:  [ 5  0 0 -5 -3 -4 ]
//   check: [ -3 0 0 -4 -5 -1 ]
func (m *Matcher) Foobar(dsize int) {
	if dsize == 0 {
		return
	}

	m.Base = append(m.Base, make([]int, dsize)...)
	m.Check = append(m.Check, make([]int, dsize)...)
	m.Fail = append(m.Fail, make([]int, dsize)...)

	lastFreeState := m.lastFreeState()
	firstFreeState := m.firstFreeState()
	for i := len(m.Check) - dsize; i < len(m.Check); i++ {
		if lastFreeState == -1 {
			m.Check[0] = -1 * i
			m.Base[i] = -1 * i
			m.Check[i] = -1
			firstFreeState = i
			lastFreeState = i
		} else {
			m.Base[i] = -1 * lastFreeState
			m.Check[i] = -1
			m.Base[firstFreeState] = -1 * i
			m.Check[lastFreeState] = -1 * i
			lastFreeState = i
		}
	}
}

func (m *Matcher) nextFreeState(curFreeState int) int {
	nextState := -1 * m.Check[curFreeState]

	// state 1 can never be a free state.
	if nextState == 1 {
		return -1
	}

	return nextState
}

func (m *Matcher) firstFreeState() int {
	state := m.Check[0]
	if state != 0 {
		return -1 * state
	}
	return -1
}

func (m *Matcher) lastFreeState() int {
	firstFree := m.firstFreeState()
	if firstFree != -1 {
		return -1 * m.Base[firstFree]
	}
	return -1
}

func (m *Matcher) increaseSize(dsize int) {
	m.Base = append(m.Base, make([]int, dsize)...)
	m.Check = append(m.Check, make([]int, dsize)...)
	m.Fail = append(m.Fail, make([]int, dsize)...)
}

func (m *Matcher) hasEdge(fromState, offset int) bool {
	toState := m.Base[fromState] + offset
	return toState > 0 && toState < len(m.Check) && m.Check[toState] == fromState
	// return toState >= 0 && toState < len(m.Check) && m.Check[toState] == fromState+1
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
		for !m.hasEdge(state, offset) && state != 0 {
			state = m.Fail[state]
		}

		if m.hasEdge(state, offset) {
			state = m.Base[state] + offset
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
