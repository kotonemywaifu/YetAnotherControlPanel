package util

func MatchPointerSlice(a []*string, b string) bool {
	for _, v := range a {
		if *v == b {
			return true
		}
	}

	return false
}
