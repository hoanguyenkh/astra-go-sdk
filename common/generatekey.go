package common

import (
	"math/rand"
	"time"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateSecretKeyRandomString(n int) (string, error) {
	key, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}

	return encodeBase64(key), nil
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	digits := "0123456789"
	specials := "!@#$%&*+_-="
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + digits + specials

	buf := make([]byte, length)

	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = letters[rand.Intn(len(letters))]
	}

	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	str := string(buf)

	return str
}
