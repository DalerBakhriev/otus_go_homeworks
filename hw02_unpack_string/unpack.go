package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return s, nil
	}

	if unicode.IsDigit(rune(s[0])) {
		return "", ErrInvalidString
	}

	stringBuilder := &strings.Builder{}
	currLetter := rune(s[0])

	for _, letter := range s[1:] {
		if !unicode.IsDigit(letter) {
			if !unicode.IsDigit(currLetter) {
				stringBuilder.WriteRune(currLetter)
			}
			currLetter = letter
			continue
		}

		if unicode.IsDigit(currLetter) {
			return "", ErrInvalidString
		}
		multiplier, _ := strconv.Atoi(string(letter))
		stringBuilder.WriteString(strings.Repeat(string(currLetter), multiplier))
		currLetter = letter
	}
	stringBuilder.WriteRune(currLetter)

	return stringBuilder.String(), nil
}
