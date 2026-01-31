package aliasgen

import (
	"math/rand"
	"time"
)

// GenRandomString generates random string with given size
func GenRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	alph := []rune("QWERTYUIOPASDFGHJKLZXCVBNM" + "qwertyuiopasdfghjklzxcvbnm" + "0123456789" + "_-")

	b := make([]rune, size)
	for i := range b {
		b[i] = alph[rnd.Intn(len(alph))]
	}

	return string(b)
}
