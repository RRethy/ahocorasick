package ahocorasick

import (
	"fmt"
	"reflect"
	"testing"
)

func convert(got []*Match) []Match {
	var converted []Match
	for _, matchptr := range got {
		converted = append(converted, *matchptr)
	}
	return converted
}

func TestFindAllByteSlice(t *testing.T) {
	m := compile([][]byte{
		[]byte("he"),
		[]byte("his"),
		[]byte("hers"),
		[]byte("she")},
	)
	m.findAll([]byte("ushers")) // => { "she" 1 }, { "he" 2}, { "hers" 2 }
	tests := []struct {
		patterns [][]byte
		expected []Match
		text     []byte
	}{
		{
			[][]byte{[]byte("na"), []byte("ink"), []byte("ki")},
			[]Match{{[]byte("ink"), 0}, {[]byte("ki"), 2}},
			[]byte("inking"),
		},
		{
			[][]byte{[]byte("ca"), []byte("erica"), []byte("rice")},
			[]Match{{[]byte("ca"), 3}, {[]byte("erica"), 0}},
			[]byte("erican"),
		},
		{
			[][]byte{[]byte("he"), []byte("she"), []byte("his"), []byte("hers")},
			[]Match{{[]byte("he"), 2}, {[]byte("she"), 1}, {[]byte("hers"), 2}},
			[]byte("ushers"),
		},
		{
			[][]byte{[]byte("they"), []byte("their"), []byte("theyre"), []byte("the"), []byte("tea"), []byte("te"), []byte("team"), []byte("go"), []byte("goo"), []byte("good"), []byte("oode")},
			[]Match{{[]byte("the"), 0}, {[]byte("they"), 0}, {[]byte("theyre"), 0}, {[]byte("go"), 13}, {[]byte("goo"), 13}, {[]byte("good"), 13}, {[]byte("oode"), 14}, {[]byte("te"), 19}, {[]byte("tea"), 19}, {[]byte("team"), 19}},
			[]byte("theyre not a goode team"),
		},
		{
			[][]byte{[]byte("a")},
			[]Match{{[]byte("a"), 0}, {[]byte("a"), 1}, {[]byte("a"), 2}, {[]byte("a"), 5}, {[]byte("a"), 7}, {[]byte("a"), 9}, {[]byte("a"), 11}},
			[]byte("aaabbabababa"),
		},
		{
			[][]byte{},
			[]Match{},
			[]byte("there is no patterns"),
		},
		{
			[][]byte{[]byte("锅"), []byte("持有人"), []byte("potholderz"), []byte("MF DOOM")},
			[]Match{{[]byte("potholderz"), 0}, {[]byte("MF DOOM"), 14}, {[]byte("锅"), 39}, {[]byte("持有人"), 43}},
			[]byte("potholderz by MF DOOM hot shit aw shit 锅 持有人"),
		},
	}
	for _, test := range tests {
		matcher := compile(test.patterns)
		got := matcher.findAll(test.text)
		gotConverted := convert(got)
		if !(len(got) == 0 && len(test.expected) == 0) &&
			!reflect.DeepEqual(gotConverted, test.expected) {
			t.Errorf(`
        Text:     %s
		Expected: %v
		Got:      %v
		`, test.text, test.expected, gotConverted)
		}
	}
}

