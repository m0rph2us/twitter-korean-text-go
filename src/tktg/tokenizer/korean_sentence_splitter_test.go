package tokenizer

import (
	"testing"
	"fmt"
)

func TestSplit00(t *testing.T) {
	d := Split("안녕? iphone6안녕? 세상아?")

	if !(fmt.Sprintf("%v", d) == "[{안녕? 0 3} {iphone6안녕? 4 14} {세상아? 15 19}]") {
		t.Error("Expected result should match.")
	}
}

func TestSplit01(t *testing.T) {
	d := Split("그런데, 누가 그러는데, 루루가 있대. 그렇대? 그렇지! 아리고 이럴수가!!!!! 그래...")

	if !(fmt.Sprintf("%v", d) ==
		"[{그런데, 누가 그러는데, 루루가 있대. 0 21} {그렇대? 22 26} {그렇지! 27 31} {아리고 이럴수가!!!!! 32 45} {그래... 46 51}]") {
		t.Error("Expected result should match.")
	}
}

func TestSplit02(t *testing.T) {
	d := Split("이게 말이 돼?! 으하하하 ㅋㅋㅋㅋㅋㅋㅋ…    ")

	if !(fmt.Sprintf("%v", d) == "[{이게 말이 돼?! 0 9} {으하하하 ㅋㅋㅋㅋㅋㅋㅋ… 10 23}]") {
		t.Error("Expected result should match.")
	}
}
