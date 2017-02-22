package phrase_extractor

import (
	"fmt"
	"reflect"
	"strings"
	"tktg/shared"
	"tktg/tokenizer"
	"tktg/util"
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var minCharsPerPhraseChunkWithoutSpaces = 2
var minPhrasesPerPhraseChunk = 3

var maxCharsPerPhraseChunkWithoutSpaces = 30
var maxPhrasesPerPhraseChunk = 8

var modifyingPredicateEndings data_structure.Set = data_structure.NewLinkedSet('ㄹ', 'ㄴ')
var modifyingPredicateExceptions data_structure.Set = data_structure.NewLinkedSet('만')

var phraseTokens data_structure.Set = data_structure.NewLinkedSet(util.Noun, util.ProperNoun, util.Space)
var conjunctionJosa data_structure.Set = data_structure.NewLinkedSet("와", "과", "의")

type KoreanPhrase struct {
	Tokens []tokenizer.KoreanToken
	Pos    util.KoreanPos
}

func (t *KoreanPhrase) offset() int {
	return t.Tokens[0].Offset
}

func (t *KoreanPhrase) text() string {
	v := shared.ConvertSliceOfInterfaceToSliceOfString(
		scala.Map(t.Tokens, func(key interface{}, value interface{}) interface{} {
			v := value.(tokenizer.KoreanToken)
			return v.Text
		}).([]interface{}))

	return strings.Join(v, "")
}

func (t *KoreanPhrase) length() int {
	v := scala.Map(t.Tokens, func(key interface{}, value interface{}) interface{} {
		v := value.(tokenizer.KoreanToken)
		return utf8.RuneCountInString(v.Text)
	}).([]interface{})

	return scala.FoldLeft(v, 0, func(folded interface{}, key interface{}, value interface{}) interface{} {
		return folded.(int) + value.(int)
	}).(int)
}

type PhraseBuffer struct {
	phrases []KoreanPhrase
	curTrie []util.KoreanPosTrie
	ending  util.KoreanPos
}

type KoreanPhraseChunk []KoreanPhrase

var phraseHeadPoses data_structure.Set = data_structure.NewLinkedSet(
	util.Adjective, util.Noun, util.ProperNoun, util.Alpha, util.Number)

var phraseTailPoses data_structure.Set = data_structure.NewLinkedSet(
	util.Noun, util.ProperNoun, util.Alpha, util.Number)

var collapsingRules data_structure.Map

var collapseTrie []util.KoreanPosTrie

func init() {
	collapsingRules = data_structure.NewLinkedMap()
	collapsingRules.Put("D0p*N1s0", util.Noun)
	collapsingRules.Put("n*a+n*", util.Noun)
	collapsingRules.Put("n+", util.Noun)
	collapsingRules.Put("v*V1r*e0", util.Verb)
	collapsingRules.Put("v*J1r*e0", util.Adjective)

	collapseTrie = util.GetTrie(collapsingRules)
}

func trimPhraseChunk(phrases KoreanPhraseChunk) KoreanPhraseChunk {
	a := scala.DropWhile(phrases, func(value interface{}) bool {
		return !phraseHeadPoses.Contains(value.(KoreanPhrase).Pos)
	})

	b := scala.Reverse(a)

	c := scala.DropWhile(b, func(value interface{}) bool {
		return !phraseTailPoses.Contains(value.(KoreanPhrase).Pos)
	})

	trimNonNouns := ConvertSliceOfInterfaceToSliceOfKoreanPhrase(scala.Reverse(c).([]interface{}))

	trimSpacesFromPhrase := func(phrases []KoreanPhrase) []KoreanPhrase {
		return ConvertSliceOfInterfaceToSliceOfKoreanPhrase(
			scala.Map(
				scala.ZipWithIndex(phrases), func(key interface{}, value interface{}) interface{} {
					phrase := value.([]interface{})[0].(KoreanPhrase)
					i := value.([]interface{})[1].(int)

					if len(phrases) == 1 {
						return trimPhrase(phrase)
					} else if i == 0 {
						a := tokenizer.ConvertSliceOfInterfaceToSliceOfKoreanToken(
							scala.DropWhile(phrase.Tokens, func(x interface{}) bool {
								return x.(tokenizer.KoreanToken).Pos == util.Space
							}).([]interface{}))

						return KoreanPhrase{a, phrase.Pos}
					} else if i == len(phrases)-1 {
						a := scala.Reverse(phrase.Tokens)

						b := scala.DropWhile(a, func(x interface{}) bool {
							return x.(tokenizer.KoreanToken).Pos == util.Space
						}).([]interface{})

						c := tokenizer.ConvertSliceOfInterfaceToSliceOfKoreanToken(
							scala.Reverse(b).([]interface{}))
						return KoreanPhrase{c, phrase.Pos}
					} else {
						return phrase
					}
				}).([]interface{}))
	}

	return KoreanPhraseChunk(trimSpacesFromPhrase(trimNonNouns))
}

func trimPhrase(phrase KoreanPhrase) KoreanPhrase {
	a := scala.DropWhile(phrase.Tokens, func(value interface{}) bool {
		return value.(tokenizer.KoreanToken).Pos == util.Space
	})

	b := scala.Reverse(a)

	c := scala.DropWhile(b, func(value interface{}) bool {
		return value.(tokenizer.KoreanToken).Pos == util.Space
	})

	d := tokenizer.ConvertSliceOfInterfaceToSliceOfKoreanToken(scala.Reverse(c).([]interface{}))

	return KoreanPhrase{d, phrase.Pos}
}

func isProperPhraseChunk(phraseChunk KoreanPhraseChunk) bool {
	if len(phraseChunk) == 0 {
		return false
	}

	c := phraseChunk[len(phraseChunk)-1]

	if len(c.Tokens) == 0 {
		return false
	}

	lastToken := c.Tokens[len(c.Tokens)-1]
	notEndingInNonPhraseSuffix := !(lastToken.Pos == util.Suffix && lastToken.Text == "적")

	phraseChunkWithoutSpaces := ConvertSliceOfInterfaceToSliceOfKoreanPhrase(
		scala.Filter(phraseChunk, func(value interface{}) bool {
			return value.(KoreanPhrase).Pos != util.Space
		}).([]interface{}))

	a := scala.Map(phraseChunkWithoutSpaces, func(key interface{}, value interface{}) interface{} {
		a := value.(KoreanPhrase)
		return (&a).length()
	})

	b := scala.FoldLeft(a, 0, func(folded interface{}, key interface{}, value interface{}) interface{} {
		return folded.(int) + value.(int)
	}).(int)

	checkMaxLength := len(phraseChunkWithoutSpaces) <= maxPhrasesPerPhraseChunk &&
		b <= maxCharsPerPhraseChunkWithoutSpaces

	checkMinLength := len(phraseChunkWithoutSpaces) >= minPhrasesPerPhraseChunk ||
		(len(phraseChunkWithoutSpaces) < minPhrasesPerPhraseChunk &&
			b >= minCharsPerPhraseChunkWithoutSpaces)

	checkMinLengthPerToken := scala.Exists(phraseChunkWithoutSpaces, func(value interface{}) bool {
		a := value.(KoreanPhrase)
		return (&a).length() > 1
	})

	isRightLength := checkMaxLength && checkMinLength && checkMinLengthPerToken

	return isRightLength && notEndingInNonPhraseSuffix
}

type tmpTrieTuple struct {
	curTrie  util.KoreanPosTrie
	nextTrie []util.KoreanPosTrie
}

func collapsePos(tokens []tokenizer.KoreanToken) []KoreanPhrase {
	getTries := func(token tokenizer.KoreanToken, trie []util.KoreanPosTrie) tmpTrieTuple {
		a := tokenizer.ConvertSliceOfInterfaceToSliceOfKoreanPosTrie(
			scala.Filter(trie, func(value interface{}) bool {
				return value.(util.KoreanPosTrie).CurPos == token.Pos
			}).([]interface{}))

		curTrie := a[0]

		nextTrie := tokenizer.ConvertSliceOfInterfaceToSliceOfKoreanPosTrie(
			scala.Map(curTrie.NextTrie, func(key interface{}, value interface{}) interface{} {
				v := value.(util.KoreanPosTrie)
				if v.Compare(util.SelfNode) {
					return curTrie
				} else {
					return v
				}
			}).([]interface{}))

		return tmpTrieTuple{curTrie, nextTrie}
	}

	getInit := func(phraseBuffer PhraseBuffer) []KoreanPhrase {
		if len(phraseBuffer.phrases) == 0 {
			return []KoreanPhrase{}
		} else {
			x := []KoreanPhrase{}
			return append(x, phraseBuffer.phrases[0:len(phraseBuffer.phrases)-1]...)
		}
	}

	return scala.FoldLeft(tokens, PhraseBuffer{[]KoreanPhrase{}, collapseTrie, util.None},
		func(folded interface{}, key interface{}, value interface{}) interface{} {
			output := folded.(PhraseBuffer)
			token := value.(tokenizer.KoreanToken)

			if scala.Exists(output.curTrie, func(value interface{}) bool {
				return value.(util.KoreanPosTrie).CurPos == token.Pos
			}) {
				t := getTries(token, output.curTrie)
				ct := t.curTrie
				nt := t.nextTrie

				if len(output.phrases) == 0 ||
					fmt.Sprintf("%v", output.curTrie) == fmt.Sprintf("%v", collapseTrie) {
					a := []KoreanPhrase{}
					a = append(a, output.phrases...)
					a = append(a, KoreanPhrase{
						[]tokenizer.KoreanToken{token},
						ct.Ending,
					})

					return PhraseBuffer{a, nt, ct.Ending}
				} else {
					a := []KoreanPhrase{}
					a = append(a, getInit(output)...)

					b := []tokenizer.KoreanToken{}
					b = append(b, output.phrases[len(output.phrases)-1].Tokens...)
					b = append(b, token)

					a = append(a, KoreanPhrase{
						b, ct.Ending,
					})

					return PhraseBuffer{a, nt, ct.Ending}
				}
			} else if scala.Exists(collapseTrie, func(value interface{}) bool {
				return value.(util.KoreanPosTrie).CurPos == token.Pos
			}) {
				t := getTries(token, collapseTrie)
				ct := t.curTrie
				nt := t.nextTrie

				a := []KoreanPhrase{}
				a = append(a, output.phrases...)
				a = append(a, KoreanPhrase{
					[]tokenizer.KoreanToken{token},
					ct.Ending,
				})

				return PhraseBuffer{a, nt, ct.Ending}
			} else {
				a := []KoreanPhrase{}
				a = append(a, output.phrases...)
				a = append(a, KoreanPhrase{
					[]tokenizer.KoreanToken{token},
					token.Pos,
				})
				return PhraseBuffer{a, collapseTrie, output.ending}
			}
		}).(PhraseBuffer).phrases
}

type tmpPhraseChunkTuple struct {
	l      []KoreanPhraseChunk
	buffer data_structure.Set
}

func distinctPhrases(chunks []KoreanPhraseChunk) []KoreanPhraseChunk {
	t := scala.FoldLeft(chunks, tmpPhraseChunkTuple{[]KoreanPhraseChunk{}, data_structure.NewLinkedSet()},
		func(folded interface{}, key interface{}, value interface{}) interface{} {
			l := folded.(tmpPhraseChunkTuple).l
			buffer := folded.(tmpPhraseChunkTuple).buffer
			chunk := value.(KoreanPhraseChunk)

			phraseText := strings.Join(
				shared.ConvertSliceOfInterfaceToSliceOfString(
					scala.Map(chunk, func(key interface{}, value interface{}) interface{} {
						return strings.Join(shared.ConvertSliceOfInterfaceToSliceOfString(
							scala.Map(value.(KoreanPhrase).Tokens,
								func(key interface{}, value interface{}) interface{} {
									return value.(tokenizer.KoreanToken).Text
								}).([]interface{})), "")
					}).([]interface{})), "")

			if buffer.Contains(phraseText) {
				return tmpPhraseChunkTuple{l, buffer}
			} else {
				a := []KoreanPhraseChunk{}
				a = append(a, chunk)
				a = append(a, l...)

				b := data_structure.NewLinkedSet(buffer.ToSlice()...)
				b.Add(phraseText)

				return tmpPhraseChunkTuple{a, b}
			}
		}).(tmpPhraseChunkTuple)

	return ConvertSliceOfInterfaceToSliceOfKoreanPhraseChunk(scala.Reverse(t.l).([]interface{}))
}

type tmpCollpaseNounPhraseTuple struct {
	output []KoreanPhrase
	buffer []KoreanPhrase
}

type tmpNewBufferTuple struct {
	output []KoreanPhraseChunk
	buffer []KoreanPhraseChunk
}

func getCandidatePhraseChunks(phrases KoreanPhraseChunk, filterSpam bool /* default false */) []KoreanPhraseChunk {
	isNotSpam := func(phrase KoreanPhrase) bool {
		return !filterSpam || !scala.Exists(phrase.Tokens, func(value interface{}) bool {
			return util.SpamNouns.Contains(value.(tokenizer.KoreanToken).Text)
		})
	}

	isNonNounPhraseCondidate := func(phrase KoreanPhrase) bool {
		trimmed := trimPhrase(phrase)

		t := trimmed.Tokens[len(trimmed.Tokens)-1].Text
		a := []rune(t)
		lastRune := a[len(a)-1]

		isModifyingPredicate := (trimmed.Pos == util.Verb || trimmed.Pos == util.Adjective) &&
			modifyingPredicateEndings.Contains(util.DecomposeHangul(lastRune).Coda) &&
			!modifyingPredicateExceptions.Contains(lastRune)

		isConjunction := trimmed.Pos == util.Josa && conjunctionJosa.Contains(t)

		isAlphaNumeric := trimmed.Pos == util.Alpha || trimmed.Pos == util.Number

		return isAlphaNumeric || isModifyingPredicate || isConjunction
	}

	collapseNounPhrases := func(phrases KoreanPhraseChunk) KoreanPhraseChunk {
		flatMap := func(value []KoreanPhrase, f func(value interface{}) interface{}) []tokenizer.KoreanToken {
			s := data_structure.NewLinkedSet()
			for _, v := range value {
				x := f(v)

				to := reflect.TypeOf(x)

				if to.Kind() == reflect.Array || to.Kind() == reflect.Slice {
					for _, d := range x.([]tokenizer.KoreanToken) {
						s.Add(d)
					}
				}
			}
			return tokenizer.ConvertSliceOfInterfaceToSliceOfKoreanToken(s.ToSlice())
		}

		t := scala.FoldLeft(phrases, tmpCollpaseNounPhraseTuple{[]KoreanPhrase{}, []KoreanPhrase{}},
			func(folded interface{}, key interface{}, value interface{}) interface{} {
				output := folded.(tmpCollpaseNounPhraseTuple).output
				buffer := folded.(tmpCollpaseNounPhraseTuple).buffer
				phrase := value.(KoreanPhrase)

				if phrase.Pos == util.Noun || phrase.Pos == util.ProperNoun {
					a := []KoreanPhrase{}
					a = append(a, buffer...)
					a = append(a, phrase)

					return tmpCollpaseNounPhraseTuple{
						output,
						a,
					}
				} else {
					tmpPhrases := []KoreanPhrase{phrase}
					if len(buffer) > 0 {
						tmpPhrases = []KoreanPhrase{
							{
								flatMap(buffer, func(value interface{}) interface{} {
									return value.(KoreanPhrase).Tokens
								}),
								util.Noun,
							},
							phrase,
						}
					}

					a := []KoreanPhrase{}
					a = append(a, output...)
					a = append(a, tmpPhrases...)

					return tmpCollpaseNounPhraseTuple{a, []KoreanPhrase{}}
				}
			}).(tmpCollpaseNounPhraseTuple)

		output := t.output
		buffer := t.buffer

		if len(buffer) > 0 {
			a := []KoreanPhrase{}
			a = append(a, output...)
			a = append(a, KoreanPhrase{
				flatMap(buffer, func(value interface{}) interface{} {
					return value.(KoreanPhrase).Tokens
				}),
				util.Noun,
			})
			return a
		} else {
			return output
		}
	}

	collapsePhrases := func(phrases KoreanPhraseChunk) []KoreanPhraseChunk {
		addPhraseToBuffer := func(phrase KoreanPhrase, buffer []KoreanPhraseChunk) []KoreanPhraseChunk {
			return ConvertSliceOfInterfaceToSliceOfKoreanPhraseChunk(
				scala.Map(buffer, func(key interface{}, value interface{}) interface{} {
					a := KoreanPhraseChunk{}
					a = append(a, value.(KoreanPhraseChunk)...)
					a = append(a, phrase)

					return a
				}).([]interface{}))
		}

		t := scala.FoldLeft(phrases, tmpNewBufferTuple{[]KoreanPhraseChunk{}, []KoreanPhraseChunk{{}}},
			func(folded interface{}, key interface{}, value interface{}) interface{} {
				output := folded.(tmpNewBufferTuple).output
				buffer := folded.(tmpNewBufferTuple).buffer
				phrase := value.(KoreanPhrase)

				if phraseTokens.Contains(phrase.Pos) && isNotSpam(phrase) {
					bufferWithThisPhrase := addPhraseToBuffer(phrase, buffer)

					if phrase.Pos == util.Noun || phrase.Pos == util.ProperNoun {
						a := []KoreanPhraseChunk{}
						a = append(a, output...)
						a = append(a, bufferWithThisPhrase...)
						return tmpNewBufferTuple{a, bufferWithThisPhrase}
					} else {
						return tmpNewBufferTuple{output, bufferWithThisPhrase}
					}
				} else if isNonNounPhraseCondidate(phrase) {
					return tmpNewBufferTuple{output, addPhraseToBuffer(phrase, buffer)}
				} else {
					a := []KoreanPhraseChunk{}
					a = append(a, output...)
					a = append(a, buffer...)

					return tmpNewBufferTuple{a, []KoreanPhraseChunk{{}}}
				}
			}).(tmpNewBufferTuple)

		output := t.output
		buffer := t.buffer

		if len(buffer) > 0 {
			a := []KoreanPhraseChunk{}
			a = append(a, output...)
			a = append(a, buffer...)
			return a
		} else {
			return output
		}
	}

	a := scala.Filter(phrases, func(value interface{}) bool {
		phrase := value.(KoreanPhrase)

		trimmed := trimPhrase(phrase)

		return (phrase.Pos == util.Noun || phrase.Pos == util.ProperNoun) && isNotSpam(phrase) &&
			((&trimmed).length() >= minCharsPerPhraseChunkWithoutSpaces || len(trimmed.Tokens) >= minPhrasesPerPhraseChunk)
	})

	getSingleTokenNouns := ConvertSliceOfInterfaceToSliceOfKoreanPhraseChunk(
		scala.Map(a, func(key interface{}, value interface{}) interface{} {
			phrase := value.(KoreanPhrase)
			return KoreanPhraseChunk([]KoreanPhrase{trimPhrase(phrase)})
		}).([]interface{}))

	nounPhrases := collapseNounPhrases(phrases)

	phraseCollapsed := collapsePhrases(nounPhrases)

	b := scala.Map(phraseCollapsed, func(key interface{}, value interface{}) interface{} {
		return trimPhraseChunk(value.(KoreanPhraseChunk))
	}).([]interface{})

	c := ConvertSliceOfInterfaceToSliceOfKoreanPhraseChunk(b)

	c = append(c, getSingleTokenNouns...)

	return distinctPhrases(c)
}

func ExtractPhrases(tokens []tokenizer.KoreanToken, filterSpam bool, /* default false */
	addHashtags bool /* default true */) []KoreanPhrase {

	flatMap := func(value []tokenizer.KoreanToken, f func(value interface{}) interface{}) []KoreanPhrase {
		s := []KoreanPhrase{}
		for _, v := range value {
			x := f(v)

			to := reflect.TypeOf(x)

			if to.Kind() == reflect.Array || to.Kind() == reflect.Slice {
				s = append(s, x.([]KoreanPhrase)...)
			}
		}
		return s
	}

	hashtags := flatMap(tokens, func(value interface{}) interface{} {
		t := value.(tokenizer.KoreanToken)
		if t.Pos == util.Hashtag {
			return KoreanPhrase{[]tokenizer.KoreanToken{t}, util.Hashtag}
		} else if t.Pos == util.CashTag {
			return KoreanPhrase{[]tokenizer.KoreanToken{t}, util.CashTag}
		} else {
			return util.None
		}
	})

	collapsed := collapsePos(tokens)

	candidates := getCandidatePhraseChunks(collapsed, filterSpam)

	permutatedCandidates := permutateCadidates(candidates)

	flatMap2 := func(value KoreanPhraseChunk, f func(value interface{}) interface{}) []tokenizer.KoreanToken {
		s := []tokenizer.KoreanToken{}
		for _, v := range value {
			x := f(v)

			to := reflect.TypeOf(x)

			if to.Kind() == reflect.Array || to.Kind() == reflect.Slice {
				s = append(s, x.([]tokenizer.KoreanToken)...)
			}
		}
		return s
	}

	phrases := ConvertSliceOfInterfaceToSliceOfKoreanPhrase(
		scala.Map(permutatedCandidates, func(key interface{}, value interface{}) interface{} {
			phraseChunk := value.(KoreanPhraseChunk)
			return KoreanPhrase{
				flatMap2(
					trimPhraseChunk(phraseChunk),
					func(value interface{}) interface{} {
						return value.(KoreanPhrase).Tokens
					}),
				util.Noun,
			}
		}).([]interface{}))

	if addHashtags {
		a := []KoreanPhrase{}
		a = append(a, phrases...)
		a = append(a, hashtags...)

		return a
	} else {
		return phrases
	}
}

func permutateCadidates(candidates []KoreanPhraseChunk) []KoreanPhraseChunk {
	flatMap := func(value []KoreanPhraseChunk, f func(value interface{}) interface{}) []KoreanPhraseChunk {
		s := []KoreanPhraseChunk{}
		for _, v := range value {
			x := f(v)

			to := reflect.TypeOf(x)

			if to.Kind() == reflect.Array || to.Kind() == reflect.Slice {
				s = append(s, x.([]KoreanPhraseChunk)...)
			}
		}
		return s
	}

	a := flatMap(candidates, func(value interface{}) interface{} {
		phrases := value.(KoreanPhraseChunk)

		if len(phrases) > minPhrasesPerPhraseChunk {
			a := []int{}
			for i := 0; len(phrases)-minPhrasesPerPhraseChunk >= i; i++ {
				a = append(a, i)
			}

			return ConvertSliceOfInterfaceToSliceOfKoreanPhraseChunk(
				scala.Map(a, func(key interface{}, value interface{}) interface{} {
					i := value.(int)
					return trimPhraseChunk(phrases[i:])
				}).([]interface{}))
		} else {
			return []KoreanPhraseChunk{phrases}
		}
	})

	permutated := ConvertSliceOfInterfaceToSliceOfKoreanPhraseChunk(
		scala.Filter(a, func(value interface{}) bool {
			phraseChunk := value.(KoreanPhraseChunk)
			return isProperPhraseChunk(phraseChunk)
		}).([]interface{}))

	return distinctPhrases(permutated)
}
