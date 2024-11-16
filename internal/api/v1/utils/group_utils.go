package utils

import (
	"math/rand"
)

func GenerateGroupInviteCode() string {
	var letterRunes = []rune("0123456789")

	b := make([]rune, 4)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
