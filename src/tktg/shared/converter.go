package shared

func ConvertSliceOfInterfaceToSliceOfString(si []interface{}) []string {
	ret := []string{}
	for _, v := range si {
		ret = append(ret, v.(string))
	}
	return ret
}
