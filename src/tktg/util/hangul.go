package util

import (
	"github.com/acidd15/go-scala-util/src/scala"
)

var hangulBase rune = 0xAC00

var onsetBase rune = 21 * 28
var vowelBase rune = 28

var onsetList []rune = []rune{
	'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ',
	'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

var vowelList []rune = []rune{
	'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ', 'ㅔ',
	'ㅕ', 'ㅖ', 'ㅗ', 'ㅘ', 'ㅙ', 'ㅚ',
	'ㅛ', 'ㅜ', 'ㅝ', 'ㅞ', 'ㅟ', 'ㅠ',
	'ㅡ', 'ㅢ', 'ㅣ',
}

var codaList []rune = []rune{
	' ', 'ㄱ', 'ㄲ', 'ㄳ', 'ㄴ', 'ㄵ', 'ㄶ', 'ㄷ',
	'ㄹ', 'ㄺ', 'ㄻ', 'ㄼ', 'ㄽ', 'ㄾ', 'ㄿ', 'ㅀ',
	'ㅁ', 'ㅂ', 'ㅄ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅊ',
	'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

var doubleCodas map[rune]DoubleCoda = map[rune]DoubleCoda{
	'ㄳ': {'ㄱ', 'ㅅ'},
	'ㄵ': {'ㄴ', 'ㅈ'},
	'ㄶ': {'ㄴ', 'ㅎ'},
	'ㄺ': {'ㄹ', 'ㄱ'},
	'ㄻ': {'ㄹ', 'ㅁ'},
	'ㄼ': {'ㄹ', 'ㅂ'},
	'ㄽ': {'ㄹ', 'ㅅ'},
	'ㄾ': {'ㄹ', 'ㅌ'},
	'ㄿ': {'ㄹ', 'ㅍ'},
	'ㅀ': {'ㄹ', 'ㅎ'},
	'ㅄ': {'ㅂ', 'ㅅ'},
}

var onsetMap map[rune]int
var vowelMap map[rune]int
var CodaMap map[rune]int

type HangulRune struct {
	Onset rune
	Vowel rune
	Coda  rune
}

type DoubleCoda struct {
	first  rune
	second rune
}

func init() {
	onsetMap = convertMap(scala.ToMap(scala.ZipWithIndex(onsetList)).(map[interface{}]interface{}))
	vowelMap = convertMap(scala.ToMap(scala.ZipWithIndex(vowelList)).(map[interface{}]interface{}))
	CodaMap = convertMap(scala.ToMap(scala.ZipWithIndex(codaList)).(map[interface{}]interface{}))
}

func convertMap(m map[interface{}]interface{}) map[rune]int {
	v := map[rune]int{}
	for key, value := range m {
		v[key.(rune)] = value.(int)
	}
	return v
}

func DecomposeHangul(c rune) HangulRune {
	_, e1 := onsetMap[c]
	_, e2 := vowelMap[c]
	_, e3 := CodaMap[c]

	if !(!e1 && !e2 && !e3) {
		panic("Input character is not a valid Korean character")
	}

	u := c - hangulBase

	return HangulRune{
		onsetList[u/onsetBase],
		vowelList[(u%onsetBase)/vowelBase],
		codaList[u%vowelBase],
	}
}

func isHangulMatch(s HangulRune, d HangulRune) bool {
	return (s.Onset == d.Onset && s.Vowel == d.Vowel && s.Coda == d.Coda)
}

func hasCoda(c rune) bool {
	return (c-hangulBase)%vowelBase > 0
}

func ComposeHangulFull(onset rune, vowel rune, coda rune) rune {
	if !(onset != ' ' && vowel != ' ') {
		panic("Input characters are not valid")
	}

	return hangulBase + (rune(onsetMap[onset]) * onsetBase) + (rune(vowelMap[vowel]) * vowelBase) + rune(CodaMap[coda])
}

func composeHangul(hc HangulRune) rune {
	return ComposeHangulFull(hc.Onset, hc.Vowel, hc.Coda)
}
