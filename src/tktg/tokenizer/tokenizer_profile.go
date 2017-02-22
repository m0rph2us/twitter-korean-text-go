package tokenizer

import (
	"tktg/util"
	"tktg/data_structure"
)

type TokenizerProfile struct {
	tokenCount          float64
	unknown             float64
	wordCount           float64
	freq                float64
	unknownCoverage     float64
	exactMatch          float64
	allNoun             float64
	unknownPosCount     float64
	determinerPosCount  float64
	exclamationPosCount float64
	initialPostPosition float64
	haVerb              float64
	preferredPattern    float64
	preferredPatterns   [][]interface{}
	spaceGuide          data_structure.Set /* int */
	spaceGuidePenalty   float64
}

var defaultProfile TokenizerProfile = TokenizerProfile{
	tokenCount:          0.18,
	unknown:             0.3,
	wordCount:           0.3,
	freq:                0.2,
	unknownCoverage:     0.5,
	exactMatch:          0.5,
	allNoun:             0.1,
	unknownPosCount:     10.0,
	determinerPosCount:  -0.01,
	exclamationPosCount: 0.01,
	initialPostPosition: 0.2,
	haVerb:              0.3,
	preferredPattern:    0.6,
	preferredPatterns: [][]interface{}{
		{util.Noun, util.Josa},
		{util.ProperNoun, util.Josa},
	},
	spaceGuide:        data_structure.NewLinkedSet(),
	spaceGuidePenalty: 3.0,
}
