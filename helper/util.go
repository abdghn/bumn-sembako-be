package helper

import (
	"crypto/rand"
	"fmt"
	"strings"
)

const random string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"

func Randstring(n int) (string, error) {
	s, r := make([]rune, n), []rune(random)
	for i := range s {
		p, err := rand.Prime(rand.Reader, len(r))
		if err != nil {
			return "", fmt.Errorf("random string n %d: %w", n, err)
		}
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s), nil
}

func ContainString(alpha, value string) bool {
	for _, char := range value {
		if !strings.Contains(alpha, strings.ToLower(string(char))) {
			return true
		}
	}

	return false
}
