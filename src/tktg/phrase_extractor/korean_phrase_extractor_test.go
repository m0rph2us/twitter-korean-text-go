package phrase_extractor

import (
	"fmt"
	"testing"
	"tktg/tokenizer"
)

func TestExtractor01(t *testing.T) {
	d := ExtractPhrases(
		tokenizer.Tokenize("블랙프라이데이: 이날 미국의 수백만 소비자들은 크리스마스 선물을 할인된 가격에 사는 것을 주 목적으로 블랙프라이데이 쇼핑을 한다."),
		false,
		true)

	s := ""
	for _, v := range d {
		s += fmt.Sprintf("%v(%v:%v,%v),", v.text(), v.Pos, v.offset(), v.length())
	}

	if !(s == "블랙프라이데이(2:0,7),이날(2:9,2),이날 미국(2:9,5),이날 미국의 수백만(2:9,10),미국의 수백만(2:12,7),"+
		"수백만(2:16,3),이날 미국의 수백만 소비자들(2:9,15),미국의 수백만 소비자들(2:12,12),수백만 소비자들(2:16,8),"+
		"크리스마스(2:26,5),크리스마스 선물(2:26,8),할인(2:36,2),할인된 가격(2:36,6),가격(2:40,2),주 목적(2:50,4),"+
		"블랙프라이데이 쇼핑(2:57,10),미국(2:12,2),소비자들(2:20,4),선물(2:32,2),목적(2:52,2),쇼핑(2:65,2),") {
		t.Error("Expected result should match.")
	}
}

func TestExtractor02(t *testing.T) {
	d := ExtractPhrases(
		tokenizer.Tokenize("성탄절 쇼핑 성탄절 쇼핑 성탄절 쇼핑 성탄절 쇼핑"),
		false,
		true)

	s := ""
	for _, v := range d {
		s += fmt.Sprintf("%v(%v:%v,%v),", v.text(), v.Pos, v.offset(), v.length())
	}

	if !(s == "성탄절(2:0,3),성탄절 쇼핑(2:0,6),성탄절 쇼핑 성탄절(2:0,10),"+
		"쇼핑 성탄절(2:4,6),성탄절 쇼핑 성탄절 쇼핑(2:0,13),쇼핑 성탄절 쇼핑(2:4,9),"+
		"성탄절 쇼핑 성탄절 쇼핑 성탄절(2:0,17),쇼핑 성탄절 쇼핑 성탄절(2:4,13),성탄절 쇼핑 성탄절 "+
		"쇼핑 성탄절 쇼핑(2:0,20),쇼핑 성탄절 쇼핑 성탄절 쇼핑(2:4,16),성탄절 쇼핑 성탄절 쇼핑 "+
		"성탄절 쇼핑 성탄절(2:0,24),쇼핑 성탄절 쇼핑 성탄절 쇼핑 성탄절(2:4,20),성탄절 쇼핑 성탄절 "+
		"쇼핑 성탄절 쇼핑 성탄절 쇼핑(2:0,27),쇼핑 성탄절 쇼핑 성탄절 쇼핑 성탄절 쇼핑(2:4,23),"+
		"쇼핑(2:4,2),") {
		t.Error("Expected result should match.")
	}
}

func TestExtractor03(t *testing.T) {
	toks := tokenizer.Tokenize("떡볶이 3,444,231원 + 400원.")

	d := ExtractPhrases(
		toks,
		false,
		true)

	s := ""
	for _, v := range d {
		s += fmt.Sprintf("%v(%v:%v,%v),", v.text(), v.Pos, v.offset(), v.length())
	}

	if !(s == "떡볶이(2:0,3),떡볶이 3,444,231원(2:0,14),400원(2:17,4),3,444,231원(2:4,10),") {
		fmt.Printf("toks: %#v\n", toks)
		fmt.Printf("d: %#v\n", d)
		fmt.Printf("s: %v\n", s)
		t.Error("Expected result should match.")
	}
}
