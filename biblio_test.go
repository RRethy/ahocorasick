package main

import (
	"reflect"
	"testing"
)

func TestCompile(t *testing.T) {
	tests := []struct {
		input  []string
		g      map[int]map[rune]int
		f      map[int]int
		output map[int]map[string]bool
	}{
		{
			[]string{"he", "she", "his", "hers"},
			map[int]map[rune]int{
				0: {'h': 1, 's': 3},
				1: {'e': 2, 'i': 6},
				2: {'r': 8},
				3: {'h': 4},
				4: {'e': 5},
				5: {},
				6: {'s': 7},
				7: {},
				8: {'s': 9},
				9: {},
			},
			map[int]int{
				0: 0,
				1: 0,
				2: 0,
				3: 0,
				4: 1,
				5: 2,
				6: 0,
				7: 3,
				8: 0,
				9: 3,
			},
			map[int]map[string]bool{
				2: map[string]bool{"he": true},
				5: map[string]bool{"she": true, "he": true},
				7: map[string]bool{"his": true},
				9: map[string]bool{"hers": true},
			},
		},
	}
	for _, test := range tests {
		got := Compile(test.input)

		// check for correctly compiled state machine
		if !reflect.DeepEqual(test.g, got.g) {
			t.Errorf(`
Expected: %v
Got:      %v`, test.g, got.g)
		}
		if !reflect.DeepEqual(test.f, got.f) {
			t.Errorf(`
Expected: %v
Got:      %v`, test.f, got.f)
		}
		if !reflect.DeepEqual(test.output, got.output) {
			t.Errorf(`
Expected: %v
Got:      %v`, test.output, got.output)
		}
	}
}

func TestParse(t *testing.T) {
}
