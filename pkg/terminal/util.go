package terminal

import "unicode"

func TrimSpaceLeft(r []rune) []rune {
	i := len(r)
	for j, v := range r {
		if !unicode.IsSpace(v) {
			i = j
			break
		}
	}
	return r[i:]
}

func HasPrefix(pref, val []rune) bool {
	if len(pref) > len(val) {
		return false
	}
	return RuneEquals(pref, val[:len(pref)])
}

func RuneEquals(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
