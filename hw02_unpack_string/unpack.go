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
	prevLetter := []rune(s)[0]
	escaping := true // флаг экранирования

	for _, letter := range s[1:] {
		if unicode.IsDigit(letter) {
			if unicode.IsDigit(prevLetter) && !escaping {
				return "", ErrInvalidString
			}
			if (string(prevLetter) != `\`) || (string(prevLetter) == `\` && escaping) {
				multiplier, _ := strconv.Atoi(string(letter))
				stringBuilder.WriteString(strings.Repeat(string(prevLetter), multiplier))
				escaping = false
			} else {
				escaping = true
			}
			prevLetter = letter
			continue
		}

		if string(letter) == `\` {
			if escaping {
				stringBuilder.WriteRune(prevLetter)
			}
			escaping = !escaping
			prevLetter = letter
			continue
		}
		if string(prevLetter) == `\` {
			return "", ErrInvalidString
		}
		if !unicode.IsDigit(prevLetter) || escaping {
			stringBuilder.WriteRune(prevLetter)
		}
		prevLetter = letter
	}
	if !unicode.IsDigit(prevLetter) || escaping {
		stringBuilder.WriteRune(prevLetter)
	}

	return stringBuilder.String(), nil
}
