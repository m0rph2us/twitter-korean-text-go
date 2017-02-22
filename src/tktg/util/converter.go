package util

func ConvertSliceOfInterfaceToSliceOfHangulRune(value []interface{}) []HangulRune {
	x := []HangulRune{}
	for _, v := range value {
		x = append(x, v.(HangulRune))
	}
	return x
}

func ConvertSliceOfInterfaceToSliceOfRune(value []interface{}) []rune {
	x := []rune{}
	for _, v := range value {
		x = append(x, v.(rune))
	}
	return x
}
