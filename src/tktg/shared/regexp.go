package shared

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

func ReplaceAll(re pcre.Regexp, text string, f func(m *pcre.Matcher) string) string {
	ret := ""

	start := 0
	for {
		x := []byte(text[start:])
		offset := re.FindIndex(x, 0)

		if offset == nil {
			break
		}

		pt := text[start:start+offset[0]]
		if len(pt) > 0 {
			ret += pt
		}

		ret += f(re.Matcher(x, 0))

		start += offset[1]
	}

	if start < len(text) {
		ret += string(text[start:])
	}

	return ret
}
