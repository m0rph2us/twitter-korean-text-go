package phrase_extractor

func ConvertSliceOfInterfaceToSliceOfKoreanPhrase(si []interface{}) []KoreanPhrase {
	ret := []KoreanPhrase{}
	for _, v := range si {
		ret = append(ret, v.(KoreanPhrase))
	}
	return ret
}

func ConvertSliceOfInterfaceToSliceOfKoreanPhraseChunk(si []interface{}) []KoreanPhraseChunk {
	ret := []KoreanPhraseChunk{}
	for _, v := range si {
		ret = append(ret, v.(KoreanPhraseChunk))
	}
	return ret
}
