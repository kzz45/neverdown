package syncer

import "strings"

func refactorDomains(in []string, rules []string) []string {
	res := make([]string, 0)
	for _, v := range in {
		res = append(res, refactorDomain(v, rules))
	}
	return res
}

func refactorDomain(in string, rules []string) string {
	if len(rules) < 2 {
		return in
	}
	if len(rules) >= 2 {
		in = strings.Replace(in, rules[0], rules[1], -1)
	}
	if len(rules) < 4 {
		return in
	}
	return strings.Replace(in, rules[2], rules[3], -1)
}
