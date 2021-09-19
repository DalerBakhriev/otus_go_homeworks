package main

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		name          string
		offset        int64
		limit         int64
		expectedError error
	}{
		{
			name:   "out_offset0_limit0",
			offset: 0,
			limit:  0,
		},
		{
			name:   "out_offset0_limit10",
			offset: 0,
			limit:  10,
		},
		{
			name:   "out_offset0_limit1000",
			offset: 0,
			limit:  1000,
		},
		{
			name:   "out_offset0_limit10000",
			offset: 0,
			limit:  10000,
		},
		{
			name:   "out_offset100_limit1000",
			offset: 100,
			limit:  1000,
		},
		{
			name:   "out_offset6000_limit1000",
			offset: 6000,
			limit:  1000,
		},
	}

	for _, tc := range testCases {
		inputFileName := path.Join("testdata", "input.txt")
		t.Run(tc.name, func(t *testing.T) {
			dstFile, err := os.CreateTemp("/tmp", "test")
			assert.NoError(t, err)
			defer os.Remove(dstFile.Name())

			err = Copy(inputFileName, dstFile.Name(), tc.offset, tc.limit)
			require.NoError(t, err)

			refFileName := path.Join("testdata", strings.Join([]string{tc.name, "txt"}, "."))
			refFile, err := os.Open(refFileName)
			assert.NoError(t, err)
			defer refFile.Close()

			refContent, err := io.ReadAll(refFile)
			assert.NoError(t, err)

			resultContent, err := io.ReadAll(dstFile)
			assert.NoError(t, err)
			defer dstFile.Close()

			require.Equal(t, refContent, resultContent)
		})
	}
}
