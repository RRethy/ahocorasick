package biblio

import (
	"bufio"
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
		map[int][][]byte{},
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
		map[int][][]byte{},
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
		map[int][][]byte{},
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
		map[int][][]byte{},
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
		map[int][][]byte{},
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
		map[int][][]byte{},
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

func benchCompileByteSlices(b *testing.B, every int) {
	patterns, err := readLines("./words.txt", every)
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		CompileByteSlices(patterns)
	}
}

func benchFindAllByteSlice(b *testing.B, every int) {
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

	for i := 0; i < b.N; i++ {
		m := CompileByteSlices(patterns)
		m.FindAllByteSlice(text)
	}
}

func BenchmarkCompileByteSlicesMassive(b *testing.B) {
	benchCompileByteSlices(b, 10)
}

func BenchmarkFindAllByteSliceMassive(b *testing.B) {
	benchFindAllByteSlice(b, 1)
}

func BenchmarkCompileByteSlicesLarge(b *testing.B) {
	benchCompileByteSlices(b, 100)
}

func BenchmarkFindAllByteSliceLarge(b *testing.B) {
	benchFindAllByteSlice(b, 10)
}

func BenchmarkCompileByteSlicesMedium(b *testing.B) {
	benchCompileByteSlices(b, 1000)
}

func BenchmarkFindAllByteSliceMedium(b *testing.B) {
	benchFindAllByteSlice(b, 100)
}

func BenchmarkCompileByteSlicesSmall(b *testing.B) {
	benchCompileByteSlices(b, 10000)
}

func BenchmarkFindAllByteSliceSmall(b *testing.B) {
	benchFindAllByteSlice(b, 1000)
}
