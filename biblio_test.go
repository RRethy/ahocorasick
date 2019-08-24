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
