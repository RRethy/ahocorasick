package biblio

import (
	"reflect"
	"testing"
)

func TestFindAll(t *testing.T) {
	tests := []struct {
		patterns [][]byte
		expected []Match
		text     []byte
	}{
		{
			[][]byte{[]byte("he"), []byte("she"), []byte("his"), []byte("hers")},
			[]Match{{[]byte("he"), 3}, {[]byte("she"), 3}, {[]byte("hers"), 5}},
			[]byte("ushers"),
		},
		{
			[][]byte{[]byte("they"), []byte("their"), []byte("theyre"), []byte("the"), []byte("tea"), []byte("te"), []byte("team"), []byte("go"), []byte("goo"), []byte("good"), []byte("oode")},
			[]Match{{[]byte("the"), 2}, {[]byte("they"), 3}, {[]byte("theyre"), 5}, {[]byte("go"), 14}, {[]byte("goo"), 15}, {[]byte("good"), 16}, {[]byte("oode"), 17}, {[]byte("te"), 20}, {[]byte("tea"), 21}, {[]byte("team"), 22}},
			[]byte("theyre not a goode team"),
		},
		// {
		// 	[][]byte{[]byte("a")},
		// 	[]Match{{[]byte("a"), 0}, {[]byte("a"), 1}, {[]byte("a"), 2}, {[]byte("a"), 5}, {[]byte("a"), 7}, {[]byte("a"), 9}, {[]byte("a"), 11}},
		// 	[]byte("aaabbabababa"),
		// },
		// {
		// 	[][]byte{},
		// 	[]Match{},
		// 	[]byte("there is no patterns"),
		// },
		// {
		// 	[][]byte{[]byte("锅"), []byte("持有人"), []byte("potholderz"), []byte("MF DOOM")},
		// 	[]Match{{[]byte("potholderz"), 0}, {[]byte("MF DOOM"), 14}, {[]byte("锅"), 39}, {[]byte("持有人"), 43}},
		// 	[]byte("potholderz by MF DOOM hot shit aw shit 锅 持有人"),
		// },
	}
	for _, test := range tests {
		matcher := Compile(test.patterns)
		got := matcher.FindAll(test.text)
		if !(len(got) == 0 && len(test.expected) == 0) && !reflect.DeepEqual(got, test.expected) {
			t.Errorf(`
Expected: %v
Got:      %v
`, test.expected, got)
		}
	}
}

func BenchmarkFindAll(b *testing.B) {
	matcher := Compile([][]byte{
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
		matcher.FindAll([]byte(`
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
