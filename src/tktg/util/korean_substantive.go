package util

import (
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var josaHeadForCoda data_structure.Set

var josaHeadForNoCoda data_structure.Set

var numberRunes data_structure.Set

var numberLastRunes data_structure.Set

func init() {
	josaHeadForCoda = data_structure.NewLinkedSet('은', '이', '을', '과', '아')
	josaHeadForNoCoda = data_structure.NewLinkedSet('는', '가', '를', '와', '야', '여', '라')
	numberRunes = data_structure.NewLinkedSet(
		'일', '이', '삼', '사', '오', '육', '칠', '팔',
		'구', '천', '백', '십', '해', '경', '조', '억', '만',
	)
	numberLastRunes = data_structure.NewLinkedSet(
		'일', '이', '삼', '사', '오', '육', '칠', '팔',
		'구', '천', '백', '십', '해', '경', '조', '억',
		'만', '원', '배', '분', '초',
	)
}

func isJosaAttachable(prevRune rune, headRune rune) bool {
	return hasCoda(prevRune) && !josaHeadForNoCoda.Contains(headRune) ||
		(!hasCoda(prevRune) && !josaHeadForCoda.Contains(headRune))
}

func IsName(chunk string) bool {
	if nameDictionary["full_name"].Contains(chunk) ||
		nameDictionary["given_name"].Contains(chunk) {
		return true
	}

	l := utf8.RuneCountInString(chunk)
	if l == 3 {
		return (nameDictionary["family_name"].Contains(string([]rune(chunk)[0])) &&
			nameDictionary["given_name"].Contains(string([]rune(chunk)[1:3])))
	} else if l == 4 {
		return (nameDictionary["family_name"].Contains(string([]rune(chunk)[0:2])) &&
			nameDictionary["given_name"].Contains(string([]rune(chunk)[2:4])))
	}

	return false
}

func IsKoreanNumber(chunk string) bool {
	d := []int{}
	l := utf8.RuneCountInString(chunk)
	for i := 0; l >= i; i++ {
		d = append(d, i)
	}

	return scala.FoldLeft(d, true, func(folded interface{}, key interface{}, value interface{}) interface{} {
		if value.(int) < l-1 {
			return folded.(bool) && numberRunes.Contains(chunk[value.(int)])
		} else {
			return folded.(bool) && numberLastRunes.Contains(chunk[value.(int)])
		}
	}).(bool)
}

func IsKoreanNameVariation(chunk string) bool {
	//nounDict := koreanDictionary[Noun]

	if IsName(chunk) {
		return true
	}

	l := utf8.RuneCountInString(chunk)

	if l < 3 || l > 5 {
		return false
	}

	decomposed := ConvertSliceOfInterfaceToSliceOfHangulRune(
		scala.Map([]rune(chunk), func(key interface{}, value interface{}) interface{} {
			return DecomposeHangul(value.(rune))
		}).([]interface{}))

	dl := len(decomposed)
	lastRune := HangulRune(decomposed[dl-1 : dl][0])

	if _, ok := CodaMap[lastRune.Onset]; ok == false {
		return false
	}

	if lastRune.Onset == 'ㅇ' || lastRune.Vowel != 'ㅣ' || lastRune.Coda != ' ' {
		return false
	}

	if HangulRune(decomposed[dl-2 : dl-1][0]).Coda != ' ' {
		return false
	}

	recovered := string(ConvertSliceOfInterfaceToSliceOfRune(
		scala.Map(scala.ZipWithIndex(decomposed), func(key interface{}, value interface{}) interface{} {
			i := key.(int)
			v := value.([]interface{})
			if i == l-1 {
				return '이'
			} else if i == l-2 {
				return ComposeHangulFull(v[0].(HangulRune).Onset, v[0].(HangulRune).Vowel, lastRune.Onset)
			} else {
				return composeHangul(v[0].(HangulRune))
			}
		}).([]interface{})))

	rl := utf8.RuneCountInString(recovered)
	if IsName(recovered) {
		return true
	} else if IsName(string([]rune(recovered)[0 : rl-1])) {
		return true
	}

	return false
}
