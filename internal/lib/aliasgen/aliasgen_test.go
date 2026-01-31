package aliasgen_test

import (
	"testing"

	"github.com/sillkiw/url-shorten/internal/lib/aliasgen"
)

func TestNewRandomString(t *testing.T) {
	type testCase struct {
		name string
		size int
	}

	cases := []testCase{
		{
			name: "len 1",
			size: 1,
		},
		{
			name: "len 10",
			size: 10,
		},
		{
			name: "len 3",
			size: 3,
		},
		{
			name: "len 25",
			size: 25,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := len(aliasgen.GenRandomString(tc.size))
			if actual != tc.size {
				t.Fatalf("expected %d, got %d", tc.size, actual)
			}

		})
	}
}
