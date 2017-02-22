package tokenizer

import (
	"testing"
)

func TestDetokenizer(t *testing.T) {
	s := Detokenize([]string{"연세", "대학교", "보건", "대학원", "에", "오신", "것", "을", "환영", "합니다", "!"})

	if !(s == "연세대학교 보건 대학원에 오신것을 환영합니다!") {
		t.Error("Expected result should match.")
	}

	s = Detokenize([]string{"와", "!!!", "iPhone", "6+", "가", ",", "드디어", "나왔다", "!"})

	if !(s == "와!!! iPhone 6+ 가, 드디어 나왔다!") {
		t.Error("Expected result should match.")
	}

	s = Detokenize([]string{"뭐", "완벽", "하진", "않", "지만", "그럭저럭", "쓸", "만", "하군", "..."})

	if !(s == "뭐 완벽하진 않지만 그럭저럭 쓸 만하군...") {
		t.Error("Expected result should match.")
	}
}

func TestDetokenizer01(t *testing.T) {
	s := Detokenize([]string{""})

	if !(s == "") {
		t.Error("Expected result should match.")
	}

	s = Detokenize([]string{})

	if !(s == "") {
		t.Error("Expected result should match.")
	}

	s = Detokenize([]string{"완벽"})

	if !(s == "완벽") {
		t.Error("Expected result should match.")
	}

	s = Detokenize([]string{"이"})

	if !(s == "이") {
		t.Error("Expected result should match.")
	}

	s = Detokenize([]string{"이", "제품을", "사용하겠습니다"})

	if !(s == "이 제품을 사용하겠습니다") {
		t.Error("Expected result should match.")
	}
}
