package id

import "strings"

const (
	_              = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base33Alphabet = "0123456789ABCDEFGHJKLMNPRSTUVWXYZ" // without I, O, Q
)

// appendBase33 convert n to base33 string
func AppendBase33(w *strings.Builder, n int64, maxLen int) {
	s := getBase33String(n, maxLen)
	w.WriteString(reverseString(s))
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// getBase33String convert n to base33 string, up to maxLen
func getBase33String(n int64, maxLen int) string {
	var sb strings.Builder
	for i := 0; i < maxLen; i++ {
		sb.WriteByte(base33Alphabet[n%33])
		n /= 33
	}
	return sb.String()
}

// TranslateBase33 convert base33 string to int64
func TranslateBase33(s string) int64 {
	var n int64
	for i := len(s) - 1; i >= 0; i-- {
		n = n*33 + int64(strings.IndexByte(base33Alphabet, s[i]))
	}
	return n
}