func TestIncreaseSize(t *testing.T) {
	m := &Matcher{
		[]int{5, 0, 0},
		[]int{0, 0, 0},
		[]int{0, 0, 0},
		[][]int{},
	}
	m.increaseSize(1)
	if !reflect.DeepEqual(m.base, []int{5, 0, 0, -3}) {
		t.Errorf("Got: %v\n", m.base)
	}
	if !reflect.DeepEqual(m.check, []int{-3, 0, 0, -1}) {
		t.Errorf("Got: %v\n", m.check)
	}

	m.increaseSize(1)
	if !reflect.DeepEqual(m.base, []int{5, 0, 0, -4, -3}) {
		t.Errorf("Got: %v\n", m.base)
	}
	if !reflect.DeepEqual(m.check, []int{-3, 0, 0, -4, -1}) {
		t.Errorf("Got: %v\n", m.check)
	}

	m.increaseSize(1)
	if !reflect.DeepEqual(m.base, []int{5, 0, 0, -5, -3, -4}) {
		t.Errorf("Got: %v\n", m.base)
	}
	if !reflect.DeepEqual(m.check, []int{-3, 0, 0, -4, -5, -1}) {
		t.Errorf("Got: %v\n", m.check)
	}

	m = &Matcher{
		[]int{5, 0, 0},
		[]int{0, 0, 0},
		[]int{0, 0, 0},
		[][]int{},
	}
	m.increaseSize(3)
	if !reflect.DeepEqual(m.base, []int{5, 0, 0, -5, -3, -4}) {
		t.Errorf("Got: %v\n", m.base)
	}
	if !reflect.DeepEqual(m.check, []int{-3, 0, 0, -4, -5, -1}) {
		t.Errorf("Got: %v\n", m.check)
	}

	m.increaseSize(3)
	if !reflect.DeepEqual(m.base, []int{5, 0, 0, -8, -3, -4, -5, -6, -7}) {
		t.Errorf("Got: %v\n", m.base)
	}
	if !reflect.DeepEqual(m.check, []int{-3, 0, 0, -4, -5, -6, -7, -8, -1}) {
		t.Errorf("Got: %v\n", m.check)
	}

	m = &Matcher{
		[]int{0},
		[]int{0},
		[]int{0},
		[][]int{},
	}
	m.increaseSize(5)
	if !reflect.DeepEqual(m.base, []int{0, -5, -1, -2, -3, -4}) {
		t.Errorf("Got: %v\n", m.base)
	}
	if !reflect.DeepEqual(m.check, []int{-1, -2, -3, -4, -5, -1}) {
		t.Errorf("Got: %v\n", m.check)
	}

	m = &Matcher{
		[]int{-103, -1867},
		[]int{0, 0},
		[]int{},
		[][]int{},
	}
	m.increaseSize(5)
	if !reflect.DeepEqual(m.base, []int{-103, -1867, -6, -2, -3, -4, -5}) {
		t.Errorf("Got: %v\n", m.base)
	}
	if !reflect.DeepEqual(m.check, []int{-2, 0, -3, -4, -5, -6, -1}) {
		t.Errorf("Got: %v\n", m.check)
	}
}

func TestNextFreeState(t *testing.T) {
	m := &Matcher{
		[]int{5, 0, 0, -3},
		[]int{-3, 0, 0, -1},
		[]int{},
		[][]int{},
	}
	nextState := m.nextFreeState(3)
	if nextState != -1 {
		t.Errorf("Got: %d\n", nextState)
	}

	m.increaseSize(3)
	nextState = m.nextFreeState(3)
	if nextState != 4 {
		t.Errorf("Got: %d\n", nextState)
	}
}

func TestOccupyState(t *testing.T) {
	m := &Matcher{
		[]int{5, 0, 0, -3},
		[]int{-3, 0, 0, -1},
		[]int{},
		[][]int{},
	}
	m.increaseSize(5)
	m.occupyState(3, 1)
	m.occupyState(4, 1)
	m.occupyState(8, 1)
	m.occupyState(6, 1)
	m.occupyState(5, 1)
	m.occupyState(7, 1)
	if !reflect.DeepEqual(m.base, []int{5, 0, 0, -1867, -1867, -1867, -1867, -1867, -1867}) {
		t.Errorf("Got: %v\n", m.base)
	}
	if !reflect.DeepEqual(m.check, []int{0, 0, 0, 1, 1, 1, 1, 1, 1}) {
		t.Errorf("Got: %v\n", m.check)
	}
}

func ExampleMatcher_FindAllByteSlice() {
	matcher := CompileByteSlices([][]byte{
		[]byte("he"),
		[]byte("she"),
		[]byte("his"),
		[]byte("hers"),
		[]byte("she"),
	})
	fmt.Print(matcher.FindAllByteSlice([]byte("ushers")))

	// Output:
	// [{ "he" 2 } { "she" 1 } { "she" 1 } { "hers" 2 }]
}

func ExampleMatcher_FindAllString() {
	matcher := CompileStrings([]string{
		"he",
		"she",
		"his",
		"hers",
		"she",
	})
	fmt.Print(matcher.FindAllString("ushers"))

	// Output:
	// [{ "he" 2 } { "she" 1 } { "she" 1 } { "hers" 2 }]
}
