package util

import "math/rand"

var charset = [...]string{
	"2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H",
	"J", "K", "L", "M", "N", "P", "Q", "R",
	"S", "T", "U", "V", "W", "X", "Y", "Z",
}

func GenerateCaptcha(length int) string {
	var i = 1
	ret := ""
	for i <= length {
		i++
		ret += charset[rand.Intn(len(charset))]
	}
	return ret
}
