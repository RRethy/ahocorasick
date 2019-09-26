package biblio

import (
	"bufio"
	"bytes"
	bobu "github.com/BobuSumisu/aho-corasick"
	anknown "github.com/anknown/ahocorasick"
	_ "github.com/cloudflare/ahocorasick"
	_ "github.com/iohub/ahocorasick"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestFindAllByteSlice(t *testing.T) {
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
		matcher := CompileByteSlices(test.patterns)
		got := matcher.FindAllByteSlice(test.text)
		if !(len(got) == 0 && len(test.expected) == 0) && !reflect.DeepEqual(got, test.expected) {
			t.Errorf(`
        Text:     %s
		Expected: %v
		Got:      %v
		`, test.text, test.expected, got)
		}
	}
}

func TestIncreaseSize(t *testing.T) {
	m := &Matcher{
		[]int{5, 0, 0},
		[]int{0, 0, 0},
		[]int{0, 0, 0},
		map[int][]*[]byte{},
	}
	m.increaseSize(1)
	if !reflect.DeepEqual(m.Base, []int{5, 0, 0, -3}) {
		t.Errorf("Got: %v\n", m.Base)
	}
	if !reflect.DeepEqual(m.Check, []int{-3, 0, 0, -1}) {
		t.Errorf("Got: %v\n", m.Check)
	}

	m.increaseSize(1)
	if !reflect.DeepEqual(m.Base, []int{5, 0, 0, -4, -3}) {
		t.Errorf("Got: %v\n", m.Base)
	}
	if !reflect.DeepEqual(m.Check, []int{-3, 0, 0, -4, -1}) {
		t.Errorf("Got: %v\n", m.Check)
	}

	m.increaseSize(1)
	if !reflect.DeepEqual(m.Base, []int{5, 0, 0, -5, -3, -4}) {
		t.Errorf("Got: %v\n", m.Base)
	}
	if !reflect.DeepEqual(m.Check, []int{-3, 0, 0, -4, -5, -1}) {
		t.Errorf("Got: %v\n", m.Check)
	}

	m = &Matcher{
		[]int{5, 0, 0},
		[]int{0, 0, 0},
		[]int{0, 0, 0},
		map[int][]*[]byte{},
	}
	m.increaseSize(3)
	if !reflect.DeepEqual(m.Base, []int{5, 0, 0, -5, -3, -4}) {
		t.Errorf("Got: %v\n", m.Base)
	}
	if !reflect.DeepEqual(m.Check, []int{-3, 0, 0, -4, -5, -1}) {
		t.Errorf("Got: %v\n", m.Check)
	}

	m.increaseSize(3)
	if !reflect.DeepEqual(m.Base, []int{5, 0, 0, -8, -3, -4, -5, -6, -7}) {
		t.Errorf("Got: %v\n", m.Base)
	}
	if !reflect.DeepEqual(m.Check, []int{-3, 0, 0, -4, -5, -6, -7, -8, -1}) {
		t.Errorf("Got: %v\n", m.Check)
	}

	m = &Matcher{
		[]int{0},
		[]int{0},
		[]int{0},
		map[int][]*[]byte{},
	}
	m.increaseSize(5)
	if !reflect.DeepEqual(m.Base, []int{0, -5, -1, -2, -3, -4}) {
		t.Errorf("Got: %v\n", m.Base)
	}
	if !reflect.DeepEqual(m.Check, []int{-1, -2, -3, -4, -5, -1}) {
		t.Errorf("Got: %v\n", m.Check)
	}

	m = &Matcher{
		[]int{-103, -1867},
		[]int{0, 0},
		[]int{},
		map[int][]*[]byte{},
	}
	m.increaseSize(5)
	if !reflect.DeepEqual(m.Base, []int{-103, -1867, -6, -2, -3, -4, -5}) {
		t.Errorf("Got: %v\n", m.Base)
	}
	if !reflect.DeepEqual(m.Check, []int{-2, 0, -3, -4, -5, -6, -1}) {
		t.Errorf("Got: %v\n", m.Check)
	}
}

func TestNextFreeState(t *testing.T) {
	m := &Matcher{
		[]int{5, 0, 0, -3},
		[]int{-3, 0, 0, -1},
		[]int{},
		map[int][]*[]byte{},
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
		map[int][]*[]byte{},
	}
	m.increaseSize(5)
	m.occupyState(3, 1)
	m.occupyState(4, 1)
	m.occupyState(8, 1)
	m.occupyState(6, 1)
	m.occupyState(5, 1)
	m.occupyState(7, 1)
	if !reflect.DeepEqual(m.Base, []int{5, 0, 0, -1867, -1867, -1867, -1867, -1867, -1867}) {
		t.Errorf("Got: %v\n", m.Base)
	}
	if !reflect.DeepEqual(m.Check, []int{0, 0, 0, 1, 1, 1, 1, 1, 1}) {
		t.Errorf("Got: %v\n", m.Check)
	}
}

func readLines(fname string, every int) ([][]byte, error) {
	file, err := os.Open(fname)
	defer file.Close()
	var pattens [][]byte
	if err != nil {
		return pattens, err
	}

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		if i%every == 0 {
			pattens = append(pattens, scanner.Bytes())
		}
		i++
	}
	return pattens, nil
}

func readBytes(fname string, every int) ([]byte, error) {
	file, err := os.Open(fname)
	defer file.Close()
	var bytes []byte
	if err != nil {
		return bytes, err
	}

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		if i%every == 0 {
			bytes = append(bytes, scanner.Bytes()...)
		}
		i++
	}
	return bytes, nil
}

