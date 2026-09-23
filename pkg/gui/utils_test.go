package gui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlAddSchema(t *testing.T) {
	tMatrix := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "HTTPS",
			Input:    "https://example.com",
			Expected: "https://example.com",
		},
		{
			Name:     "HTTP",
			Input:    "http://example.com",
			Expected: "http://example.com",
		},
		{
			Name:     "None",
			Input:    "example.com",
			Expected: "https://example.com",
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert.Equal(t, tCase.Expected, urlAddSchema(tCase.Input))
		})
	}
}
