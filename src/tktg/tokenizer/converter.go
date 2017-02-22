package tokenizer

import "tktg/util"

func ConvertSliceOfInterfaceToSliceChunkMatch(value interface{}) []ChunkMatch {
	ret := []ChunkMatch{}
	for _, v := range value.([]interface{}) {
		ret = append(ret, v.(ChunkMatch))
	}
	return ret
}

func ConvertSliceOfInterfaceToSliceOfKoreanToken(v []interface{}) []KoreanToken {
	x := []KoreanToken{}
	for _, v := range v {
		x = append(x, v.(KoreanToken))
	}
	return x
}

func ConvertSliceOfInterfaceToSliceOfCandidateParse(value []interface{}) []CandidateParse {
	ret := []CandidateParse{}
	for _, v := range value {
		ret = append(ret, v.(CandidateParse))
	}
	return ret
}

func ConvertSliceOfInterfaceToSliceOfPossibleTrie(value []interface{}) []PossibleTrie {
	ret := []PossibleTrie{}
	for _, v := range value {
		ret = append(ret, v.(PossibleTrie))
	}
	return ret
}

func ConvertSliceOfInterfaceToSliceOfKoreanPosTrie(value []interface{}) []util.KoreanPosTrie {
	ret := []util.KoreanPosTrie{}
	for _, v := range value {
		ret = append(ret, v.(util.KoreanPosTrie))
	}
	return ret
}
