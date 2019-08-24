// Copyright 2019 Adam P. Regasz-Rethy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package biblio implements the aho-corasick pattern matching algorithm.  For
// more information on the algorithm, see the academic paper published by
// Alfred V. Aho and Margaret J. Corasick which is freely available at
// https://cr.yp.to/bib/1975/aho.pdf.
//
// Terminology used corresponds to the terminology used in the aforementioned
// paper.
// This package uses uses a modified version of algorithm 2 to construct the
// goto function, trie.tree, and partially constructs the output function,
// Biblio.output. Then it uses a modified version of algorithm 4 to construct
// the next function, Biblio.next, which lets us ignore algorithm 3 and the
// failure function. This has a negative effect on compile time which is
// benchmarked and can be viewed by running `git log --grep=BenchmarkCompile`.
// The benefit of using algorithm 4 and the next function is the faster parsing
// since we don't need to do any failure transitions, this has benchmarks
// available for viewing with `git log --grep=BenchmarkParse`.
//
// UTF-8 is fully supported since all transitions are represented using the
// rune type.
//
// There are two public facing datatypes, biblio.Biblio and biblio.Match. To
// construct a Biblio, one must first call biblio.Compile and pass it a slice
// of strings which are the patterns that will be looked for. Under the hood,
// biblio.Compile will construct a finite state machine. Then, a string can be
// passed to biblio.FindAll, which will then return a slice of biblio.Matches.
// For example:
// bib := biblio.Compile([]string{"foo", "bar", "baz"})
// bib.FindAll("foo bar baz bot") => {"foo", 0, 3}, {"bar", 4, 3}, {"baz", 8, 3}
//
// Let Σ be the alphabet for the patterns passed to biblio.Compile such that
//   the alphabet represents the set of characters used in the patterns.
// Let total_pats_len be the sum of the lengths of all patterns passed to
//   biblio.Compile
// Let max be the length of the longest pattern passed to biblio.Compile
// let n be the number of patterns passed to biblio.Compile
// Let m be the number of matches
//
// Time Complexity:
// biblio.Compile: O(|Σ|*total_pats_len)
// biblio.Biblio.FindAll(text): Θ(len(text))
//
// Space Complexity:
// biblio.Compile: O(|Σ|*n*max)
// biblio.Biblio.FindAll: Θ(m)
package biblio

// Biblio is the representation of a compiled set of patterns.
type Biblio struct {
	next   map[int]map[rune]int    // state => { transition character => state }
	output map[int]map[string]bool // state => set of words which terminate at state
}

// trie is a simple trie representation used to construct the required dfa
type trie struct {
	tree map[int]map[rune]int
}

// Match is a representation of a pattern found in the text.
// Word is the pattern that was matched
// Index is the index in the text of the first character of the matched pattern
type Match struct {
	Word  string
	Index int
}

// Compile creates a Biblio for use in parsing. A state machine will be created
// which can be used for linear time parsing of a text.
func Compile(words []string) *Biblio {
	biblio := new(Biblio)
	if len(words) == 0 {
		return biblio
	}

	// create the trie from each word in words
	t := trie{}
	t.tree = map[int]map[rune]int{}
	biblio.output = map[int]map[string]bool{}
	for _, word := range words {
		t.add(word, &biblio.output)
	}

	// t.tree is a subgraph of biblio.next. It is used as the starting graph
	// for it which then as additional edges added in buildNextFunc.
	biblio.next = map[int]map[rune]int{}
	for state, transition := range t.tree {
		biblio.next[state] = map[rune]int{}
		for c, dest := range transition {
			biblio.next[state][c] = dest
		}
	}

	biblio.buildNextFunc(&t)

	return biblio
}

// FindAll returns a slice of biblio.Match which represent each pattern found in
// text
func (biblio *Biblio) FindAll(text string) (matches []Match) {
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
			matches = append(matches, Match{word, i - len(word) + 1})
		}
	}
	return
}

// add adds a word to the trie t and modifies output which holds information on
// which state is a word terminating state in the trie
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

// builds the finited state machine, biblio.next
func (biblio *Biblio) buildNextFunc(t *trie) {
	type statedata struct {
		state           int
		transition      rune
		parentFailState int
		level           int
	}

	// will bfs the trie
	queue := make([]statedata, 1)[:]
	for len(queue) > 0 {
		// pop front of queue
		data := queue[0]
		queue = queue[1:]

		failstate := 0
		// if data.state is in level 1 or 0, it has a failure state of 0 for
		// all transitions
		//
		// if not then check if the parent of the current state has failure
		// state which has the appropriate transition. This occurs for
		// substring matches. For example, if we were going down the trie
		// looking at "abc", then fail, then we could continue down the path in
		// the trie holding "bc"
		//
		// if not, then check if state 0 has the appropriate transition. This
		// occurs for substring matches starting at the current character. For
		// example, if we were going down the trie looking at "abc", then fail,
		// then we could continue down the path in the trie holding "c"
		//
		// if not, then 0 is the default fail state
		//
		// the trie in the previous examples is in fact the next dfa, but can
		// be treated as a trie for explanations since there is a valid trie
		// which is a subgraph of the next dfa
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

		// union the transitions in the fail state with the current state, if
		// they share a transition character, then the transition of the
		// current state takes precedence
		for c, state := range biblio.next[failstate] {
			if _, ok := biblio.next[data.state][c]; !ok {
				biblio.next[data.state][c] = state
			}
		}

		// continue adding to the queue since we are doing a bfs
		for c, state := range t.tree[data.state] {
			queue = append(queue, statedata{state, c, failstate, data.level + 1})
		}
	}
}
