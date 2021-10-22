package hw09structvalidator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		// meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in              interface{}
		expectedErr     error
		expectedErrsNum int
	}{
		{
			in:              Response{Code: 300},
			expectedErr:     errNotInSet,
			expectedErrsNum: 1,
		},
		{
			in:              App{Version: "long_version_tag"},
			expectedErr:     errLength,
			expectedErrsNum: 1,
		},
		{
			in: User{
				ID:     strings.Repeat("q", 36),
				Name:   "tstName",
				Age:    51,
				Role:   "admin",
				Phones: []string{strings.Repeat("1", 11)},
				Email:  "smth@mail.ru",
			},
			expectedErr:     errMaxThreshold,
			expectedErrsNum: 1,
		},
		{
			in: User{
				ID:     strings.Repeat("q", 36),
				Name:   "tstName",
				Age:    17,
				Role:   "admin",
				Phones: []string{strings.Repeat("1", 11)},
				Email:  "smth@mail.ru",
			},
			expectedErr:     errMinThreshold,
			expectedErrsNum: 1,
		},
		{
			in: User{
				ID:     strings.Repeat("q", 36),
				Name:   "tstName",
				Age:    18,
				Role:   "smth",
				Phones: []string{strings.Repeat("1", 11)},
				Email:  "smth@mail.ru",
			},
			expectedErr:     errNotInSet,
			expectedErrsNum: 1,
		},
		{
			in: User{
				ID:     strings.Repeat("q", 36),
				Name:   "tstName",
				Age:    18,
				Role:   "admin",
				Phones: []string{strings.Repeat("1", 11), strings.Repeat("1", 9)},
				Email:  "smth@mail.ru",
			},
			expectedErr:     errLength,
			expectedErrsNum: 1,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			// t.Parallel()
			err := Validate(tt.in)
			assert.Error(t, err)
			var validErrs ValidationErrors
			assert.ErrorAs(t, err, &validErrs)
			assert.Len(t, validErrs, tt.expectedErrsNum)
			assert.Contains(t, validErrs[0].Error(), tt.expectedErr.Error())
		})
	}
}
