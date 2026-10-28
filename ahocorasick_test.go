package ahocorasick

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"slices"
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
	m.FindByteSliceN([]byte("ushers"), -1) // => { "she" 1 }, { "he" 2}, { "hers" 2 }
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
	for i, test := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			matcher := compile(test.patterns)
			got := matcher.FindByteSliceN(test.text, -1)
			gotConverted := convert(got)
			if !(len(got) == 0 && len(test.expected) == 0) &&
				!reflect.DeepEqual(gotConverted, test.expected) {
				t.Errorf(`
        Text:     %s
		Expected: %v
		Got:      %v
		`, test.text, test.expected, gotConverted)
			}
		})
	}
}

func TestFindNStreaming(t *testing.T) {
	matcher := CompileStrings([]string{"he", "she", "his", "hers"})
	chunks := []string{"us", "h", "ers"}
	consumed := 0
	text := func(yield func(byte) bool) {
		for _, chunk := range chunks {
			for i := range len(chunk) {
				consumed++
				if !yield(chunk[i]) {
					return
				}
			}
		}
	}

	got := convert(matcher.FindN(text, 2))
	want := []Match{{[]byte("he"), 2}, {[]byte("she"), 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindN() = %v, want %v", got, want)
	}
	if consumed != 4 {
		t.Errorf("consumed %d bytes, want 4", consumed)
	}
}

func TestFindReadNStopsAtLimit(t *testing.T) {
	matcher := CompileStrings([]string{"he", "she", "his", "hers"})
	reader := bytes.NewReader([]byte("ushers"))

	got, err := matcher.FindReadN(reader, 2)
	want := []Match{{[]byte("he"), 2}, {[]byte("she"), 1}}
	if err != nil {
		t.Fatalf("FindReadN() error = %v", err)
	}
	if got := convert(got); !reflect.DeepEqual(got, want) {
		t.Errorf("FindReadN() = %v, want %v", got, want)
	}
	if reader.Len() != 2 {
		t.Errorf("reader has %d bytes remaining, want 2", reader.Len())
	}
}

func TestFindReadNPropagatesError(t *testing.T) {
	matcher := CompileStrings([]string{"he", "she", "his", "hers"})
	readErr := errors.New("read failed")
	text := []byte("ushers")
	reader := &errorByteReader{data: text, err: readErr}

	got, err := matcher.FindReadN(reader, -1)
	if !errors.Is(err, readErr) {
		t.Fatalf("FindReadN() error = %v, want %v", err, readErr)
	}
	if want := matcher.FindByteSliceN(text, -1); !reflect.DeepEqual(got, want) {
		t.Errorf("FindReadN() matches = %v, want %v", got, want)
	}
}

type errorByteReader struct {
	data []byte
	err  error
}

func (r *errorByteReader) ReadByte() (byte, error) {
	if len(r.data) == 0 {
		return 0, r.err
	}
	b := r.data[0]
	r.data = r.data[1:]
	return b, nil
}

func TestFindAllVariants(t *testing.T) {
	matcher := CompileStrings([]string{"he", "she", "his", "hers"})
	text := []byte("ushers")
	want := matcher.FindByteSliceN(text, -1)

	if got := matcher.FindAllByteSlice(text); !reflect.DeepEqual(got, want) {
		t.Errorf("FindAllByteSlice() = %v, want %v", got, want)
	}
	if got := matcher.FindAllString(string(text)); !reflect.DeepEqual(got, want) {
		t.Errorf("FindAllString() = %v, want %v", got, want)
	}
	if got := matcher.FindAll(slices.Values(text)); !reflect.DeepEqual(got, want) {
		t.Errorf("FindAll() = %v, want %v", got, want)
	}
	if got, err := matcher.FindAllRead(bytes.NewReader(text)); err != nil {
		t.Errorf("FindAllRead() error = %v", err)
	} else if !reflect.DeepEqual(got, want) {
		t.Errorf("FindAllRead() = %v, want %v", got, want)
	}
}

func TestIncreaseSize(t *testing.T) {
	m := &Matcher{
		[]int{5, 0, 0},
		[]int{0, 0, 0},
		[]int{0, 0, 0},
		[][]int{},
		100,
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
		100,
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
		100,
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
		100,
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
		100,
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
		100,
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

func ExampleMatcher_FindByteSliceN() {
	matcher := CompileByteSlices([][]byte{
		[]byte("he"),
		[]byte("she"),
		[]byte("his"),
		[]byte("hers"),
		[]byte("she"),
	})
	fmt.Print(matcher.FindByteSliceN([]byte("ushers"), -1))

	// Output:
	// [{ "he" 2 } { "she" 1 } { "she" 1 } { "hers" 2 }]
}

func ExampleMatcher_FindStringN() {
	matcher := CompileStrings([]string{
		"he",
		"she",
		"his",
		"hers",
		"she",
	})
	fmt.Print(matcher.FindStringN("ushers", -1))

	// Output:
	// [{ "he" 2 } { "she" 1 } { "she" 1 } { "hers" 2 }]
}

func Test_trail(t *testing.T) {
	t.Run("continuous", func(t *testing.T) {
		trl := newTrail(5)
		for b := byte('a'); b < byte('e'); b++ {
			trl.next(b)
		}
		if w := string(trl.word(3)); w != "bcd" {
			t.Errorf("word '%s'", w)
		}
	})
	t.Run("wrapped around", func(t *testing.T) {
		trl := newTrail(5)
		for b := byte('a'); b < byte('h'); b++ {
			trl.next(b)
		}
		if w := string(trl.word(3)); w != "efg" {
			t.Errorf("word '%s'", w)
		}
	})
}
