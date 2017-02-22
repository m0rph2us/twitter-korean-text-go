package util

import (
	"reflect"
	"tktg/shared"
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var codasCommon []rune = []rune{
	'ㅂ', 'ㅆ', 'ㄹ', 'ㄴ', 'ㅁ',
}

var codasForContraction []rune = []rune{
	'ㅆ', 'ㄹ', 'ㅁ',
}

var codasNoPast []rune = []rune{
	'ㅂ', 'ㄹ', 'ㄴ', 'ㅁ',
}

var codasSlangConsonant []rune = []rune{
	'ㅋ', 'ㅎ',
}

var codasSlangVowel []rune = []rune{
	'ㅜ', 'ㅠ',
}

var preEomiCommon []rune = []rune{
	'거', '게', '겠', '고', '구', '기', '긴', '길',
	'네', '다', '더', '던', '도', '든', '면', '자',
	'잖', '재', '져', '죠', '지', '진', '질',
}

var preEomi_1_1 []rune = []rune{
	'야', '서', '써', '도', '준',
}

var preEomi_1_2 []rune = []rune{
	'어', '었',
}

var preEomi_1_3 []rune = []rune{
	'아', '앗',
}

var preEomi_1_4 []rune = []rune{
	'워', '웠',
}

var preEomi_1_5 []rune = []rune{
	'여', '엿',
}

var preEomi_2 []rune = []rune{
	'노', '느', '니', '냐',
}

var preEomi_3 []rune = []rune{
	'러', '려', '며',
}

var preEomi_4 []rune = []rune{
	'으',
}

var preEomi_5 []rune = []rune{
	'은',
}

var preEomi_6 []rune = []rune{
	'는',
}

var preEomi_7 []rune = []rune{
	'운',
}

var preEomiRespect []rune = []rune{
	'세', '시', '실', '신', '셔', '습', '셨', '십',
}

var preEomiVowel []rune = []rune{}

func init() {
	preEomiVowel = append(preEomiVowel, preEomiCommon...)
	preEomiVowel = append(preEomiVowel, preEomi_2...)
	preEomiVowel = append(preEomiVowel, preEomi_3...)
	preEomiVowel = append(preEomiVowel, preEomiRespect...)
}

func addPreEomi(lastRune rune, runesToAdd []rune) []string {
	x := []string{}

	for _, v := range runesToAdd {
		x = append(x, string(lastRune)+string(v))
	}

	return x
}

func conjugatePredicatesToSet(words data_structure.Set, isAdjective bool) data_structure.Set {
	expanded := conjugatePredicated(words, isAdjective)
	return expanded
}

func getInit(s string) string {
	l := utf8.RuneCountInString(s)

	if l == 0 {
		panic("argument must have some string.")
	}

	return string([]rune(s)[0 : l-1])
}

func getLast(s string) rune {
	l := utf8.RuneCountInString(s)

	if l == 0 {
		panic("argument must have some string.")
	}

	return ([]rune(s))[l-1 : l][0]
}

func getPreEomi_2_6(c rune) []string {
	apeArg1 := []rune{}
	apeArg1 = append(apeArg1, preEomi_2...)
	apeArg1 = append(apeArg1, preEomi_6...)

	return addPreEomi(c, apeArg1)
}

func getPreEomi_1_4_7(c rune) []string {
	apeArg1 := []rune{}
	apeArg1 = append(apeArg1, preEomi_1_4...)
	apeArg1 = append(apeArg1, preEomi_7...)

	return addPreEomi(c, apeArg1)
}

func conjugatePredicated(words data_structure.Set, isAdjective bool) data_structure.Set {
	flatMap := func(set data_structure.Set, f func(value interface{}) interface{}) data_structure.Set {
		strSet := data_structure.NewLinkedSet()
		for v := range set.Iter() {
			x := f(v)

			to := reflect.TypeOf(x)

			if to.Kind() == reflect.Array || to.Kind() == reflect.Slice {
				for _, d := range x.([]string) {
					strSet.Add(d)
				}
			} else {
				strSet.Add(x.(string))
			}
		}
		return strSet
	}

	expanded := flatMap(words, func(word interface{}) interface{} {
		w := word.(string)
		init := getInit(w)
		lastRune := getLast(w)
		lastRuneString := string(lastRune)
		lastRuneDecomposed := DecomposeHangul(lastRune)
		expandedLast := []string{}

		o := lastRuneDecomposed.Onset
		v := lastRuneDecomposed.Vowel
		c := lastRuneDecomposed.Coda

		// Cases without codas
		// 하다, special case
		if isHangulMatch(lastRuneDecomposed, HangulRune{'ㅎ', 'ㅏ', ' '}) {
			var endings []string

			if isAdjective {
				endings = []string{"합", "해", "히", "하"}
			} else {
				endings = []string{"합", "해"}
			}

			apeArg1 := []rune{}
			apeArg1 = append(apeArg1, preEomiCommon...)
			apeArg1 = append(apeArg1, preEomi_2...)
			apeArg1 = append(apeArg1, preEomi_6...)
			apeArg1 = append(apeArg1, preEomiRespect...)

			ape1 := addPreEomi(lastRune, apeArg1)

			ccm := scala.Map(codasCommon, func(key interface{}, value interface{}) interface{} {
				c := value.(rune)
				if c == 'ㅆ' {
					return string(ComposeHangulFull('ㅎ', 'ㅐ', c))
				} else {
					return string(ComposeHangulFull('ㅎ', 'ㅏ', c))
				}
			}).([]interface{})

			apeArg2 := []rune{}
			apeArg2 = append(apeArg2, preEomiVowel...)
			apeArg2 = append(apeArg2, preEomi_1_5...)
			apeArg2 = append(apeArg2, preEomi_6...)

			ape2 := addPreEomi('하', apeArg2)

			ape3 := addPreEomi('해', preEomi_1_1)

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(ccm)...)
			expandedLast = append(expandedLast, ape2...)
			expandedLast = append(expandedLast, ape3...)
			expandedLast = append(expandedLast, endings...)
		} else
		// 쏘다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅗ', ' '}) {
			apeArg1 := []rune{}
			apeArg1 = append(apeArg1, preEomiVowel...)
			apeArg1 = append(apeArg1, preEomi_2...)
			apeArg1 = append(apeArg1, preEomi_1_3...)
			apeArg1 = append(apeArg1, preEomi_6...)

			ape1 := addPreEomi(lastRune, apeArg1)

			cnpm := scala.Map(codasNoPast, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, 'ㅗ', value.(rune)))
			}).([]interface{})

			s := []string{
				string(ComposeHangulFull(o, 'ㅘ', ' ')),
				string(ComposeHangulFull(o, 'ㅘ', 'ㅆ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(cnpm)...)
			expandedLast = append(expandedLast, s...)
		} else
		// 맞추다 겨누다 재우다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅜ', ' '}) {
			apeArg1 := []rune{}
			apeArg1 = append(apeArg1, preEomiVowel...)
			apeArg1 = append(apeArg1, preEomi_1_2...)
			apeArg1 = append(apeArg1, preEomi_2...)
			apeArg1 = append(apeArg1, preEomi_6...)

			ape1 := addPreEomi(lastRune, apeArg1)

			cnpm := scala.Map(codasNoPast, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, 'ㅜ', value.(rune)))
			}).([]interface{})

			s := []string{
				string(ComposeHangulFull(o, 'ㅝ', ' ')),
				string(ComposeHangulFull(o, 'ㅝ', 'ㅆ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(cnpm)...)
			expandedLast = append(expandedLast, s...)
		} else
		// 치르다, 구르다, 굴르다, 뜨다, 모으다, 고르다, 골르다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅡ', ' '}) {
			ape1 := getPreEomi_2_6(lastRune)

			cnpm := scala.Map(codasNoPast, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, 'ㅡ', value.(rune)))
			}).([]interface{})

			s := []string{
				string(ComposeHangulFull(lastRuneDecomposed.Onset, 'ㅝ', ' ')),
				string(ComposeHangulFull(lastRuneDecomposed.Onset, 'ㅓ', ' ')),
				string(ComposeHangulFull(lastRuneDecomposed.Onset, 'ㅏ', ' ')),
				string(ComposeHangulFull(lastRuneDecomposed.Onset, 'ㅝ', 'ㅆ')),
				string(ComposeHangulFull(lastRuneDecomposed.Onset, 'ㅓ', 'ㅆ')),
				string(ComposeHangulFull(lastRuneDecomposed.Onset, 'ㅏ', 'ㅆ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(cnpm)...)
			expandedLast = append(expandedLast, s...)
		} else
		// 사귀다
		if isHangulMatch(lastRuneDecomposed, HangulRune{'ㄱ', 'ㅜ', ' '}) {
			ape1 := getPreEomi_2_6(lastRune)

			cnpm := scala.Map(codasNoPast, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull('ㄱ', 'ㅜ', value.(rune)))
			}).([]interface{})

			s := []string{
				string(ComposeHangulFull('ㄱ', 'ㅕ', ' ')),
				string(ComposeHangulFull('ㄱ', 'ㅕ', 'ㅆ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(cnpm)...)
			expandedLast = append(expandedLast, s...)
		} else
		// 쥐다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅜ', ' '}) {
			cnpm := scala.Map(codasNoPast, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, 'ㅜ', value.(rune)))
			}).([]interface{})

			ape1 := getPreEomi_2_6(lastRune)

			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(cnpm)...)
			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, []string{lastRuneString}...)
		} else
		// 마시다, 엎드리다, 치다, 이다, 아니다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅣ', ' '}) {
			apeArg1 := []rune{}
			apeArg1 = append(apeArg1, preEomi_1_2...)
			apeArg1 = append(apeArg1, preEomi_2...)
			apeArg1 = append(apeArg1, preEomi_6...)

			ape1 := addPreEomi(lastRune, apeArg1)

			s := []string{
				string(ComposeHangulFull(o, 'ㅣ', 'ㅂ')) + "니",
				string(ComposeHangulFull(o, 'ㅕ', ' ')),
				string(ComposeHangulFull(o, 'ㅕ', 'ㅆ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, s...)
		} else
		// 꿰다, 꾀다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, v, ' '}) && (v == 'ㅞ' || v == 'ㅚ' || v == 'ㅙ') {
			ape1 := getPreEomi_2_6(lastRune)

			ccm := scala.Map(codasCommon, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, v, value.(rune)))
			}).([]interface{})

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(ccm)...)
			expandedLast = append(expandedLast, []string{lastRuneString}...)
		} else
		// All other vowel endings: 둘러서다, 켜다, 세다, 캐다, 차다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, v, ' '}) {
			ccm := scala.Map(codasCommon, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, v, value.(rune)))
			}).([]interface{})

			apeArg1 := []rune{}
			apeArg1 = append(apeArg1, preEomiVowel...)
			apeArg1 = append(apeArg1, preEomi_1_1...)
			apeArg1 = append(apeArg1, preEomi_2...)
			apeArg1 = append(apeArg1, preEomi_6...)

			ape1 := addPreEomi(lastRune, apeArg1)

			for _, d := range ccm {
				expandedLast = append(expandedLast, d.(string))
			}
			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, []string{lastRuneString}...)
		} else
		// Cases with codas
		// 만들다, 알다, 풀다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, v, 'ㄹ'}) &&
			((o == 'ㅁ' && v == 'ㅓ') || v == 'ㅡ' || v == 'ㅏ' || v == 'ㅜ') {
			apeArg1 := []rune{}
			apeArg1 = append(apeArg1, preEomi_1_2...)
			apeArg1 = append(apeArg1, preEomi_3...)

			ape1 := addPreEomi(lastRune, apeArg1)

			apeArg2 := []rune{}
			apeArg2 = append(apeArg2, preEomi_2...)
			apeArg2 = append(apeArg2, preEomi_1_2...)
			apeArg2 = append(apeArg2, preEomiRespect...)

			ape2 := addPreEomi(ComposeHangulFull(o, v, ' '), apeArg2)

			s := []string{
				string(ComposeHangulFull(o, v, 'ㄻ')),
				string(ComposeHangulFull(o, v, 'ㄴ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, ape2...)
			expandedLast = append(expandedLast, s...)
		} else
		// 낫다, 뺴앗다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅏ', 'ㅅ'}) {
			ape1 := getPreEomi_2_6(lastRune)

			apeArg2 := []rune{}
			apeArg2 = append(apeArg2, preEomi_4...)
			apeArg2 = append(apeArg2, preEomi_5...)

			ape2 := addPreEomi(ComposeHangulFull(o, 'ㅏ', ' '), apeArg2)

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, ape2...)
			expandedLast = append(expandedLast, []string{lastRuneString}...)
		} else
		// 묻다
		if isHangulMatch(lastRuneDecomposed, HangulRune{'ㅁ', 'ㅜ', 'ㄷ'}) {
			ape1 := getPreEomi_2_6(lastRune)

			s := []string{
				string(ComposeHangulFull('ㅁ', 'ㅜ', 'ㄹ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, s...)
		} else
		// 붇다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅜ', 'ㄷ'}) {
			ape1 := getPreEomi_2_6(lastRune)

			apeArg2 := []rune{}
			apeArg2 = append(apeArg2, preEomi_1_2...)
			apeArg2 = append(apeArg2, preEomi_1_4...)
			apeArg2 = append(apeArg2, preEomi_4...)
			apeArg2 = append(apeArg2, preEomi_5...)

			ape2 := addPreEomi(ComposeHangulFull(o, 'ㅜ', ' '), apeArg2)

			s := []string{
				string(ComposeHangulFull(o, 'ㅜ', 'ㄹ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, ape2...)
			expandedLast = append(expandedLast, s...)

		} else
		// 눕다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅜ', 'ㅂ'}) {
			ape1 := getPreEomi_2_6(lastRune)

			apeArg2 := []rune{}
			apeArg2 = append(apeArg2, preEomi_1_4...)
			apeArg2 = append(apeArg2, preEomi_4...)
			apeArg2 = append(apeArg2, preEomi_5...)

			ape2 := addPreEomi(ComposeHangulFull(o, 'ㅜ', ' '), apeArg2)

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, ape2...)
			expandedLast = append(expandedLast, []string{lastRuneString}...)
		} else
		// 간지럽다, 갑작스럽다 -> 갑작스런
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅓ', 'ㅂ'}) && isAdjective {
			ape1 := getPreEomi_1_4_7(ComposeHangulFull(o, 'ㅓ', ' '))

			s := []string{
				string(ComposeHangulFull(o, 'ㅓ', ' ')),
				string(ComposeHangulFull(o, 'ㅜ', 'ㄴ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, s...)
		} else
		// 아름답다, 가볍다, 덥다, 간지럽다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, v, 'ㅂ'}) && isAdjective {
			ape1 := getPreEomi_1_4_7(ComposeHangulFull(o, v, ' '))

			s := []string{
				string(ComposeHangulFull(o, v, ' ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, s...)
		} else
		// 놓다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, 'ㅗ', 'ㅎ'}) {
			ape1 := getPreEomi_2_6(lastRune)

			ccm := scala.Map(codasCommon, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, 'ㅗ', ' '))
			}).([]interface{})

			s := []string{
				string(ComposeHangulFull(o, 'ㅘ', ' ')),
				string(ComposeHangulFull(o, 'ㅗ', ' ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(ccm)...)
			expandedLast = append(expandedLast, s...)
		} else
		// 파랗다, 퍼렇다, 어떻다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, v, 'ㅎ'}) && isAdjective {
			ccm := scala.Map(codasCommon, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, v, value.(rune)))
			}).([]interface{})

			cfcm := scala.Map(codasForContraction, func(key interface{}, value interface{}) interface{} {
				return string(ComposeHangulFull(o, 'ㅐ', value.(rune)))
			}).([]interface{})

			s := []string{
				string(ComposeHangulFull(o, 'ㅐ', ' ')),
				string(ComposeHangulFull(o, v, ' ')),
				lastRuneString,
			}

			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(ccm)...)
			expandedLast = append(expandedLast, shared.ConvertSliceOfInterfaceToSliceOfString(cfcm)...)
			expandedLast = append(expandedLast, s...)

		} else
		// 1 char with coda adjective, 있다, 컸다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, v, c}) &&
			(utf8.RuneCountInString(w) == 1 || (isAdjective && c == 'ㅆ')) {
			apeArg1 := []rune{}
			apeArg1 = append(apeArg1, preEomiCommon...)
			apeArg1 = append(apeArg1, preEomi_1_2...)
			apeArg1 = append(apeArg1, preEomi_1_3...)
			apeArg1 = append(apeArg1, preEomi_2...)
			apeArg1 = append(apeArg1, preEomi_4...)
			apeArg1 = append(apeArg1, preEomi_5...)
			apeArg1 = append(apeArg1, preEomi_6...)

			ape1 := addPreEomi(lastRune, apeArg1)

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, []string{lastRuneString}...)
		} else
		// 1 char with coda adjective, 밝다
		if isHangulMatch(lastRuneDecomposed, HangulRune{o, v, c}) &&
			(utf8.RuneCountInString(w) == 1 && isAdjective) {
			apeArg1 := []rune{}
			apeArg1 = append(apeArg1, preEomiCommon...)
			apeArg1 = append(apeArg1, preEomi_1_2...)
			apeArg1 = append(apeArg1, preEomi_1_3...)
			apeArg1 = append(apeArg1, preEomi_2...)
			apeArg1 = append(apeArg1, preEomi_4...)
			apeArg1 = append(apeArg1, preEomi_5...)

			ape1 := addPreEomi(lastRune, apeArg1)

			expandedLast = append(expandedLast, ape1...)
			expandedLast = append(expandedLast, []string{lastRuneString}...)
		} else {
			expandedLast = append(expandedLast, []string{lastRuneString}...)
		}

		return shared.ConvertSliceOfInterfaceToSliceOfString(
			scala.Map(expandedLast, func(key interface{}, value interface{}) interface{} {
				return init + value.(string)
			}).([]interface{}))
	})

	if isAdjective {
		return expanded
	}

	expanded.Remove("아니")
	expanded.Remove("입")
	expanded.Remove("입니")
	expanded.Remove("나는")

	return expanded
}
