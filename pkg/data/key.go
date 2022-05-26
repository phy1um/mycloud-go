package data

import "math/rand"

var chars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@#&*+"

func RandomKey() string {
	out := ""
	for i := 0; i < 6; i++ {
		n := rand.Int() % len(chars)
		out += string(chars[n])
	}
	return out
}
