package main

import (
	"math/rand"
	"time"
)

// create random string of length=n
func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano())) // start seed
	str := make([]rune, length)
	for i := range str {
		str[i] = letters[rnd.Intn(len(letters))]
	}
	return string(str)
}

// split string into chunks of character length len
func chunkify(text string, len int) []string {
	var s string
	var set []string
	a := []rune(text)
	for i, r := range a {
		s = s + string(r)

		if (i > 0) && (i%len == 0) {
			set = append(set, s)
			s = ""
		}
	}

	set = append(set, s)
	return set
}
