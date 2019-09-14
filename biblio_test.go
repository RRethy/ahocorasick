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
			[]Match{{[]byte("potholderz"), 9}, {[]byte("MF DOOM"), 20}, {[]byte("锅"), 41}, {[]byte("持有人"), 51}},
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

func BenchmarkCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Compile([][]byte{
			[]byte("pssahe"),
			[]byte("pssashe"),
			[]byte("pssathey"),
			[]byte("pssatheir"),
			[]byte("pssawhere"),
			[]byte("pssabear"),
			[]byte("pssataratula"),
			[]byte("pssaadam"),
			[]byte("pssaregard-rethy"),
			[]byte("pssapanda"),
			[]byte("pssabear"),
			[]byte("pssagolang"),
			[]byte("pssahis"),
			[]byte("pssahers"),
			[]byte("pssaher"),
			[]byte("psahe"),
			[]byte("psashe"),
			[]byte("psathey"),
			[]byte("psatheir"),
			[]byte("psawhere"),
			[]byte("psabear"),
			[]byte("psataratula"),
			[]byte("psaadam"),
			[]byte("psaregard-rethy"),
			[]byte("psapanda"),
			[]byte("psabear"),
			[]byte("psagolang"),
			[]byte("psahis"),
			[]byte("psahers"),
			[]byte("psaher"),
			[]byte("pshe"),
			[]byte("psshe"),
			[]byte("psthey"),
			[]byte("pstheir"),
			[]byte("pswhere"),
			[]byte("psbear"),
			[]byte("pstaratula"),
			[]byte("psadam"),
			[]byte("psregard-rethy"),
			[]byte("pspanda"),
			[]byte("psbear"),
			[]byte("psgolang"),
			[]byte("pshis"),
			[]byte("pshers"),
			[]byte("psher"),
			[]byte("psahe"),
			[]byte("psashe"),
			[]byte("psathey"),
			[]byte("psatheir"),
			[]byte("psawhere"),
			[]byte("psabear"),
			[]byte("psataratula"),
			[]byte("psaadam"),
			[]byte("psaregard-rethy"),
			[]byte("psapanda"),
			[]byte("psabear"),
			[]byte("psagolang"),
			[]byte("psahis"),
			[]byte("psahers"),
			[]byte("psaher"),
			[]byte("pahe"),
			[]byte("pashe"),
			[]byte("pathey"),
			[]byte("patheir"),
			[]byte("pawhere"),
			[]byte("pabear"),
			[]byte("pataratula"),
			[]byte("paadam"),
			[]byte("paregard-rethy"),
			[]byte("papanda"),
			[]byte("pabear"),
			[]byte("pagolang"),
			[]byte("pahis"),
			[]byte("pahers"),
			[]byte("paher"),
			[]byte("phe"),
			[]byte("pshe"),
			[]byte("pthey"),
			[]byte("ptheir"),
			[]byte("pwhere"),
			[]byte("pbear"),
			[]byte("ptaratula"),
			[]byte("padam"),
			[]byte("pregard-rethy"),
			[]byte("ppanda"),
			[]byte("pbear"),
			[]byte("pgolang"),
			[]byte("phis"),
			[]byte("phers"),
			[]byte("pher"),
			[]byte("ssahe"),
			[]byte("ssashe"),
			[]byte("ssathey"),
			[]byte("ssatheir"),
			[]byte("ssawhere"),
			[]byte("ssabear"),
			[]byte("ssataratula"),
			[]byte("ssaadam"),
			[]byte("ssaregard-rethy"),
			[]byte("ssapanda"),
			[]byte("ssabear"),
			[]byte("ssagolang"),
			[]byte("ssahis"),
			[]byte("ssahers"),
			[]byte("ssaher"),
			[]byte("sahe"),
			[]byte("sashe"),
			[]byte("sathey"),
			[]byte("satheir"),
			[]byte("sawhere"),
			[]byte("sabear"),
			[]byte("sataratula"),
			[]byte("saadam"),
			[]byte("saregard-rethy"),
			[]byte("sapanda"),
			[]byte("sabear"),
			[]byte("sagolang"),
			[]byte("sahis"),
			[]byte("sahers"),
			[]byte("saher"),
			[]byte("she"),
			[]byte("sshe"),
			[]byte("sthey"),
			[]byte("stheir"),
			[]byte("swhere"),
			[]byte("sbear"),
			[]byte("staratula"),
			[]byte("sadam"),
			[]byte("sregard-rethy"),
			[]byte("spanda"),
			[]byte("sbear"),
			[]byte("sgolang"),
			[]byte("shis"),
			[]byte("shers"),
			[]byte("sher"),
			[]byte("sahe"),
			[]byte("sashe"),
			[]byte("sathey"),
			[]byte("satheir"),
			[]byte("sawhere"),
			[]byte("sabear"),
			[]byte("sataratula"),
			[]byte("saadam"),
			[]byte("saregard-rethy"),
			[]byte("sapanda"),
			[]byte("sabear"),
			[]byte("sagolang"),
			[]byte("sahis"),
			[]byte("sahers"),
			[]byte("saher"),
			[]byte("ahe"),
			[]byte("ashe"),
			[]byte("athey"),
			[]byte("atheir"),
			[]byte("awhere"),
			[]byte("abear"),
			[]byte("ataratula"),
			[]byte("aadam"),
			[]byte("aregard-rethy"),
			[]byte("apanda"),
			[]byte("abear"),
			[]byte("agolang"),
			[]byte("ahis"),
			[]byte("ahers"),
			[]byte("aher"),
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
