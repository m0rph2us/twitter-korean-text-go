package tokenizer

import (
	"fmt"
	"testing"
	"tktg/util"
	"tktg/data_structure"
)

func TestTokenize01(t *testing.T) {
	v := Tokenize("쵸귀여운개루루")

	d := []KoreanToken{
		{"쵸", util.VerbPrefix, 0, 1, false},
		{"귀여운", util.Adjective, 1, 3, false},
		{"개", util.Noun, 4, 1, false},
		{"루루", util.ProperNoun, 5, 2, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize02(t *testing.T) {
	v := Tokenize("개루루야")

	d := []KoreanToken{
		{"개", util.Noun, 0, 1, false},
		{"루루", util.ProperNoun, 1, 2, false},
		{"야", util.Josa, 3, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize03(t *testing.T) {
	v := Tokenize("쵸귀여운")

	d := []KoreanToken{
		{"쵸", util.VerbPrefix, 0, 1, false},
		{"귀여운", util.Adjective, 1, 3, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize04(t *testing.T) {
	v := Tokenize("이사람의")

	d := []KoreanToken{
		{"이", util.Determiner, 0, 1, false},
		{"사람", util.Noun, 1, 2, false},
		{"의", util.Josa, 3, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize05(t *testing.T) {
	v := Tokenize("엄청작아서귀엽다")

	d := []KoreanToken{
		{"엄청", util.Adverb, 0, 2, false},
		{"작아", util.Adjective, 2, 2, false},
		{"서", util.Eomi, 4, 1, false},
		{"귀엽", util.Adjective, 5, 2, false},
		{"다", util.Eomi, 7, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize06(t *testing.T) {
	v := Tokenize("안녕하셨어요")

	d := []KoreanToken{
		{"안녕하셨", util.Adjective, 0, 4, false},
		{"어요", util.Eomi, 4, 2, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize07(t *testing.T) {
	v := Tokenize("그리고")

	d := []KoreanToken{
		{"그리고", util.Conjunction, 0, 3, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize08(t *testing.T) {
	v := Tokenize("안녕ㅋㅋ")

	d := []KoreanToken{
		{"안녕", util.Noun, 0, 2, false},
		{"ㅋㅋ", util.KoreanParticle, 2, 2, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize09(t *testing.T) {
	v := Tokenize("라고만")

	d := []KoreanToken{
		{"라고만", util.Eomi, 0, 3, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize10(t *testing.T) {
	v := Tokenize("개컁컁아")

	d := []KoreanToken{
		{"개컁컁", util.ProperNoun, 0, 3, true},
		{"아", util.Josa, 3, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize11(t *testing.T) {
	v := Tokenize("안녕하세요쿛툐캬님")

	d := []KoreanToken{
		{"안녕하세", util.Adjective, 0, 4, false},
		{"요", util.Eomi, 4, 1, false},
		{"쿛툐캬", util.ProperNoun, 5, 3, true},
		{"님", util.Suffix, 8, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize12(t *testing.T) {
	v := Tokenize("이승기가")

	d := []KoreanToken{
		{"이승기", util.ProperNoun, 0, 3, false},
		{"가", util.Josa, 3, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize13(t *testing.T) {
	v := Tokenize("야이건뭐")

	d := []KoreanToken{
		{"야", util.Exclamation, 0, 1, false},
		{"이건", util.Noun, 1, 2, false},
		{"뭐", util.Noun, 3, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize14(t *testing.T) {
	v := Tokenize("아이럴수가")

	d := []KoreanToken{
		{"아", util.Exclamation, 0, 1, false},
		{"이럴", util.Adjective, 1, 2, false},
		{"수", util.PreEomi, 3, 1, false},
		{"가", util.Eomi, 4, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize15(t *testing.T) {
	v := Tokenize("보다가")

	d := []KoreanToken{
		{"보다", util.Verb, 0, 2, false},
		{"가", util.Eomi, 2, 1, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize16(t *testing.T) {
	v := Tokenize("하...")

	d := []KoreanToken{
		{"하", util.Exclamation, 0, 1, false},
		{"...", util.Punctuation, 1, 3, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize17(t *testing.T) {
	v := Tokenize("시전하는")

	d := []KoreanToken{
		{"시전", util.Noun, 0, 2, false},
		{"하는", util.Verb, 2, 2, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize18(t *testing.T) {
	v := Tokenize("해쵸쵸쵸쵸쵸쵸쵸쵸춏")

	d := []KoreanToken{
		{"해쵸쵸쵸쵸쵸쵸쵸쵸춏", util.Noun, 0, 10, true},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize19(t *testing.T) {
	v := Tokenize("뇬뇨뇬뇨뇬뇨뇬뇨츄쵸")

	d := []KoreanToken{
		{"뇬뇨뇬뇨뇬뇨뇬뇨", util.ProperNoun, 0, 8, true},
		{"츄쵸", util.ProperNoun, 8, 2, true},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize20(t *testing.T) {
	util.AddWordsToDictionary(util.Noun, []string{"뇬뇨", "츄쵸"})

	if !util.KoreanDictionary.Get(util.Noun).(data_structure.Set).Contains("뇬뇨") {
		t.Error("Expected result should match.")
	}

	if !util.KoreanDictionary.Get(util.Noun).(data_structure.Set).Contains("츄쵸") {
		t.Error("Expected result should match.")
	}

	v := Tokenize("뇬뇨뇬뇨뇬뇨뇬뇨츄쵸")

	d := []KoreanToken{
		{"뇬뇨", util.Noun, 0, 2, false},
		{"뇬뇨", util.Noun, 2, 2, false},
		{"뇬뇨", util.Noun, 4, 2, false},
		{"뇬뇨", util.Noun, 6, 2, false},
		{"츄쵸", util.Noun, 8, 2, false},
	}

	if !(fmt.Sprintf("%v", v) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestTokenize21(t *testing.T) {
	c := Tokenize("블랙프라이데이: 이날 미국의 수백만 소비자들은 크리스마스 선물을 할인된 가격에 사는 것을 주 목적으로 블랙프라이데이 쇼핑을 한다.")

	d := []KoreanToken{
		{"블랙프라이데이", util.Noun, 0, 7, false},
		{":", util.Punctuation, 7, 1, false},
		{" ", util.Space, 8, 1, false},
		{"이", util.Determiner, 9, 1, false},
		{"날", util.Noun, 10, 1, false},
		{" ", util.Space, 11, 1, false},
		{"미국", util.ProperNoun, 12, 2, false},
		{"의", util.Josa, 14, 1, false},
		{" ", util.Space, 15, 1, false},
		{"수백만", util.Noun, 16, 3, false},
		{" ", util.Space, 19, 1, false},
		{"소비자", util.Noun, 20, 3, false},
		{"들", util.Suffix, 23, 1, false},
		{"은", util.Josa, 24, 1, false},
		{" ", util.Space, 25, 1, false},
		{"크리스마스", util.Noun, 26, 5, false},
		{" ", util.Space, 31, 1, false},
		{"선물", util.Noun, 32, 2, false},
		{"을", util.Josa, 34, 1, false},
		{" ", util.Space, 35, 1, false},
		{"할인", util.Noun, 36, 2, false},
		{"된", util.Verb, 38, 1, false},
		{" ", util.Space, 39, 1, false},
		{"가격", util.Noun, 40, 2, false},
		{"에", util.Josa, 42, 1, false},
		{" ", util.Space, 43, 1, false},
		{"사는", util.Verb, 44, 2, false},
		{" ", util.Space, 46, 1, false},
		{"것", util.Noun, 47, 1, false},
		{"을", util.Josa, 48, 1, false},
		{" ", util.Space, 49, 1, false},
		{"주", util.Noun, 50, 1, false},
		{" ", util.Space, 51, 1, false},
		{"목적", util.Noun, 52, 2, false},
		{"으로", util.Josa, 54, 2, false},
		{" ", util.Space, 56, 1, false},
		{"블랙프라이데이", util.Noun, 57, 7, false},
		{" ", util.Space, 64, 1, false},
		{"쇼핑", util.Noun, 65, 2, false},
		{"을", util.Josa, 67, 1, false},
		{" ", util.Space, 68, 1, false},
		{"한", util.Verb, 69, 1, false},
		{"다", util.Eomi, 70, 1, false},
		{".", util.Punctuation, 71, 1, false},
	}

	if !(fmt.Sprintf("%v", c) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}

func TestCollapseNouns(t *testing.T) {
	c := collapseNouns([]KoreanToken{
		{"마", util.Noun, 0, 1, false},
		{"코", util.Noun, 1, 1, false},
		{"토", util.Noun, 2, 1, false},
	})

	d := []KoreanToken{
		{"마코토", util.Noun, 0, 3, true},
	}

	if !(fmt.Sprintf("%v", c) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}

	c = collapseNouns([]KoreanToken{
		{"마", util.Noun, 0, 1, false},
		{"코", util.Noun, 1, 1, false},
		{"토", util.Noun, 2, 1, false},
		{"를", util.Josa, 3, 1, false},
	})

	d = []KoreanToken{
		{"마코토", util.Noun, 0, 3, true},
		{"를", util.Josa, 3, 1, false},
	}

	if !(fmt.Sprintf("%v", c) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}

	c = collapseNouns([]KoreanToken{
		{"개", util.NounPrefix, 0, 1, false},
		{"마", util.Noun, 1, 1, false},
		{"코", util.Noun, 2, 1, false},
		{"토", util.Noun, 3, 1, false},
	})

	d = []KoreanToken{
		{"개", util.NounPrefix, 0, 1, false},
		{"마코토", util.Noun, 1, 3, true},
	}

	if !(fmt.Sprintf("%v", c) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}

	c = collapseNouns([]KoreanToken{
		{"마", util.Noun, 0, 1, false},
		{"코", util.Noun, 1, 1, false},
		{"토", util.Noun, 2, 1, false},
		{"사람", util.Noun, 3, 2, false},
	})

	d = []KoreanToken{
		{"마코토", util.Noun, 0, 3, true},
		{"사람", util.Noun, 3, 2, false},
	}

	if !(fmt.Sprintf("%v", c) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}

	c = collapseNouns([]KoreanToken{
		{"마", util.Noun, 0, 1, false},
		{"코", util.Noun, 1, 1, false},
		{"사람", util.Noun, 2, 2, false},
		{"토", util.Noun, 4, 1, false},
	})

	d = []KoreanToken{
		{"마코", util.Noun, 0, 2, true},
		{"사람", util.Noun, 2, 2, false},
		{"토", util.Noun, 4, 1, false},
	}

	if !(fmt.Sprintf("%v", c) == fmt.Sprintf("%v", d)) {
		t.Error("Expected result should match.")
	}
}
