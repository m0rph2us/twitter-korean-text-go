package tokenizer

import (
	"regexp"
	"strings"
	"testing"
	"tktg/util"
)

func TestRegexp(t *testing.T) {
	matched := posPatterns[util.URL].MatcherString("", 0)

	if !(matched.Matches() == false) {
		t.Error("Expected result should match.")
	}
}

func TestRegExpKorean(t *testing.T) {
	matched := posPatterns[util.Korean].MatcherString("가나다라마바사", 0)

	if !(matched.Matches() == true) {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.Korean].MatcherString("12345", 0)

	if !(matched.Matches() == false) {
		t.Error("Expected result should match.")
	}
}

func TestRegExpAlpha(t *testing.T) {
	matched := posPatterns[util.Alpha].MatcherString("abcdef", 0)

	if !(matched.Matches() == true) {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.Alpha].MatcherString("12345", 0)

	if !(matched.Matches() == true) {
		t.Error("Expected result should match.")
	}
}

func TestRegExpNumber(t *testing.T) {
	matched := posPatterns[util.Number].MatcherString("3억3천만원", 0)

	if !(matched.Matches() == true) {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.Number].MatcherString("가나다라", 0)

	if !(matched.Matches() == false) {
		t.Error("Expected result should match.")
	}
}

func TestRegExpKoreanParticle(t *testing.T) {
	matched := posPatterns[util.KoreanParticle].MatcherString("ㅋㅋ", 0)

	if matched.Matches() != true {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.KoreanParticle].MatcherString("가나다라", 0)

	if matched.Matches() != false {
		t.Error("Expected result should match.")
	}
}

func TestRegExpPunctuation(t *testing.T) {
	matched := posPatterns[util.Punctuation].MatcherString("!@!~", 0)

	if matched.Matches() != true {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.KoreanParticle].MatcherString("가나다라", 0)

	if matched.Matches() != false {
		t.Error("Expected result should match.")
	}
}

func TestRegExpURL(t *testing.T) {
	matched := posPatterns[util.URL].MatcherString(
		"http://news.kukinews.com/article/view.asp?page=1&gCode=soc&arcid=0008599913&code=41121111", 0)

	if matched.Matches() != true {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.URL].MatcherString("가나다라", 0)

	if matched.Matches() != false {
		t.Error("Expected result should match.")
	}
}

func TestRegExpEmail(t *testing.T) {
	matched := posPatterns[util.Email].MatcherString("acid.acidrain@gmail.com", 0)

	if matched.Matches() != true {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.Email].MatcherString("가나다라", 0)

	if matched.Matches() != false {
		t.Error("Expected result should match.")
	}
}

func TestRegExpHashTag(t *testing.T) {
	matched := posPatterns[util.Hashtag].MatcherString("#korean_tokenizer_rocks", 0)

	if matched.Matches() != true {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.Hashtag].MatcherString("가나다라", 0)

	if matched.Matches() != false {
		t.Error("Expected result should match.")
	}
}

func TestRegExpScreenName(t *testing.T) {
	matched := posPatterns[util.ScreenName].MatcherString("@edeng", 0)

	if matched.Matches() != true {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.ScreenName].MatcherString("가나다라", 0)

	if matched.Matches() != false {
		t.Error("Expected result should match.")
	}
}

func TestRegExpCashTag(t *testing.T) {
	matched := posPatterns[util.CashTag].MatcherString("$edeng", 0)

	if matched.Matches() != true {
		t.Error("Expected result should match.")
	}

	matched = posPatterns[util.CashTag].MatcherString("가나다라", 0)

	if matched.Matches() != false {
		t.Error("Expected result should match.")
	}
}

func TestRegExp(t *testing.T) {

	x := regexp.MustCompile(`(1)2`)

	testStr := "121212"

	loc := x.FindAllStringIndex(testStr, -1)

	if !(loc[0][0] == 0 && loc[0][1] == 2 &&
		loc[1][0] == 2 && loc[1][1] == 4 &&
		loc[2][0] == 4 && loc[2][1] == 6) {
		t.Error("Expected result should match.")
	}
}

