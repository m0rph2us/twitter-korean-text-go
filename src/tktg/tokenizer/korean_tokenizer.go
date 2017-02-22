package tokenizer

import (
	"reflect"
	"sort"
	"tktg/util"
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var topNPerState int = 5
var maxTraceback int = 8

type KoreanToken struct {
	Text    string
	Pos     util.KoreanPos
	Offset  int
	Length  int
	Unknown bool
}

func (t KoreanToken) toString() string {
	unknownStar := ""
	if t.Unknown {
		unknownStar = "*"
	}

	return t.Text + unknownStar + "(" + string(int(t.Pos)) + ": " + string(t.Offset) + "," + string(t.Length) + ")"
}

func (t KoreanToken) copyWithNewPos(pos util.KoreanPos) KoreanToken {
	return KoreanToken{
		t.Text, pos, t.Offset, t.Length, t.Unknown,
	}
}

type CandidateParse struct {
	parse   ParsedChunk
	curTrie []util.KoreanPosTrie
	ending  util.KoreanPos // Option
	order   int
}

type PossibleTrie struct {
	curTrie util.KoreanPosTrie
	words   int
}

var sequenceDefinition data_structure.Map

var koreanPosTrie []util.KoreanPosTrie

func init() {
	sequenceDefinition = data_structure.NewLinkedMap()
	sequenceDefinition.Put("D0p*N1s0j0", util.Noun)
	sequenceDefinition.Put("v*V1r*e0", util.Verb)
	sequenceDefinition.Put("v*J1r*e0", util.Adjective)
	sequenceDefinition.Put("A1", util.Adverb)
	sequenceDefinition.Put("C1", util.Conjunction)
	sequenceDefinition.Put("E+", util.Exclamation)
	sequenceDefinition.Put("j1", util.Josa)

	koreanPosTrie = util.GetTrie(sequenceDefinition)
}

func Tokenize(text string) []KoreanToken {
	return tokenize(text, defaultProfile)
}

func tokenize(text string, profile TokenizerProfile) []KoreanToken {
	flatMap := func(value []KoreanToken, f func(value interface{}) interface{}) []KoreanToken {
		s := []KoreanToken{}
		for _, v := range value {
			x := f(v)

			to := reflect.TypeOf(x)

			if to.Kind() == reflect.Array || to.Kind() == reflect.Slice {
				s = append(s, x.([]KoreanToken)...)
			}
		}
		return s
	}

	return flatMap(chunk(text),
		func(value interface{}) interface{} {
			if value.(KoreanToken).Pos == util.Korean {
				parsed := parseKoreanChunk(value.(KoreanToken), profile)
				return collapseNouns(parsed)
			}
			return []KoreanToken{value.(KoreanToken)}
		},
	)
}

type ByScoreAndTieBreaker []CandidateParse

func (t ByScoreAndTieBreaker) Len() int {
	return len(t)
}

func (t ByScoreAndTieBreaker) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByScoreAndTieBreaker) Less(i, j int) bool {
	if t[i].parse.score < t[j].parse.score {
		return true
	}
	if t[i].parse.score > t[j].parse.score {
		return false
	}
	if t[i].parse.posTieBreaker() < t[j].parse.posTieBreaker() {
		return true
	}
	if t[i].parse.posTieBreaker() > t[j].parse.posTieBreaker() {
		return false
	}
	return t[i].order < t[j].order
}

