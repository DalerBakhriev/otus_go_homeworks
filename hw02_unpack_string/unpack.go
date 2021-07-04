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

	if unicode.IsDigit([]rune(s)[0]) {
		return "", ErrInvalidString
	}

	stringBuilder := &strings.Builder{}
	prevLetter := []rune(s)[0]
	prevLetterIsEscaped := true // флаг экранирования

	for _, currLetter := range s[1:] {
		switch {
		case unicode.IsDigit(currLetter):
			if unicode.IsDigit(prevLetter) && !prevLetterIsEscaped {
				return "", ErrInvalidString
			}
			if string(prevLetter) != `\` || prevLetterIsEscaped {
				multiplier, _ := strconv.Atoi(string(currLetter))
				stringBuilder.WriteString(strings.Repeat(string(prevLetter), multiplier))
				prevLetterIsEscaped = false
			} else {
				prevLetterIsEscaped = true
			}

		case string(currLetter) == `\`:
			if prevLetterIsEscaped {
				stringBuilder.WriteRune(prevLetter)
			}
			prevLetterIsEscaped = !prevLetterIsEscaped

		default:
			if string(prevLetter) == `\` {
				return "", ErrInvalidString
			}
			if !unicode.IsDigit(prevLetter) || prevLetterIsEscaped {
				stringBuilder.WriteRune(prevLetter)
			}
		}

		prevLetter = currLetter
	}
	if !unicode.IsDigit(prevLetter) || prevLetterIsEscaped {
		stringBuilder.WriteRune(prevLetter)
	}

	return stringBuilder.String(), nil
}
