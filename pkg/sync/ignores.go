package sync

func determineLeftOuterUnion(left, right []string) []string {
	var result []string
	exists := func(needle string, haystack []string) bool {
		for _, v := range haystack {
			if needle == v {
				return true
			}
		}
		return false
	}
	for _, v := range left {
		if !exists(v, right) {
			result = append(result, v)
		}
	}
	return result
}
