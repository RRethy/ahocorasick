package main

import (
	"reflect"
	"testing"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		input  []string
		next   map[int]map[rune]int
		output map[int]map[string]bool
	}{
		{
			[]string{"theyre", "they", "the"},
			map[int]map[rune]int{
				0: {'t': 1},
				1: {'h': 2, 't': 1},
				2: {'e': 3, 't': 1},
				3: {'y': 4, 't': 1},
				4: {'r': 5, 't': 1},
				5: {'e': 6, 't': 1},
				6: {'t': 1},
			},
			map[int]map[string]bool{
				3: map[string]bool{"the": true},
				4: map[string]bool{"they": true},
				6: map[string]bool{"theyre": true},
			},
		},
		{
			[]string{"he", "she", "his", "hers"},
			map[int]map[rune]int{
				0: {'h': 1, 's': 3},
				1: {'h': 1, 'e': 2, 's': 3, 'i': 6},
				2: {'h': 1, 's': 3, 'r': 8},
				3: {'s': 3, 'h': 4},
				4: {'h': 1, 's': 3, 'e': 5, 'i': 6},
				5: {'h': 1, 's': 3, 'r': 8},
				6: {'h': 1, 's': 7},
				7: {'h': 4, 's': 3},
				8: {'h': 1, 's': 9},
				9: {'h': 4, 's': 3},
			},
			map[int]map[string]bool{
				2: map[string]bool{"he": true},
				5: map[string]bool{"she": true, "he": true},
				7: map[string]bool{"his": true},
				9: map[string]bool{"hers": true},
			},
		},
		{
			[]string{},
			map[int]map[rune]int{},
			map[int]map[string]bool{},
		},
		{
			[]string{"h"},
			map[int]map[rune]int{
				0: {'h': 1},
				1: {'h': 1},
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

func TestParse(t *testing.T) {
	tests := []struct {
		patterns []string
		expected []Match
		text     string
	}{
		{
			[]string{"he", "she", "his", "hers"},
			[]Match{{"she", 3}, {"he", 3}, {"hers", 5}},
			"ushers",
		},
		{
			[]string{"they", "their", "theyre", "the", "tea", "te", "team", "go", "goo", "good"},
			[]Match{{"the", 2}, {"they", 3}, {"theyre", 5}, {"go", 14}, {"goo", 15}, {"good", 16}, {"te", 19}, {"tea", 20}, {"team", 21}},
			"theyre not a good team",
		},
	}
	for _, test := range tests {
		bib := Compile(test.patterns)
		got := bib.Parse(test.text)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf(`
Expected: %v
Got:      %v
`, test.expected, got)
		}
	}
}

func BenchmarkCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Compile([]string{
			"he",
			"she",
			"they",
			"their",
			"where",
			"bear",
			"taratula",
			"adam",
			"regard-rethy",
			"panda",
			"bear",
			"golang",
			"his",
			"hers",
			"her",
		})
	}
}

func BenchmarkParse(b *testing.B) {
	bib := Compile([]string{
		"he",
		"she",
		"they",
		"their",
		"where",
		"bear",
		"taratula",
		"adam",
		"regard-rethy",
		"panda",
		"bear",
		"golang",
		"his",
		"hers",
		"her",
	})
	for i := 0; i < b.N; i++ {
		bib.Parse(`
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
	`)
	}
}
