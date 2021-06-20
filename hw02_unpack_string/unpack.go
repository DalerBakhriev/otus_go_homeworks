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
	escaping := true // флаг экранирования

	for _, letter := range s[1:] {
		if unicode.IsDigit(letter) {
			if unicode.IsDigit(currLetter) && !escaping {
				return "", ErrInvalidString
			}
			if (string(currLetter) != `\`) || (string(currLetter) == `\` && escaping) {
				multiplier, _ := strconv.Atoi(string(letter))
				stringBuilder.WriteString(strings.Repeat(string(currLetter), multiplier))
				escaping = false
			} else {
				escaping = true
			}
			currLetter = letter
			continue
		}

		if string(letter) == `\` {
			if escaping {
				stringBuilder.WriteRune(currLetter)
			}
			escaping = !escaping
			currLetter = letter
			continue
		}
		if string(currLetter) == `\` {
			return "", ErrInvalidString
		}
		if !unicode.IsDigit(currLetter) || escaping {
			stringBuilder.WriteRune(currLetter)
		}
		currLetter = letter
	}
	if !unicode.IsDigit(currLetter) || escaping {
		stringBuilder.WriteRune(currLetter)
	}

	return stringBuilder.String(), nil
}
