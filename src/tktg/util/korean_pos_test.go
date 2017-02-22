package util

import (
	"fmt"
	"testing"
)

func TestBuildTrie00(t *testing.T) {
	b := buildTrie("p1N+", Noun)

	c := []KoreanPosTrie{
		{
			NounPrefix,
			[]KoreanPosTrie{
				{
					Noun,
					[]KoreanPosTrie{
						SelfNode,
					},
					Noun,
				},
			},
			None,
		},
	}

	if !(fmt.Sprintf("%v", b) == fmt.Sprintf("%v", c)) {
		t.Error("Expected result should match.")
	}
}

func TestBuildTrie01(t *testing.T) {
	b := buildTrie("p1N*", Noun)

	c := []KoreanPosTrie{
		{
			NounPrefix,
			[]KoreanPosTrie{
				{
					Noun,
					[]KoreanPosTrie{
						SelfNode,
					},
					Noun,
				},
			},
			Noun,
		},
	}

	if !(fmt.Sprintf("%v", b) == fmt.Sprintf("%v", c)) {
		t.Error("Expected result should match.")
	}

	b = buildTrie("N+s0", Noun)

	c = []KoreanPosTrie{
		{
			Noun,
			[]KoreanPosTrie{
				SelfNode,
				{
					Suffix,
					[]KoreanPosTrie{},
					Noun,
				},
			},
			Noun,
		},
	}

	if !(fmt.Sprintf("%v", b) == fmt.Sprintf("%v", c)) {
		t.Error("Expected result should match.")
	}
}

func TestBuildTrie02(t *testing.T) {
	b := buildTrie("A+V+A0", Verb)

	c := []KoreanPosTrie{
		{
			Adverb,
			[]KoreanPosTrie{
				SelfNode,
				{
					Verb,
					[]KoreanPosTrie{
						SelfNode,
						KoreanPosTrie{
							Adverb,
							[]KoreanPosTrie{},
							Verb,
						},
					},
					Verb,
				},
			},
			None,
		},
	}

	if !(fmt.Sprintf("%v", b) == fmt.Sprintf("%v", c)) {
		t.Error("Expected result should match.")
	}
}