func readRunes(filename string, every int) ([][]rune, error) {
	dict := [][]rune{}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(f)
	i := 0
	for {
		l, err := r.ReadBytes('\n')
		if err != nil || err == io.EOF {
			break
		}
		l = bytes.TrimSpace(l)

		if i%every == 0 {
			dict = append(dict, bytes.Runes(l))
		}
		i++
	}

	return dict, nil
}

func benchBiblioCompileByteSlices(b *testing.B, every int) {
	patterns, err := readLines("./words.txt", every)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		CompileByteSlices(patterns)
	}
}

func benchBiblioFindAllByteSlice(b *testing.B, every int) {
	patterns, err := readLines("./words.txt", 1000)
	if err != nil {
		b.Error(err)
		return
	}

	text, err := readBytes("./war-and-peace.txt", every)
	if err != nil {
		b.Error(err)
		return
	}

	m := CompileByteSlices(patterns)
	for i := 0; i < b.N; i++ {
		m.FindAllByteSlice(text)
	}
}

func BenchmarkBiblio1ModCompileByteSlices(b *testing.B) {
	benchBiblioCompileByteSlices(b, 1)
}

func BenchmarkBiblio10ModCompileByteSlices(b *testing.B) {
	benchBiblioCompileByteSlices(b, 10)
}

func BenchmarkBiblio100ModCompileByteSlices(b *testing.B) {
	benchBiblioCompileByteSlices(b, 100)
}

func BenchmarkBiblio1000ModCompileByteSlices(b *testing.B) {
	benchBiblioCompileByteSlices(b, 1000)
}

func BenchmarkBiblio10000ModCompileByteSlices(b *testing.B) {
	benchBiblioCompileByteSlices(b, 10000)
}

func BenchmarkBiblio1ModFindAllByteSlice(b *testing.B) {
	benchBiblioFindAllByteSlice(b, 1)
}

func BenchmarkBiblio10ModFindAllByteSlice(b *testing.B) {
	benchBiblioFindAllByteSlice(b, 10)
}

func BenchmarkBiblio100ModFindAllByteSlice(b *testing.B) {
	benchBiblioFindAllByteSlice(b, 100)
}

func BenchmarkBiblio1000ModFindAllByteSlice(b *testing.B) {
	benchBiblioFindAllByteSlice(b, 1000)
}

// BOBU
func benchBobuCompile(b *testing.B, every int) {
	patterns, err := readLines("./words.txt", every)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		bobu.NewTrieBuilder().AddPatterns(patterns).Build()
	}
}

func benchBobuFind(b *testing.B, every int) {
	patterns, err := readLines("./words.txt", 1000)
	if err != nil {
		b.Error(err)
		return
	}

	text, err := readBytes("./war-and-peace.txt", every)
	if err != nil {
		b.Error(err)
		return
	}

	t := bobu.NewTrieBuilder().AddPatterns(patterns).Build()
	for i := 0; i < b.N; i++ {
		t.Match(text)
	}
}

func BenchmarkBobu1ModCompile(b *testing.B) {
	benchBobuCompile(b, 1)
}

func BenchmarkBobu10ModCompile(b *testing.B) {
	benchBobuCompile(b, 10)
}

func BenchmarkBobu100ModCompile(b *testing.B) {
	benchBobuCompile(b, 100)
}

func BenchmarkBobu1000ModCompile(b *testing.B) {
	benchBobuCompile(b, 1000)
}

func BenchmarkBobu10000ModCompile(b *testing.B) {
	benchBobuCompile(b, 10000)
}

func BenchmarkBobu1ModFind(b *testing.B) {
	benchBobuFind(b, 1)
}

func BenchmarkBobu10ModFind(b *testing.B) {
	benchBobuFind(b, 10)
}

func BenchmarkBobu100ModFind(b *testing.B) {
	benchBobuFind(b, 100)
}

func BenchmarkBobu1000ModFind(b *testing.B) {
	benchBobuFind(b, 1000)
}

// ANKNOWN
func benchAnknownCompile(b *testing.B, every int) {
	patterns, err := readRunes("./words.txt", every)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		m := new(anknown.Machine)
		m.Build(patterns)
	}
}

func benchAnknownFind(b *testing.B, every int) {
	patterns, err := readRunes("./words.txt", 1000)
	if err != nil {
		b.Error(err)
		return
	}
	m := new(anknown.Machine)
	m.Build(patterns)

	text, err := readBytes("./war-and-peace.txt", every)
	if err != nil {
		b.Error(err)
		return
	}
	textRunes := bytes.Runes(text)

	for i := 0; i < b.N; i++ {
		m.MultiPatternSearch(textRunes, false)
	}
}

func BenchmarkAnknown1ModCompile(b *testing.B) {
	benchAnknownCompile(b, 1)
}

func BenchmarkAnknown10ModCompile(b *testing.B) {
	benchAnknownCompile(b, 10)
}

func BenchmarkAnknown100ModCompile(b *testing.B) {
	benchAnknownCompile(b, 100)
}

func BenchmarkAnknown1000ModCompile(b *testing.B) {
	benchAnknownCompile(b, 1000)
}

func BenchmarkAnknown10000ModCompile(b *testing.B) {
	benchAnknownCompile(b, 10000)
}

func BenchmarkAnknown1ModFind(b *testing.B) {
	benchAnknownFind(b, 1)
}

func BenchmarkAnknown10ModFind(b *testing.B) {
	benchAnknownFind(b, 10)
}

func BenchmarkAnknown100ModFind(b *testing.B) {
	benchAnknownFind(b, 100)
}

func BenchmarkAnknown1000ModFind(b *testing.B) {
	benchAnknownFind(b, 1000)
}
