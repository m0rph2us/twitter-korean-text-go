package tokenizer

import (
	"tktg/util"
	"unicode/utf8"
	"tktg/data_structure"
	"github.com/acidd15/go-scala-util/src/scala"
)

var suffixes data_structure.Set

var preferredBeforeHaVerb data_structure.Set

type ParsedChunk struct {
	posNodes []KoreanToken
	words    int
	profile  TokenizerProfile
	score    float64
}

func init() {
	suffixes = data_structure.NewLinkedSet(util.Suffix, util.Eomi, util.Josa, util.PreEomi)
	preferredBeforeHaVerb = data_structure.NewLinkedSet(util.Noun, util.ProperNoun, util.VerbPrefix)
}

func NewParsedChunk(posNodes []KoreanToken, words int, profile TokenizerProfile) ParsedChunk {
	p := ParsedChunk{
		posNodes,
		words,
		profile,
		0.0,
	}

	(&p).init()

	return p
}

func (t *ParsedChunk) init() {
	t.score = float64(t.countTokens())*float64(t.profile.tokenCount) +
		float64(t.countUnknowns())*float64(t.profile.unknown) +
		float64(t.words)*float64(t.profile.wordCount) +
		float64(t.getUnknownCoverage())*float64(t.profile.unknownCoverage) +
		float64(t.getFreqScore())*float64(t.profile.freq) +
		float64(t.countPos(util.Unknown))*float64(t.profile.unknownPosCount) +
		float64(t.isExactMatch())*float64(t.profile.exactMatch) +
		float64(t.isAllNouns())*float64(t.profile.allNoun) +
		float64(t.isPreferredPattern())*float64(t.profile.preferredPattern) +
		float64(t.countPos(util.Determiner))*float64(t.profile.determinerPosCount) +
		float64(t.countPos(util.Exclamation))*float64(t.profile.exclamationPosCount) +
		float64(t.isInitialPostPosition())*float64(t.profile.initialPostPosition) +
		float64(t.isNounHa())*float64(t.profile.haVerb) +
		float64(t.hasSpaceOutOfGuide())*float64(t.profile.spaceGuidePenalty)
}

func (t *ParsedChunk) countUnknowns() int {
	return scala.Count(t.posNodes, func(value interface{}) bool {
		return value.(KoreanToken).Unknown
	})
}

func (t *ParsedChunk) countTokens() int {
	return len(t.posNodes)
}

func (t *ParsedChunk) isInitialPostPosition() int {
	if len(t.posNodes) > 0 && suffixes.Contains(t.posNodes[0].Pos) {
		return 1
	}
	return 0
}

func (t *ParsedChunk) isExactMatch() int {
	if len(t.posNodes) == 1 {
		return 0
	}
	return 1
}

func (t *ParsedChunk) hasSpaceOutOfGuide() int {
	if t.profile.spaceGuide.Size() == 0 {
		return 0
	} else {
		v := scala.Filter(t.posNodes, func(value interface{}) bool {
			return !suffixes.Contains(value.(KoreanToken).Pos)
		}).([]interface{})

		return scala.Count(v, func(value interface{}) bool {
			return !t.profile.spaceGuide.Contains(value.(KoreanToken).Offset)
		})
	}
}

func (t *ParsedChunk) isAllNouns() int {
	if scala.Exists(t.posNodes, func(value interface{}) bool {
		return value.(KoreanToken).Pos != util.Noun && value.(KoreanToken).Pos != util.ProperNoun
	}) {
		return 1
	}
	return 0
}

func (t *ParsedChunk) isPreferredPattern() int {
	mappedPosNodes := scala.Map(t.posNodes, func(key interface{}, value interface{}) interface{} {
		return value.(KoreanToken).Pos
	}).([]interface{})

	isExist := scala.Exists(t.profile.preferredPatterns, func(value interface{}) bool {
		if len(mappedPosNodes) == 2 &&
			mappedPosNodes[0].(util.KoreanPos) == value.([]interface{})[0].(util.KoreanPos) &&
			mappedPosNodes[1].(util.KoreanPos) == value.([]interface{})[1].(util.KoreanPos) {
			return true
		}
		return false
	})

	if len(t.posNodes) == 2 && isExist {
		return 0
	}
	return 1
}

func (t *ParsedChunk) isNounHa() int {
	if len(t.posNodes) >= 2 &&
		preferredBeforeHaVerb.Contains(t.posNodes[0].Pos) &&
		t.posNodes[1].Pos == util.Verb &&
		([]rune(t.posNodes[1].Text))[0] == '하' {
		return 0
	}
	return 1
}

func (t *ParsedChunk) posTieBreaker() int {
	x := scala.Map(t.posNodes, func(key interface{}, value interface{}) interface{} {
		return int(value.(KoreanToken).Pos)
	}).([]interface{})

	sum := 0
	for _, v := range x {
		sum += v.(int)
	}

	return sum
}

func (t *ParsedChunk) getUnknownCoverage() int {
	return scala.FoldLeft(t.posNodes, 0, func(folded interface{}, key interface{}, value interface{}) interface{} {
		if value.(KoreanToken).Unknown {
			return folded.(int) + utf8.RuneCountInString(value.(KoreanToken).Text)
		}
		return folded.(int)
	}).(int)
}

func (t *ParsedChunk) getFreqScore() float64 {
	return scala.FoldLeft(t.posNodes, 0.0, func(folded interface{}, key interface{}, value interface{}) interface{} {
		if value.(KoreanToken).Pos == util.Noun || value.(KoreanToken).Pos == util.ProperNoun {
			f, e := util.KoreanEntityFreq[value.(KoreanToken).Text]
			if e == false {
				f = 0
			}
			return folded.(float64) + (1.0 - f)
		}
		return folded.(float64) + 1.0
	}).(float64) / float64(len(t.posNodes))
}

func (t *ParsedChunk) add(that ParsedChunk) ParsedChunk {
	x := []KoreanToken{}
	x = append(x, t.posNodes...)
	x = append(x, that.posNodes...)
	return NewParsedChunk(x, t.words+that.words, t.profile)
}

func (t *ParsedChunk) countPos(pos util.KoreanPos) int {
	return scala.Count(t.posNodes, func(value interface{}) bool {
		return value.(KoreanToken).Pos == pos
	})
}
