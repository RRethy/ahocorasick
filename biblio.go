package biblio

import (
	"fmt"
)

const (
	// LEAF represents a leaf on the trie
	// This must be <255 since the offsets used are in [0,255]
	// This should only appear in the Base array since the Check array uses
	// negative values to represent free states.
	LEAF = -1867
)

// Matcher is the pattern matching state machine.
type Matcher struct {
	Base   []int   // base array in the double array trie
	Check  []int   // check array in the double array trie
	Fail   []int   // fail function
	Output [][]int // output function

}

func (m *Matcher) String() string {
	return fmt.Sprintf(`
Base:   %v
Check:  %v
Fail:   %v
Output: %v
`, m.Base, m.Check, m.Fail, m.Output)
}

// CompileByteSlices compiles a Matcher from a slice of byte slices. This Matcher can be
// used to find occurrences of each pattern in a text.
func CompileByteSlices(words [][]byte) *Matcher {
	m := new(Matcher)
	m.Base = make([]int, 2048)[:1]
	m.Check = make([]int, 2048)[:1]
	m.Fail = make([]int, 2048)[:1]
	m.Output = make([][]int, 2048)[:1]

	type wordTuple struct {
		i      int
		length int
	}

	// Represents a node in the implicit trie of words
	type trienode struct {
		state      int
		depth      int
		wordTuples []wordTuple
	}
	queue := make([]trienode, 256)[:1]
	wordTuples := make([]wordTuple, len(words)) // TODO [:]
	for i, word := range words {
		wordTuples[i] = wordTuple{i, len(word)}
	}
	queue[0] = trienode{0, 0, wordTuples}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if len(node.wordTuples) == 0 {
			m.Base[node.state] = LEAF
			continue
		}

		// TODO this is kinda confusing code
		var wordTuples [256][]wordTuple
		for _, tuple := range node.wordTuples {
			edge := words[tuple.i][node.depth]
			wordTuples[edge] = append(wordTuples[edge], tuple)
		}

		var edges []byte
		for edge, tuples := range wordTuples {
			if len(tuples) > 0 {
				edges = append(edges, byte(edge))
			}
		}

		// Calculate a suitable Base value where each edge will fit into the
		// double array trie
		base := m.findBase(edges)
		m.Base[node.state] = base

		for _, edge := range edges {
			offset := int(edge)
			newState := base + offset

			m.occupyState(newState, node.state)

			// level 0 and level 1 should fail to state 0
			if node.depth > 0 {
				m.setFailState(newState, node.state, offset)
			}
			m.unionFailOutput(newState, m.Fail[newState])

			// Add the child nodes to the queue to continue down the BFS
			newnode := trienode{newState, node.depth + 1, []wordTuple{}}
			for _, tuple := range wordTuples[edge] {
				if newnode.depth < tuple.length {
					newnode.wordTuples = append(newnode.wordTuples, tuple)
				} else {
					m.Output[newState] = append(m.Output[newState], tuple.length)
				}
			}
			queue = append(queue, newnode)
		}
	}

	return m
}

// CompileStrings compiles a Matcher from a slice of strings. This Matcher can
// be used to find occurrences of each pattern in a text.
func CompileStrings(words []string) *Matcher {
	var wordByteSlices [][]byte
	for _, word := range words {
		wordByteSlices = append(wordByteSlices, []byte(word))
	}
	return CompileByteSlices(wordByteSlices)
}

// occupyState will correctly occupy state so it maintains the
// index=check[base[index]+offset] identity. It will also update the
// bidirectional link of free states correctly.
// Note: This MUST be used instead of simply modifying the check array directly
// which is break the bidirectional link of free states.
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

// setFailState sets the output of the fail function for input state. It will
// traverse up the fail states of it's ancestors until it reaches a fail state
// with a transition for offset.
func (m *Matcher) setFailState(state, parentState, offset int) {
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

// unionFailOutput unions the output function for failState with the output
// function for state and sets the result as the output function for state.
// This allows us to match substrings, commenting out this body would match
// every word that is not a substring.
func (m *Matcher) unionFailOutput(state, failState int) {
	m.Output[state] = append([]int{}, m.Output[failState]...)
}

// findBase finds a base value which has free states in the positions that
// correspond to each edge transition in edges. If this does not exist, then
// base and check (and the fail array for consistency) will be extended just
// enough to fit each transition.
// The extension will maintain the bidirectional link of free states.
func (m *Matcher) findBase(edges []byte) int {
	if len(edges) == 0 {
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
				m.increaseSize(width - len(m.Check) + freeState + 1)
			}
			return freeState - min
		}

		freeState = m.nextFreeState(freeState)
	}
	freeState = len(m.Check)
	m.increaseSize(width + 1)
	return freeState - min
}

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
func (m *Matcher) increaseSize(dsize int) {
	if dsize == 0 {
		return
	}

	m.Base = append(m.Base, make([]int, dsize)...)
	m.Check = append(m.Check, make([]int, dsize)...)
	m.Fail = append(m.Fail, make([]int, dsize)...)
	m.Output = append(m.Output, make([][]int, dsize)...)

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

// nextFreeState uses the nature of the bidirectional link to determine the
// closest free state at a larger index. Since the check array holds the
// negative index of the next free state, except for the last free state which
// has a value of -1, negating this value is the next free state.
func (m *Matcher) nextFreeState(curFreeState int) int {
	nextState := -1 * m.Check[curFreeState]

	// state 1 can never be a free state.
	if nextState == 1 {
		return -1
	}

	return nextState
}

// firstFreeState uses the first value in the check array which points to the
// first free state. A value of 0 means there are no free states and -1 is
// returned.
func (m *Matcher) firstFreeState() int {
	state := m.Check[0]
	if state != 0 {
		return -1 * state
	}
	return -1
}

// lastFreeState uses the base value of the first free state which points the
// last free state.
func (m *Matcher) lastFreeState() int {
	firstFree := m.firstFreeState()
	if firstFree != -1 {
		return -1 * m.Base[firstFree]
	}
	return -1
}

// hasEdge determines if the fromState has a transition for offset.
func (m *Matcher) hasEdge(fromState, offset int) bool {
	toState := m.Base[fromState] + offset
	return toState > 0 && toState < len(m.Check) && m.Check[toState] == fromState
}

// Match represents a matched pattern in the text
type Match struct {
	Word  []byte // the matched pattern
	Index int    // the start index of the match
}

// FindAllByteSlice finds all instances of the patterns in the text.
func (m *Matcher) FindAllByteSlice(text []byte) (matches []*Match) {
	state := 0
	for i, b := range text {
		offset := int(b)
		for state != 0 && !m.hasEdge(state, offset) {
			state = m.Fail[state]
		}

		if m.hasEdge(state, offset) {
			state = m.Base[state] + offset
		}
		for _, wordlen := range m.Output[state] {
			matches = append(matches, &Match{text[i-wordlen+1 : i+1], i - wordlen + 1})
		}
	}
	return
}

// FindAllString finds all instances of the patterns in the text.
func (m *Matcher) FindAllString(text string) []*Match {
	return m.FindAllByteSlice([]byte(text))
}
