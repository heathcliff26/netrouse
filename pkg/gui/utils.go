package gui

import (
	"strings"

	"fyne.io/fyne/v2/widget"
)

func urlAddSchema(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "https://" + url
	}
	return url
}

func macEntryChangedFunction(entry *widget.Entry) func(string) {
	return func(s string) {
		cursorPos := entry.CursorColumn
		if len(s) > 17 {
			s = s[:17]
		}
		if cursorPos > len(s) {
			cursorPos = len(s)
		}
		i := 0
		for i < len(s) {
			if i == 2 || i == 5 || i == 8 || i == 11 || i == 14 {
				if s[i] != ':' {
					s = s[:i] + ":" + s[i:]
					if cursorPos >= i {
						cursorPos++
					}
				}
			} else if !strings.Contains("0123456789abcdefABCDEF", string(s[i])) {
				s = s[:i] + s[i+1:]
				if cursorPos > i {
					cursorPos--
				}
				continue
			}
			i++
		}
		entry.CursorColumn = cursorPos
		entry.SetText(s)
	}
}
