package normalizer

import (
	"testing"
)

func TestNormalize01(t *testing.T) {
	c := Normalize("무의식중에 손들어버려섴ㅋㅋㅋㅋ")

	if !(c == "무의식중에 손들어버려서ㅋㅋ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize02(t *testing.T) {
	c := Normalize("안됔ㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋ내 심장을 가격했엌ㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋㅋ")

	if !(c == "안돼ㅋㅋ내 심장을 가격했어ㅋㅋ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize03(t *testing.T) {
	c := Normalize("기억도 나지아낳ㅎㅎㅎ")

	if !(c == "기억도 나지아나ㅎㅎ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize04(t *testing.T) {
	c := Normalize("근데비싸서못머구뮤ㅠㅠ")

	if !(c == "근데비싸서못먹음ㅠㅠ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize05(t *testing.T) {
	c := Normalize("미친 존잘니뮤ㅠㅠㅠㅠ")

	if !(c == "미친 존잘님ㅠㅠ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize06(t *testing.T) {
	c := Normalize("만나무ㅜㅜㅠ")

	if !(c == "만남ㅜㅜ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize07(t *testing.T) {
	c := Normalize("가루ㅜㅜㅜㅜ")

	if !(c == "가루ㅜㅜ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize08(t *testing.T) {
	c := Normalize("최지우ㅜㅜㅜㅜ")

	if !(c == "최지우ㅜㅜ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize09(t *testing.T) {
	c := Normalize("유성우ㅠㅠㅠ")

	if !(c == "유성우ㅠㅠ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize10(t *testing.T) {
	c := Normalize("ㅎㅎㅎㅋㅋ트위터ㅋㅎㅋ월드컵ㅠㅜㅠㅜㅠ")

	if !(c == "ㅎㅎㅋㅋ트위터ㅋㅎㅋ월드컵ㅠㅜ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize11(t *testing.T) {
	c := Normalize("예뿌ㅠㅠ")

	if !(c == "예뻐ㅠㅠ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize12(t *testing.T) {
	c := Normalize("고수야고수ㅠㅠ")

	if !(c == "고수야고수ㅠㅠ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize13(t *testing.T) {
	c := Normalize("땡큐우우우우우우")

	if !(c == "땡큐우우") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize14(t *testing.T) {
	c := Normalize("구오오오오오오오오옹오오오")

	if !(c == "구오오옹오오") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize15(t *testing.T) {
	c := Normalize("훌쩍훌쩍훌쩍훌쩍훌쩍훌쩍훌쩍훌쩍훌쩍훌쩍훌쩍훌쩍훌쩍")

	if !(c == "훌쩍훌쩍") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize16(t *testing.T) {
	c := Normalize("ㅋㅎㅋㅎㅋㅎㅋㅎㅋㅎㅋㅎ")

	if !(c == "ㅋㅎㅋㅎ") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize17(t *testing.T) {
	c := Normalize("http://11111.cccccom soooooooo !!!!!!!!!!!!!!!")

	if !(c == "http://11111.cccccom soooooooo !!!!!!!!!!!!!!!") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize18(t *testing.T) {
	c := Normalize("가쟝 용기있는 사람이 머굼 되는거즤")

	if !(c == "가장 용기있는 사람이 먹음 되는거지") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize19(t *testing.T) {
	c := Normalize("오노딘가")

	if !(c == "오노디인가") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize20(t *testing.T) {
	c := Normalize("관곈지")

	if !(c == "관계인지") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize21(t *testing.T) {
	c := Normalize("생각하는건데")

	if !(c == "생각하는건데") {
		t.Error("Expected result should match.")
	}
}

func TestNormalize22(t *testing.T) {
	c := Normalize("생각(하는ㅋㅋㅋ)건데")

	if !(c == "생각(하는ㅋㅋ)거인데") {
		t.Error("Expected result should match.")
	}
}
