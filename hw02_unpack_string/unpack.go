package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	digitSymbol  = "digit"
	letterSymbol = "letter"
	slashSymbol  = "slash"
)

var ErrInvalidString = errors.New("invalid string")

func setState(symbol rune) string {
	var state string
	if string(symbol) == `\` {
		return slashSymbol
	}

	if unicode.IsDigit(symbol) {
		state = digitSymbol
	} else {
		state = letterSymbol
	}

	return state
}

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
		state := setState(currLetter)

		switch state {
		case digitSymbol:
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

		case letterSymbol:
			if string(prevLetter) == `\` {
				return "", ErrInvalidString
			}
			if !unicode.IsDigit(prevLetter) || prevLetterIsEscaped {
				stringBuilder.WriteRune(prevLetter)
			}

		case slashSymbol:
			if prevLetterIsEscaped {
				stringBuilder.WriteRune(prevLetter)
			}
			prevLetterIsEscaped = !prevLetterIsEscaped
		}

		prevLetter = currLetter
	}
	if !unicode.IsDigit(prevLetter) || prevLetterIsEscaped {
		stringBuilder.WriteRune(prevLetter)
	}

	return stringBuilder.String(), nil
}