func parseKoreanChunk(chunk KoreanToken, profile TokenizerProfile) []KoreanToken {
	for v := range util.KoreanDictionary.Iter() {
		if v.Value.(data_structure.Set).Contains(chunk.Text) {
			return []KoreanToken{
				{
					chunk.Text,
					v.Key.(util.KoreanPos),
					chunk.Offset,
					chunk.Length,
					false,
				},
			}
		}
	}

	var solutions map[int][]CandidateParse = map[int][]CandidateParse{
		0: {
			{
				NewParsedChunk([]KoreanToken{}, 1, profile),
				koreanPosTrie,
				util.None,
				0,
			},
		},
	}

	flatMap := func(value []CandidateParse, f func(value interface{}) interface{}) []CandidateParse {
		s := []CandidateParse{}
		for _, v := range value {
			x := f(v)
			s = append(s, x.([]CandidateParse)...)
		}
		return s
	}

	for end := 1; end <= chunk.Length; end++ {
		to := 0
		c := end - maxTraceback
		if to < c {
			to = c
		}

		for start := end - 1; start >= to; start-- {
			word := string([]rune(chunk.Text)[start:end])
			curSolutions := solutions[start]

			candidates := flatMap(curSolutions, func(value interface{}) interface{} {
				solution := value.(CandidateParse)

				var possiblePoses []PossibleTrie

				if solution.ending != util.None {
					a := ConvertSliceOfInterfaceToSliceOfPossibleTrie(scala.Map(solution.curTrie,
						func(key interface{}, value interface{}) interface{} {
							return PossibleTrie{value.(util.KoreanPosTrie), 0}
						},
					).([]interface{}))

					b := ConvertSliceOfInterfaceToSliceOfPossibleTrie(scala.Map(koreanPosTrie,
						func(key interface{}, value interface{}) interface{} {
							return PossibleTrie{value.(util.KoreanPosTrie), 1}
						},
					).([]interface{}))

					possiblePoses = append(a, b...)
				} else {
					possiblePoses = ConvertSliceOfInterfaceToSliceOfPossibleTrie(
						scala.Map(solution.curTrie,
							func(key interface{}, value interface{}) interface{} {
								return PossibleTrie{value.(util.KoreanPosTrie), 0}
							},
						).([]interface{}))
				}

				d := ConvertSliceOfInterfaceToSliceOfPossibleTrie(
					scala.Filter(possiblePoses, func(value interface{}) bool {
						v := value.(PossibleTrie)
						return v.curTrie.CurPos == util.Noun ||
							util.KoreanDictionary.Get(v.curTrie.CurPos).(data_structure.Set).Contains(word)
					}).([]interface{}))

				return ConvertSliceOfInterfaceToSliceOfCandidateParse(
					scala.Map(d, func(key interface{}, value interface{}) interface{} {
						v := value.(PossibleTrie)

						var candidateToAdd ParsedChunk

						if v.curTrie.CurPos == util.Noun &&
							!util.KoreanDictionary.Get(util.Noun).(data_structure.Set).Contains(word) {
							isWordName := util.IsName(word)
							isWordKoreanNameVariation := util.IsKoreanNameVariation(word)

							unknown := !isWordName &&
								!util.IsKoreanNumber(word) &&
								!isWordKoreanNameVariation

							var pos util.KoreanPos
							if unknown || isWordName || isWordKoreanNameVariation {
								pos = util.ProperNoun
							} else {
								pos = util.Noun
							}

							candidateToAdd = NewParsedChunk(
								[]KoreanToken{
									{
										word,
										pos,
										chunk.Offset + start,
										utf8.RuneCountInString(word),
										unknown,
									},
								},
								v.words,
								profile)
						} else {
							var pos util.KoreanPos
							if v.curTrie.CurPos == util.Noun && util.ProperNouns.Contains(word) {
								pos = util.ProperNoun
							} else {
								pos = v.curTrie.CurPos
							}

							candidateToAdd = NewParsedChunk(
								[]KoreanToken{
									{
										word,
										pos,
										chunk.Offset + start,
										utf8.RuneCountInString(word),
										false,
									},
								},
								v.words,
								profile)
						}

						nextTrie := ConvertSliceOfInterfaceToSliceOfKoreanPosTrie(
							scala.Map(v.curTrie.NextTrie,
								func(key interface{}, value interface{}) interface{} {
									if value.(util.KoreanPosTrie).Compare(util.SelfNode) {
										return v.curTrie
									}
									return value.(util.KoreanPosTrie)
								},
							).([]interface{}))

						return CandidateParse{
							solution.parse.add(candidateToAdd),
							nextTrie,
							v.curTrie.Ending,
							0,
						}
					}).([]interface{}))

			})

			currentSolutions := []CandidateParse{}

			if v, e := solutions[end]; e == true {
				currentSolutions = v
			}

			currentSolutions = append(currentSolutions, candidates...)

			// to avoid sorting problem
			for i, _ := range currentSolutions {
				currentSolutions[i].order = i
			}

			sort.Sort(ByScoreAndTieBreaker(currentSolutions))

			min := len(currentSolutions)
			if topNPerState < min {
				min = topNPerState
			}

			solutions[end] = currentSolutions[0:min]
		}
	}

	if len(solutions[chunk.Length]) == 0 {
		return []KoreanToken{{chunk.Text, util.Noun, 0, chunk.Length, true}}
	} else {
		prevScore := 0.0
		minIndex := -1

		for i, v := range solutions[chunk.Length] {
			if minIndex == -1 {
				minIndex = i
			} else {
				if prevScore > v.parse.score {
					minIndex = i
				}
			}
			prevScore = v.parse.score
		}

		return solutions[chunk.Length][minIndex].parse.posNodes
	}
}

type tmpCollapseNounsTuple struct {
	nodes      []KoreanToken
	collapsing bool
}

func collapseNouns(posNodes []KoreanToken) []KoreanToken {
	t := scala.FoldLeft(posNodes, tmpCollapseNounsTuple{[]KoreanToken{}, false},
		func(folded interface{}, key interface{}, value interface{}) interface{} {
			t := folded.(tmpCollapseNounsTuple)
			v := value.(KoreanToken)

			if v.Pos == util.Noun && utf8.RuneCountInString(v.Text) == 1 && t.collapsing {
				text := t.nodes[0].Text + v.Text
				offset := t.nodes[0].Offset

				x := []KoreanToken{}
				x = append(x, KoreanToken{text, util.Noun, offset, utf8.RuneCountInString(text), true})
				if len(t.nodes) > 1 {
					x = append(x, t.nodes[len(t.nodes)-1:]...)
				}

				return tmpCollapseNounsTuple{x, true}
			} else if v.Pos == util.Noun && utf8.RuneCountInString(v.Text) == 1 && !t.collapsing {
				x := []KoreanToken{}
				x = append(x, v)
				x = append(x, t.nodes...)

				return tmpCollapseNounsTuple{x, true}
			} else {
				x := []KoreanToken{}
				x = append(x, v)
				x = append(x, t.nodes...)

				return tmpCollapseNounsTuple{x, false}
			}
			return nil
		}).(tmpCollapseNounsTuple)

	return ConvertSliceOfInterfaceToSliceOfKoreanToken(scala.Reverse(t.nodes).([]interface{}))
}
