package ir

import ()

func setDefaultType(h map[string]Type, k string, t Type) (r Type) {
	var set bool
	if r, set = h[k]; !set {
		h[k] = t
		r = t
		set = true
	}
	return
}

func stackMapPush(s map[string][]Data, k string, v Data) {
	s[k] = append(s[k], v)
}

func stackMapPop(s map[string][]Data, k string) {
	n := len(s[k]) - 1
	s[k] = s[k][:n]
}

func stackMapPeek(s map[string][]Data, k string) (Data, bool) {
	l := len(s[k])
	if l == 0 {
		return nil, false
	}
	return s[k][l-1], true
}
