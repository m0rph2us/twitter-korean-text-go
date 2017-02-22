package util

import (
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

type KoreanPos int

const (
	Null        KoreanPos = iota // 0
	None        KoreanPos = iota // 1
	Noun        KoreanPos = iota // 2
	Verb        KoreanPos = iota // 3
	Adjective   KoreanPos = iota // 4
	Adverb      KoreanPos = iota // 5
	Determiner  KoreanPos = iota // 6
	Exclamation KoreanPos = iota // 7
	Josa        KoreanPos = iota // 8
	Eomi        KoreanPos = iota // 9
	PreEomi     KoreanPos = iota // 10
	Conjunction KoreanPos = iota // 11
	NounPrefix  KoreanPos = iota // 12
	VerbPrefix  KoreanPos = iota // 13
	Suffix      KoreanPos = iota // 14
	Unknown     KoreanPos = iota // 15

	// Chunk level POS
	Korean         KoreanPos = iota // 16
	Foreign        KoreanPos = iota // 17
	Number         KoreanPos = iota // 18
	KoreanParticle KoreanPos = iota // 19
	Alpha          KoreanPos = iota // 20
	Punctuation    KoreanPos = iota // 21
	Hashtag        KoreanPos = iota // 22
	ScreenName     KoreanPos = iota // 23
	Email          KoreanPos = iota // 24
	URL            KoreanPos = iota // 25
	CashTag        KoreanPos = iota // 26

	// Functional POS
	Space      KoreanPos = iota // 27
	Others     KoreanPos = iota // 28
	ProperNoun KoreanPos = iota // 29
)

var OtherPoses data_structure.Set

var shortCut map[rune]KoreanPos

type KoreanPosTrie struct {
	CurPos   KoreanPos
	NextTrie []KoreanPosTrie
	Ending   KoreanPos
}

func (t KoreanPosTrie) Compare(that KoreanPosTrie) bool {
	a := true
	b := true
	c := true

	if t.CurPos != that.CurPos {
		a = false
	}

	if t.Ending != that.Ending {
		a = false
	}

	if len(t.NextTrie) == len(that.NextTrie) {
		for i, v := range t.NextTrie {
			if v.Compare(that.NextTrie[i]) == false {
				c = false
				break
			}
		}
	}

	return a == b && b == c

}

var SelfNode KoreanPosTrie

var Predicates data_structure.Set

func init() {
	OtherPoses = data_structure.NewLinkedSet(Korean, Foreign, Number, KoreanParticle, Alpha,
		Punctuation, Hashtag, ScreenName, Email, URL, CashTag)

	// may need ordered map
	shortCut = map[rune]KoreanPos{
		'N': Noun,
		'V': Verb,
		'J': Adjective,
		'A': Adverb,
		'D': Determiner,
		'E': Exclamation,
		'C': Conjunction,

		'j': Josa,
		'e': Eomi,
		'r': PreEomi,
		'p': NounPrefix,
		'v': VerbPrefix,
		's': Suffix,

		'a': Alpha,
		'n': Number,

		'o': Others,
	}

	SelfNode = KoreanPosTrie{Null, nil, None}

	Predicates = data_structure.NewLinkedSet(Verb, Adjective)
}

func buildTrie(s string, endingPos KoreanPos) []KoreanPosTrie {
	isFinal := func(rest string) bool {
		isNextOptional := scala.FoldLeft([]rune(rest), true,
			func(folded interface{}, key interface{}, value interface{}) interface{} {
				if value.(rune) == '+' || value.(rune) == '1' {
					return false
				} else {
					return folded.(bool)
				}
			}).(interface{})

		return len(rest) == 0 || isNextOptional.(bool)
	}

	l := utf8.RuneCountInString(s)

	if l < 2 {
		return []KoreanPosTrie{}
	}

	pos := shortCut[[]rune(s)[0]]
	rule := []rune(s)[1]

	var rest string = ""
	if l > 1 {
		rest = string([]rune(s)[2:l])
	}

	var end KoreanPos = None
	if isFinal(rest) {
		end = endingPos
	}

	if rule == '+' {
		t := []KoreanPosTrie{SelfNode}
		a := buildTrie(rest, endingPos)
		t = append(t, a...)

		return []KoreanPosTrie{{pos, t, end}}
	} else if rule == '*' {
		t := []KoreanPosTrie{SelfNode}
		a := buildTrie(rest, endingPos)
		t = append(t, a...)

		x := []KoreanPosTrie{{pos, t, end}}
		y := buildTrie(rest, endingPos)
		x = append(x, y...)
		return x
	} else if rule == '1' {
		return []KoreanPosTrie{{pos, buildTrie(rest, endingPos), end}}
	} else if rule == '0' {
		x := []KoreanPosTrie{{pos, buildTrie(rest, endingPos), end}}
		y := buildTrie(rest, endingPos)
		x = append(x, y...)
		return x
	}

	panic("buildTrie panic")
}

func convertInterfaceToSliceOfKoreanPosTrie(value interface{}) []KoreanPosTrie {
	x := []KoreanPosTrie{}
	for _, v := range value.([]KoreanPosTrie) {
		x = append(x, v)
	}
	return x
}

func GetTrie(sequences data_structure.Map /*map[string]KoreanPos*/) []KoreanPosTrie {
	foldLeft := func (value interface{}, initValue interface{},
		f func(folded interface{}, key interface{}, value interface{}) interface{}) interface{} {

		v := initValue

		for d := range value.(data_structure.Map).Iter() {
			v = f(v, d.Key, d.Value)
		}

		return v
	}

	x := foldLeft(sequences, []KoreanPosTrie{}, func(folded interface{}, key interface{}, value interface{}) interface{} {
		k := append(buildTrie(key.(string), value.(KoreanPos)), folded.([]KoreanPosTrie)...)
		return k
	}).(interface{})

	return convertInterfaceToSliceOfKoreanPosTrie(x)
}
