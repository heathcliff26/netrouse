package gui

import (
	"testing"

	"fyne.io/fyne/v2/widget"
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

func TestMacEntryChangedFunction(t *testing.T) {
	tMatrix := []struct {
		Name              string
		Input             string
		CursorPos         int
		Expected          string
		ExpectedCursorPos int
	}{
		{
			Name:              "Empty",
			Input:             "",
			CursorPos:         11,
			Expected:          "",
			ExpectedCursorPos: 0,
		},
		{
			Name:              "NoColonsValidMac",
			Input:             "aabbccddeeff",
			CursorPos:         11,
			Expected:          "aa:bb:cc:dd:ee:ff",
			ExpectedCursorPos: 16,
		},
		{
			Name:              "CursorInMiddle",
			Input:             "aabbccddeeff",
			CursorPos:         2,
			Expected:          "aa:bb:cc:dd:ee:ff",
			ExpectedCursorPos: 3,
		},
		{
			Name:              "ToLong",
			Input:             "aabbccddeeff00112233",
			CursorPos:         17,
			Expected:          "aa:bb:cc:dd:ee:ff",
			ExpectedCursorPos: 17,
		},
		{
			Name:              "CursorOutOfBounds",
			Input:             "aa:bb:cc:dd:ee:ff",
			CursorPos:         100,
			Expected:          "aa:bb:cc:dd:ee:ff",
			ExpectedCursorPos: 17,
		},
		{
			Name:              "ColonsWrongPlaces",
			Input:             "aa:bb:cc:dd:e:ff",
			CursorPos:         0,
			Expected:          "aa:bb:cc:dd:ef:f",
			ExpectedCursorPos: 0,
		},
		{
			Name:              "InvalidChars",
			Input:             "aa:bb:ccg:dd:r:ff",
			CursorPos:         17,
			Expected:          "aa:bb:cc:dd:ff",
			ExpectedCursorPos: 14,
		},
	}

	entry := widget.NewEntry()
	entry.OnChanged = macEntryChangedFunction(entry)

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			entry.Text = tCase.Input
			entry.CursorColumn = tCase.CursorPos
			entry.OnChanged(entry.Text)

			assert.Equal(tCase.Expected, entry.Text)
			assert.Equal(tCase.ExpectedCursorPos, entry.CursorColumn)
		})
	}
}
