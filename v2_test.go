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
