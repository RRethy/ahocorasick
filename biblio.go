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
	next   map[int]map[rune]int    // TODO
	output map[int]map[string]bool // state => set of words which terminate at state
}

type trie struct {
	tree map[int]map[rune]int
}

// Match TODO
type Match struct {
	word  string
	index int
}

// Parse TODO
func (biblio *Biblio) Parse(text string) (matches []Match) {
	if len(biblio.output) == 0 {
		return
	}

	state := 0
	for i, c := range text {
		if next, ok := biblio.next[state][c]; ok {
			state = next
		} else {
			state = 0
		}
		for word := range biblio.output[state] {
			matches = append(matches, Match{word, i})
		}
	}
	return
}

func (t *trie) add(word string, output *map[int]map[string]bool) {
	state := 0
	if len(t.tree) == 0 {
		t.tree[0] = map[rune]int{}
	}

	for _, c := range word {
		if next, ok := t.tree[state][c]; ok {
			state = next
		} else {
			next = len(t.tree)
			t.tree[state][c] = next
			t.tree[next] = map[rune]int{}
			state = next
		}
	}

	(*output)[state] = map[string]bool{word: true}
}

// builds the failure function for each state
func (biblio *Biblio) buildFailureTransitions(t *trie) {
	type statedata struct {
		state           int
		transition      rune
		parentFailState int
		level           int
	}

	queue := make([]statedata, 1)[:]
	for len(queue) > 0 {
		data := queue[0]
		queue = queue[1:]

		failstate := 0
		if data.level > 1 {
			if state, ok := biblio.next[data.parentFailState][data.transition]; ok {
				for word := range biblio.output[state] {
					if _, ok := biblio.output[data.state]; !ok {
						biblio.output[data.state] = map[string]bool{word: true}
					} else {
						biblio.output[data.state][word] = true
					}
				}
				failstate = state
			} else if state, ok := biblio.next[0][data.transition]; ok {
				failstate = state
			}
		}

		for c, state := range biblio.next[failstate] {
			if _, ok := biblio.next[data.state][c]; !ok {
				biblio.next[data.state][c] = state
			}
		}

		for c, state := range t.tree[data.state] {
			queue = append(queue, statedata{state, c, failstate, data.level + 1})
		}
	}
}

// Compile TODO
func Compile(words []string) *Biblio {
	biblio := new(Biblio)
	if len(words) == 0 {
		return biblio
	}

	t := trie{}
	t.tree = map[int]map[rune]int{}
	biblio.output = map[int]map[string]bool{}
	for _, word := range words {
		t.add(word, &biblio.output)
	}

	biblio.next = map[int]map[rune]int{}
	for state, transition := range t.tree {
		biblio.next[state] = map[rune]int{}
		for c, dest := range transition {
			biblio.next[state][c] = dest
		}
	}

	biblio.buildFailureTransitions(&t)

	return biblio
}
