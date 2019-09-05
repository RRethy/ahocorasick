package main

import (
	"reflect"
	"testing"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		input  [][]byte
		next   []map[byte]int
		output map[int]map[string]bool
	}{
		{
			[][]byte{[]byte("theyre"), []byte("they"), []byte("the")},
			[]map[byte]int{
				{'t': 1},
				{'h': 2, 't': 1},
				{'e': 3, 't': 1},
				{'y': 4, 't': 1},
				{'r': 5, 't': 1},
				{'e': 6, 't': 1},
				{'t': 1},
			},
			map[int]map[string]bool{
				3: map[string]bool{"the": true},
				4: map[string]bool{"they": true},
				6: map[string]bool{"theyre": true},
			},
		},
		{
			[][]byte{[]byte("he"), []byte("she"), []byte("his"), []byte("hers")},
			[]map[byte]int{
				{'h': 1, 's': 3},
				{'h': 1, 'e': 2, 's': 3, 'i': 6},
				{'h': 1, 's': 3, 'r': 8},
				{'s': 3, 'h': 4},
				{'h': 1, 's': 3, 'e': 5, 'i': 6},
				{'h': 1, 's': 3, 'r': 8},
				{'h': 1, 's': 7},
				{'h': 4, 's': 3},
				{'h': 1, 's': 9},
				{'h': 4, 's': 3},
			},
			map[int]map[string]bool{
				2: map[string]bool{"he": true},
				5: map[string]bool{"she": true, "he": true},
				7: map[string]bool{"his": true},
				9: map[string]bool{"hers": true},
			},
		},
		{
			[][]byte{},
			[]map[byte]int{},
			map[int]map[string]bool{},
		},
		{
			[][]byte{[]byte("h")},
			[]map[byte]int{
				{'h': 1},
				{'h': 1},
			},
			map[int]map[string]bool{
				1: map[string]bool{"h": true},
			},
		},
	}
	for _, test := range tests {
		got := Compile(test.input)

		// check for correctly compiled state machine
		if !(len(test.next) == 0 && len(got.next) == 0) && !reflect.DeepEqual(test.next, got.next) {
			t.Errorf(`
Expected: %v
Got:      %v`, test.next, got.next)
		}
		if !(len(test.output) == 0 && len(got.output) == 0) && !reflect.DeepEqual(test.output, got.output) {
			t.Errorf(`
Expected: %v
Got:      %v`, test.output, got.output)
		}
	}
}

func TestFindAll(t *testing.T) {
	tests := []struct {
		patterns [][]byte
		expected []Match
		text     []byte
	}{
		{
			[][]byte{[]byte("he"), []byte("she"), []byte("his"), []byte("hers")},
			[]Match{{"she", 1}, {"he", 2}, {"hers", 2}},
			[]byte("ushers"),
		},
		{
			[][]byte{[]byte("they"), []byte("their"), []byte("theyre"), []byte("the"), []byte("tea"), []byte("te"), []byte("team"), []byte("go"), []byte("goo"), []byte("good"), []byte("oode")},
			[]Match{{"the", 0}, {"they", 0}, {"theyre", 0}, {"go", 13}, {"goo", 13}, {"good", 13}, {"oode", 14}, {"te", 19}, {"tea", 19}, {"team", 19}},
			[]byte("theyre not a goode team"),
		},
		{
			[][]byte{[]byte("a")},
			[]Match{{"a", 0}, {"a", 1}, {"a", 2}, {"a", 5}, {"a", 7}, {"a", 9}, {"a", 11}},
			[]byte("aaabbabababa"),
		},
		{
			[][]byte{},
			[]Match{},
			[]byte("there is no patterns"),
		},
		{
			[][]byte{[]byte("锅"), []byte("持有人"), []byte("potholderz"), []byte("MF DOOM")},
			[]Match{{"potholderz", 0}, {"MF DOOM", 14}, {"锅", 39}, {"持有人", 43}},
			[]byte("potholderz by MF DOOM hot shit aw shit 锅 持有人"),
		},
	}
	for _, test := range tests {
		bib := Compile(test.patterns)
		got := bib.FindAll(test.text)
		if !(len(got) == 0 && len(test.expected) == 0) && !reflect.DeepEqual(got, test.expected) {
			t.Errorf(`
Expected: %v
Got:      %v
`, test.expected, got)
		}
	}
}

func BenchmarkCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Compile([][]byte{
			[]byte("he"),
			[]byte("she"),
			[]byte("they"),
			[]byte("their"),
			[]byte("where"),
			[]byte("bear"),
			[]byte("taratula"),
			[]byte("adam"),
			[]byte("regard-rethy"),
			[]byte("panda"),
			[]byte("bear"),
			[]byte("golang"),
			[]byte("his"),
			[]byte("hers"),
			[]byte("her"),
		})
	}
}

func BenchmarkFindAll(b *testing.B) {
	bib := Compile([][]byte{
		[]byte("he"),
		[]byte("she"),
		[]byte("they"),
		[]byte("their"),
		[]byte("where"),
		[]byte("bear"),
		[]byte("taratula"),
		[]byte("adam"),
		[]byte("regard-rethy"),
		[]byte("panda"),
		[]byte("bear"),
		[]byte("golang"),
		[]byte("his"),
		[]byte("hers"),
		[]byte("her"),
	})
	for i := 0; i < b.N; i++ {
		bib.FindAll([]byte(`
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
ushers golang to     be rrrrrrrr tartula taratulapandawhere
	`))
	}
}
