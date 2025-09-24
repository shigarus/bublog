package bublog

import (
	"unicode"
)

// TextViewer takes Text and allows to have only a portion of it being displayed in given Private.Width and Private.Height
type TextViewer struct {
	// Inner representation of Text is already split by /n to make it easier to navigate.
	Text [][]rune

	// DisplayText has several layer structure.
	// Upper level represents a complete /n line from Text.
	// Level below is a sublines - several strings gotten from original to fix into Width.
	// Last level is rune representation of a string.
	DisplayText [][][]rune

	Width  int
	Height int

	// LineInText shows from which line current DisplayText is rendered.
	LineInText           int
	SubLineInDisplayText int

	StickedToBottom bool
}

func NewTextViewer(text []rune) *TextViewer {
	innerText := make([][]rune, 0)
	lastStart := 0
	for i, ch := range text {
		if ch == '\n' {
			innerText = append(innerText, text[lastStart:i])
			lastStart = i + 1
		}
	}
	if lastStart < len(text) {
		innerText = append(innerText, text[lastStart:len(text)])
	}
	return &TextViewer{Text: innerText}
}

func (t *TextViewer) View() [][]rune {
	// Returns exactly t.Width * t.Height arrays.
	// If there is no enough symbols in a line - fills them with spaces.
	result := make([][]rune, 0, t.Height)
	curLinePos, curSublinePos := 0, t.SubLineInDisplayText
	for i := 0; i < t.Height; i++ {
		line := make([]rune, 0, t.Width)
		if curLinePos >= len(t.DisplayText) {
			for j := 0; j < t.Width; j++ {
				line = append(line, ' ')
			}
		} else {
			curLine := t.DisplayText[curLinePos]
			curSubline := curLine[curSublinePos]
			line = append(line, curSubline...)
			for i := 0; i < t.Width-len(curSubline); i++ {
				line = append(line, ' ')
			}
			if curSublinePos == len(curLine)-1 {
				curLinePos += 1
				curSublinePos = 0
			} else {
				curSublinePos += 1
			}
		}
		result = append(result, line)
	}
	return result
}

func (t *TextViewer) ScrollTextUp() bool {
	if t.LineInText == 0 {
		return false
	}
	t.LineInText -= 1
	lineToSwap := SplitLine(t.Text[t.LineInText], t.Width)
	for i, l := range t.DisplayText {
		t.DisplayText[i], lineToSwap = lineToSwap, l
	}
	return true
}

func (t *TextViewer) ScrollTextDown() bool {
	if t.LineInText+len(t.DisplayText) >= len(t.Text) {
		return false
	}
	t.LineInText += 1
	for i := 1; i < len(t.DisplayText); i++ {
		t.DisplayText[i-1] = t.DisplayText[i]
	}
	t.DisplayText[len(t.DisplayText)-1] = SplitLine(t.Text[t.LineInText+t.Height], t.Width)
	return true
}

func (t *TextViewer) ScrollUp() bool {
	if len(t.DisplayText) == 0 {
		return false
	}
	if t.SubLineInDisplayText == 0 {
		if t.ScrollTextUp() {
			t.SubLineInDisplayText = len(t.DisplayText[0]) - 1
			return true
		}
		return false
	}
	t.SubLineInDisplayText -= 1
	return true
}

func (t *TextViewer) ScrollDown() bool {
	if len(t.DisplayText) == 0 {
		return false
	}
	if t.SubLineInDisplayText >= len(t.DisplayText[0])-1 {
		if t.ScrollTextDown() {
			t.SubLineInDisplayText = 0
			return true
		}
		return false
	}
	t.SubLineInDisplayText = +1
	return true
}

func (t *TextViewer) SwitchStickToBottom() {
	if t.StickedToBottom {
		t.StickedToBottom = false
	} else {
		t.StickedToBottom = true
		t.LineInText = len(t.Text) - t.Height
		for t.ScrollDown() {
		}
	}
}

func (t *TextViewer) AppendToText(s string) {
	t.Text = append(t.Text, []rune(s))
	if t.StickedToBottom {
		for t.ScrollDown() {
		}
	}
}

func (t *TextViewer) Recalculate() {
	if len(t.Text) == 0 {
		return
	}
	// There some over-allocation but it is to cut logic a bit.
	newDisplayText := make([][][]rune, 0, t.LineInText+t.Height)
	linesToDisplay := 0
	for i := t.LineInText; linesToDisplay <= t.Height; i += 1 {
		fittedLines := SplitLine(t.Text[i], t.Width)
		linesToDisplay += len(fittedLines)
		newDisplayText = append(newDisplayText, fittedLines)
	}
	t.DisplayText = newDisplayText
}

// SplitLine fits an array of []rune into several with max length of maxWidth.
func SplitLine(line []rune, maxWidth int) [][]rune {
	res := make([][]rune, 0)
	lastSpacePos := -1
	lastSplit := 0
	for i, ch := range line {
		if unicode.IsSpace(ch) {
			lastSpacePos = i
		}
		if i-lastSplit == maxWidth {
			if i == len(line)-1 || lastSpacePos < lastSplit {
				res = append(res, line[lastSplit:i+1])
				lastSplit = i + 1
			} else {
				res = append(res, line[lastSplit:lastSpacePos])
				lastSplit = lastSpacePos + 1
			}
		}
	}
	if lastSplit < len(line)-1 {
		res = append(res, line[lastSplit:])
	}
	return res
}

func (t *TextViewer) SetWidth(w int) {
	if w == t.Width {
		return
	}
	t.Width = w
	t.SubLineInDisplayText = 0
	t.Recalculate()
}

func (t *TextViewer) SetHeight(h int) {
	if h == t.Height {
		return
	}
	t.Height = h
	t.SubLineInDisplayText = 0
	t.Recalculate()
}
