package tokenizer

import (
	"strings"
	"tktg/util"
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var SuffixPos data_structure.Set
var PrefixPos data_structure.Set

func init() {
	SuffixPos = data_structure.NewLinkedSet(util.Josa, util.Eomi, util.PreEomi, util.Suffix, util.Punctuation)
	PrefixPos = data_structure.NewLinkedSet(util.NounPrefix, util.VerbPrefix)
}

func Detokenize(input []string) string {
	spaceGuide := getSpaceGuide(input)

	tokenizerProfile := defaultProfile
	tokenizerProfile.spaceGuide = spaceGuide

	tokenized := tokenize(strings.Join(input, ""), tokenizerProfile)

	return strings.Join(collapseTokens(tokenized), " ")
}

type tmpCollapseTokensTuple struct {
	output   []string
	isPrefix bool
}

func collapseTokens(tokenized []KoreanToken) []string {
	t := scala.FoldLeft(tokenized, tmpCollapseTokensTuple{[]string{}, false},
		func(folded interface{}, key interface{}, value interface{}) interface{} {
			output := folded.(tmpCollapseTokensTuple).output
			isPrefix := folded.(tmpCollapseTokensTuple).isPrefix
			token := value.(KoreanToken)

			l := len(output)

			if l > 0 && (isPrefix || SuffixPos.Contains(token.Pos)) {
				attached := output[l-1] + token.Text
				output[l-1] = attached
				x := []string{}

				x = append(x, output[0:l]...)
				return tmpCollapseTokensTuple{x, false}
			} else if PrefixPos.Contains(token.Pos) {
				x := []string{}
				x = append(x, output[0:l]...)
				x = append(x, token.Text)
				return tmpCollapseTokensTuple{x, true}
			} else {
				x := []string{}
				x = append(x, output[0:l]...)
				x = append(x, token.Text)
				return tmpCollapseTokensTuple{x, false}
			}
		}).(tmpCollapseTokensTuple)

	return t.output
}

type tmpSpaceGuideTuple struct {
	spaceGuide data_structure.Set
	index      int
}

func getSpaceGuide(input []string) data_structure.Set {
	t := scala.FoldLeft(input, tmpSpaceGuideTuple{data_structure.NewLinkedSet(), 0},
		func(folded interface{}, key interface{}, value interface{}) interface{} {
			output := folded.(tmpSpaceGuideTuple).spaceGuide
			i := folded.(tmpSpaceGuideTuple).index
			word := value.(string)

			l := i + utf8.RuneCountInString(word)

			s := data_structure.NewLinkedSet(output.ToSlice()...)
			s.Add(l)

			return tmpSpaceGuideTuple{s, l}
		}).(tmpSpaceGuideTuple)

	return t.spaceGuide
}
