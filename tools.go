package toolkit

import "crypto/rand"

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+-"

// Tools is the type to instantitate this module. Any varialee of this type will havee access 
// to all the methods with the receiver *Tools
type Tools struct{}

// RandomString returns a string of random character of lenght n,
// using randomStringSource as the source for the string
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x , y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	return string(s)
}