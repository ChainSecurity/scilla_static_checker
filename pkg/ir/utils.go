package ir

import (
	"fmt"
)

func setDefaultType(h map[string]Type, k string, t Type) (r Type) {
	var set bool
	if r, set = h[k]; !set {
		h[k] = t
		r = t
		set = true
	}
	return
}

func stackMapCopy(s map[string][]Data) map[string][]Data {
	sCopy := map[string][]Data{}
	for k, vals := range s {
		valsCopy := make([]Data, len(vals))
		copy(valsCopy, vals)
		sCopy[k] = valsCopy
	}
	return sCopy
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

func TypeOf(d Data) Type {
	switch x := d.(type) {
	case *Nat:
		return x.NatType
	case *DataVar:
		return x.DataType
	case *Load:
	case *AppDD:
		return TypeOf(x.To)
	case *AppTD:
		return TypeOf(x.To)
	default:
		fmt.Printf("TypeOf not implemented %T\n", d)
	}

	return nil
}

func KindOf(t Type) Kind {
	return nil
}
