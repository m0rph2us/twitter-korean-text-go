package normalizer

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"strings"
	"tktg/shared"
	"tktg/util"
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var extendedKoreanRegex = pcre.MustCompile(`([ㄱ-ㅣ가-힣]+)`, pcre.MULTILINE|pcre.UTF8)
var koreanToNomalizeRegex = pcre.MustCompile(`([가-힣]+)(ㅋ+|ㅎ+|[ㅠㅜ]+)`, pcre.MULTILINE|pcre.UTF8)
var repeatingCharRegex = pcre.MustCompile(`(.)\1{2,}|[ㅠㅜ]{2,}`, pcre.MULTILINE|pcre.UTF8)
var repeating2CharRegex = pcre.MustCompile(`(..)\1{2,}`, pcre.MULTILINE|pcre.UTF8)
var whitespaceRegex = pcre.MustCompile(`\s+`, pcre.MULTILINE|pcre.UTF8)

var codaNException = data_structure.NewLinkedSet('근', '던', '닌', '는', '텐', '인', '운', '른', '든', '픈', '은')

func Normalize(input string) string {
	return shared.ReplaceAll(extendedKoreanRegex, input, func(m *pcre.Matcher) string {
		return normalizeKoreanChunk(string(m.Group(0)))
	})
}

func normalizeKoreanChunk(input string) string {
	endingNormalized := shared.ReplaceAll(koreanToNomalizeRegex, input, func(m *pcre.Matcher) string {
		return processNormalizationCandidate(m)
	})

	exclamationNormalized := shared.ReplaceAll(repeatingCharRegex, endingNormalized, func(m *pcre.Matcher) string {
		return string([]rune(string(m.Group(0)))[0:2])
	})

	repeatingNormalized := shared.ReplaceAll(repeating2CharRegex, exclamationNormalized, func(m *pcre.Matcher) string {
		return string([]rune(string(m.Group(0)))[0:4])
	})

	codaNNormalized := normalizeCodaN(repeatingNormalized)

	typoCorrected := correctTypo(codaNNormalized)

	return shared.ReplaceAll(whitespaceRegex, typoCorrected, func(m *pcre.Matcher) string {
		return " "
	})
}

func correctTypo(chunk string) string {
	return scala.FoldLeft(util.TypoDictionaryByLength, chunk, func(folded interface{}, key interface{}, value interface{}) interface{} {
		output := folded.(string)
		wordLen := key.(int)
		typoMap := value.(map[string]string)

		slidings := []string{}
		l := utf8.RuneCountInString(output)
		for i := 0; l >= i+wordLen; i++ {
			slidings = append(slidings, string([]rune(output)[i:i+wordLen]))
		}

		return scala.FoldLeft(slidings, output, func(folded interface{}, key interface{}, value interface{}) interface{} {
			sliceOutput := folded.(string)
			slice := value.(string)
			if _, e := typoMap[slice]; e == true {
				return strings.Replace(sliceOutput, slice, typoMap[slice], -1)
			} else {
				return sliceOutput
			}
		})
	}).(string)
}

func normalizeCodaN(chunk string) string {
	cr := []rune(chunk)
	l := len(cr)
	if l < 2 {
		return chunk
	}

	lastTwo := string(cr[l-2:])
	last := cr[l-1]
	lastTwoHead := cr[l-2]

	if util.KoreanDictionary.Get(util.Noun).(data_structure.Set).Contains(chunk) ||
		util.KoreanDictionary.Get(util.Conjunction).(data_structure.Set).Contains(chunk) ||
		util.KoreanDictionary.Get(util.Adverb).(data_structure.Set).Contains(chunk) ||
		util.KoreanDictionary.Get(util.Noun).(data_structure.Set).Contains(lastTwo) ||
		lastTwoHead < '가' || lastTwoHead > '힣' || codaNException.Contains(lastTwoHead) {
		return chunk
	}

	hc := util.DecomposeHangul(lastTwoHead)

	newHead := ""
	newHead += string(cr[0 : l-2])
	newHead += string(util.ComposeHangulFull(hc.Onset, hc.Vowel, ' '))

	if hc.Coda == 'ㄴ' &&
		(last == '데' || last == '가' || last == '지') &&
		util.KoreanDictionary.Get(util.Noun).(data_structure.Set).Contains(newHead) {
		mid := "인"
		if hc.Vowel == 'ㅡ' {
			mid = "은"
		}
		return newHead + mid + string(last)
	} else {
		return chunk
	}
}

func processNormalizationCandidate(m *pcre.Matcher) string {
	chunk := string(m.Group(1))
	toNormalize := string(m.Group(2))

	normalizedChunk := ""

	cr := []rune(chunk)
	l := len(cr)

	if util.KoreanDictionary.Get(util.Noun).(data_structure.Set).Contains(chunk) ||
		util.KoreanDictionary.Get(util.Eomi).(data_structure.Set).Contains(string(cr[l-1])) ||
		util.KoreanDictionary.Get(util.Eomi).(data_structure.Set).Contains(string(cr[l-2:])) {
		normalizedChunk = chunk
	} else {
		normalizedChunk = normalizeEmotionAttachedChunk(chunk, toNormalize)
	}

	return normalizedChunk + toNormalize
}

func normalizeEmotionAttachedChunk(s string, toNormalize string) string {
	x := []rune(s)
	l := len(x)
	init := x[0 : l-1]
	last := x[l-1]

	var secondToLastDecomposed util.HangulRune
	isDefined := false

	il := len(init)
	if il > 0 {
		hc := util.DecomposeHangul(init[il-1])

		if hc.Coda == ' ' {
			isDefined = true
			secondToLastDecomposed = hc
		}
	}

	hc := util.DecomposeHangul(last)
	_, e := util.CodaMap[hc.Onset]

	if hc.Coda == 'ㅋ' || hc.Coda == 'ㅎ' {
		return string(init) + string(util.ComposeHangulFull(hc.Onset, hc.Vowel, ' '))
	} else if isDefined && hc.Vowel == []rune(toNormalize)[0] && e == true {
		return string(init[0:il-1]) + string(util.ComposeHangulFull(
			secondToLastDecomposed.Onset, secondToLastDecomposed.Vowel, hc.Onset))
	} else {
		return s
	}
}