func TestSplitBySpaceKeepingSpace(t *testing.T) {
	x := splitBySpaceKeepingSpace("this is a test")

	if !(x[0] == "this" && x[1] == " " && x[2] == "is" && x[3] == " " &&
		x[4] == "a" && x[5] == " " && x[6] == "test") {
		t.Error("Expected result should match.")
	}
}

func TestGetChunks(t *testing.T) {
	c := getChunks("안녕? iphone6안녕? 세상아?", false)

	if !(strings.Join(c, "/") == "안녕/?/ /iphone/6/안녕/?/ /세상아/?") {
		t.Error("Expected result should match.")
	}

	c = getChunks("This is an 한국어가 섞인 English tweet.", false)

	if !(strings.Join(c, "/") == "This/ /is/ /an/ /한국어가/ /섞인/ /English/ /tweet/.") {
		t.Error("Expected result should match.")
	}

	c = getChunks("이 日本것은 日本語Eng", false)

	if !(strings.Join(c, "/") == "이/ /日本/것은/ /日本語/Eng") {
		t.Error("Expected result should match.")
	}

	c = getChunks("무효이며", false)

	if !(strings.Join(c, "/") == "무효이며") {
		t.Error("Expected result should match.")
	}

	c = getChunks("#해쉬태그 이라는 것 #hash @hello 123 이런이런 #여자최애캐_5명으로_취향을_드러내자", false)

	if !(strings.Join(c, "/") == "#해쉬태그/ /이라는/ /것/ /#hash/ /@hello/ /123/ /이런이런/ /#여자최애캐_5명으로_취향을_드러내자") {
		t.Error("Expected result should match.")
	}
}

func TestGetChunksWithNumbers(t *testing.T) {
	c := getChunks("300위안짜리 밥", false)

	if !(strings.Join(c, "/") == "300위안/짜리/ /밥") {
		t.Error("Expected result should match.")
	}

	c = getChunks("200달러와 300유로", false)

	if !(strings.Join(c, "/") == "200달러/와/ /300유로") {
		t.Error("Expected result should match.")
	}

	c = getChunks("$200이나 한다", false)

	if !(strings.Join(c, "/") == "$200/이나/ /한다") {
		t.Error("Expected result should match.")
	}

	c = getChunks("300옌이었다.", false)

	if !(strings.Join(c, "/") == "300옌/이었다/.") {
		t.Error("Expected result should match.")
	}

	c = getChunks("3,453,123,123원 3억3천만원", false)

	if !(strings.Join(c, "/") == "3,453,123,123원/ /3억/3천만원") {
		t.Error("Expected result should match.")
	}

	c = getChunks("6/4 지방 선거", false)

	if !(strings.Join(c, "/") == "6/4/ /지방/ /선거") {
		t.Error("Expected result should match.")
	}

	c = getChunks("6.4 지방 선거", false)

	if !(strings.Join(c, "/") == "6.4/ /지방/ /선거") {
		t.Error("Expected result should match.")
	}

	c = getChunks("6-4 지방 선거", false)

	if !(strings.Join(c, "/") == "6-4/ /지방/ /선거") {
		t.Error("Expected result should match.")
	}

	c = getChunks("6.25 전쟁", false)

	if !(strings.Join(c, "/") == "6.25/ /전쟁") {
		t.Error("Expected result should match.")
	}

	c = getChunks("1998년 5월 28일", false)

	if !(strings.Join(c, "/") == "1998년/ /5월/ /28일") {
		t.Error("Expected result should match.")
	}

	c = getChunks("62:45의 결과", false)

	if !(strings.Join(c, "/") == "62:45/의/ /결과") {
		t.Error("Expected result should match.")
	}

	c = getChunks("여러 칸  띄어쓰기,   하나의 Space묶음으로 처리됩니다.", false)

	if !(strings.Join(c, "/") == "여러/ /칸/  /띄어쓰기/,/   /하나의/ /Space/묶음으로/ /처리됩니다/.") {
		t.Error("Expected result should match.")
	}
}
