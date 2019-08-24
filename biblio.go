package main

// For more information on the aho-corasick pattern matching algorithm for
// bibliographic search, see the following paper:
// https://cr.yp.to/bib/1975/aho.pdf

// Biblio is the container for the neccessary components for the aho-corasick
// string matching algorithm for bibliographic search
// g is the goto function
// f is the failure function
// output is the output function
type Biblio struct {
	g      map[int]map[rune]int    // state => transition => next state
	f      map[int]int             // state => fail state
	output map[int]map[string]bool // state => set of words which terminate at state
}

// Match TODO
type Match struct {
	word  string
	index int
}

// Parse TODO
func (biblio *Biblio) Parse(text string) (matches []Match) {
	state := 0
	for i, c := range text {
		for {
			if _, ok := biblio.g[state][c]; !ok {
				if state == 0 {
					break
				} else {
					state = biblio.f[state]
				}
			} else {
				state = biblio.g[state][c]
				break
			}
		}
		for word := range biblio.output[state] {
			matches = append(matches, Match{word, i})
		}
	}
	return
}

// adds word to the trie represented by biblio.g
func (biblio *Biblio) addWord(word string) {
	state := 0
	for _, c := range word {
		if _, ok := biblio.g[state]; !ok {
			biblio.g[state] = map[rune]int{}
		}
		if _, ok := biblio.g[state][c]; !ok {
			biblio.g[state][c] = len(biblio.g)
		}
		state = biblio.g[state][c]
	}
	biblio.g[state] = map[rune]int{}
	// denote that state terminates word
	biblio.output[state] = map[string]bool{word: true}
}

// builds the failure function for each state
func (biblio *Biblio) buildFailureTransitions() {
	type failtuple struct {
		state           int
		transition      rune
		parentFailState int
		level           int
	}
	queue := make([]failtuple, 1)[:] // Note: this implicitly adds root state to the queue
	// bfs the trie
	for len(queue) > 0 {
		// pop front of queue
		tuple := queue[0]
		queue = queue[1:]

		// Only nodes lower than level 1 can have a non-zero fail state
		if tuple.level > 1 {
			// Use the following algorithm to determine fail state:
			// 1. check if parent fail state has current transition, if so then
			//    transition to next state which is fail state
			// 2. check if root has current transition, if so then transition
			//    to next state which is fail state
			// 3. else root is fail state
			if failState, ok := biblio.g[tuple.parentFailState][tuple.transition]; ok {
				// the set output(failState) to the set output(current state)
				// since they are substring matches
				for word := range biblio.output[failState] {
					if _, ok := biblio.output[tuple.state]; !ok {
						biblio.output[tuple.state] = map[string]bool{}
					}
					biblio.output[tuple.state][word] = true
				}
				biblio.f[tuple.state] = failState
			} else if failState, ok := biblio.g[0][tuple.transition]; ok {
				biblio.f[tuple.state] = failState
			} else {
				biblio.f[tuple.state] = 0
			}
		} else {
			biblio.f[tuple.state] = 0
		}

		// continue down bfs traversal
		for c, childState := range biblio.g[tuple.state] {
			queue = append(queue, failtuple{childState, c, biblio.f[tuple.state], tuple.level + 1})
		}
	}
}

// Compile TODO
func Compile(words []string) *Biblio {
	biblio := new(Biblio)
	if len(words) == 0 {
		return biblio
	}

	biblio.g = map[int]map[rune]int{}
	biblio.f = map[int]int{}
	biblio.output = map[int]map[string]bool{}
	for _, word := range words {
		biblio.addWord(word)
	}
	biblio.buildFailureTransitions()

	return biblio
}
