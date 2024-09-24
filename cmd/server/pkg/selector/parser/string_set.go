package parser

import "sort"

type StringSet []string

func (ss StringSet) Contains(s string) bool {
	idx := sort.SearchStrings(ss, s)
	return idx < len(ss) && ss[idx] == s
}

func (ss StringSet) SliceCopy() []string {
	if ss == nil {
		return nil
	}
	cp := make([]string, len(ss), len(ss))
	copy(cp, ss)
	return cp
}

func ConvertToStringSetInPlace(s []string) StringSet {
	if len(s) <= 1 {
		return s
	}
	sort.Strings(s)
	out := s[0:1]
	for _, v := range s[1:] {
		if v == out[len(out)-1] {
			continue
		}
		out = append(out, v)
	}
	return out
}
