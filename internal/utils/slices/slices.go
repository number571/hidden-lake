package slices

import "slices"

func DeleteDuplicates(pStrs []string) []string {
	result := make([]string, 0, len(pStrs))
	mapping := make(map[string]struct{}, len(pStrs))
	for _, s := range pStrs {
		if _, ok := mapping[s]; ok {
			continue
		}
		mapping[s] = struct{}{}
		result = append(result, s)
	}
	return result
}

func DeleteValues(s []string, v string) []string {
	return slices.DeleteFunc(s, func(p string) bool { return p == v })
}

func AppendIfNotExist(pSlice []string, pStr string) []string {
	if slices.Contains(pSlice, pStr) {
		return pSlice
	}
	return append(pSlice, pStr)
}
